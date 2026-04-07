package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"

	"dev.azure.com/xbox/xb-tasks/internal/db"
	"dev.azure.com/xbox/xb-tasks/domain"
)

const (
	serviceName     = "team-ado-tool"
	defaultTenantID = "YOUR_TENANT_ID"
	defaultClientID = "YOUR_CLIENT_ID"
	callbackTimeout = 5 * time.Minute
	graphMeEndpoint = "https://graph.microsoft.com/v1.0/me"
	graphPhotoURL   = "https://graph.microsoft.com/v1.0/me/photo/$value"
)

type AuthService struct {
	db          *db.DB
	oauthConfig *oauth2.Config
	currentUser *domain.User
	app         *application.App
}

func NewAuthService(database *db.DB, app *application.App) *AuthService {
	tenantID := os.Getenv("TEAM_ADO_TENANT_ID")
	if tenantID == "" {
		tenantID = defaultTenantID
	}
	clientID := os.Getenv("TEAM_ADO_CLIENT_ID")
	if clientID == "" {
		clientID = defaultClientID
	}

	return &AuthService{
		db:  database,
		app: app,
		oauthConfig: &oauth2.Config{
			ClientID: clientID,
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
				TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
			},
			Scopes: []string{"User.Read", "offline_access"},
		},
	}
}

// SignIn starts the OAuth2 PKCE flow, opens the browser, and returns the authenticated user.
func (s *AuthService) SignIn() (*domain.User, error) {
	verifier := generateCodeVerifier()
	challenge := generateCodeChallenge(verifier)

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("start callback server: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	s.oauthConfig.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", port)

	authURL := s.oauthConfig.AuthCodeURL("state",
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no authorization code in callback")
			fmt.Fprint(w, "<html><body><h1>Error</h1><p>No authorization code received.</p></body></html>")
			return
		}
		codeCh <- code
		fmt.Fprint(w, "<html><body><h1>Signed in</h1><p>You can close this tab and return to the app.</p></body></html>")
	})

	server := &http.Server{Handler: mux}
	go func() {
		if serveErr := server.Serve(listener); serveErr != http.ErrServerClosed {
			errCh <- serveErr
		}
	}()
	defer server.Close()

	if err := browser.OpenURL(authURL); err != nil {
		return nil, fmt.Errorf("open browser: %w", err)
	}

	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return nil, fmt.Errorf("callback error: %w", err)
	case <-time.After(callbackTimeout):
		return nil, fmt.Errorf("authentication timed out after %v", callbackTimeout)
	}

	token, err := s.oauthConfig.Exchange(context.Background(), code,
		oauth2.SetAuthURLParam("code_verifier", verifier),
	)
	if err != nil {
		return nil, fmt.Errorf("exchange token: %w", err)
	}

	if err := saveTokens(token); err != nil {
		return nil, fmt.Errorf("save tokens: %w", err)
	}

	user, err := fetchUserProfile(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("fetch profile: %w", err)
	}

	if err := s.upsertUser(user); err != nil {
		return nil, fmt.Errorf("store user: %w", err)
	}

	s.currentUser = user
	s.app.Event.Emit("auth:state-changed", map[string]any{"authenticated": true})
	return user, nil
}

// SignOut clears stored tokens and resets the current user.
func (s *AuthService) SignOut() error {
	_ = keyring.Delete(serviceName, "access_token")
	_ = keyring.Delete(serviceName, "refresh_token")
	_ = keyring.Delete(serviceName, "token_expiry")
	_ = keyring.Delete(serviceName, "pat")
	_ = keyring.Delete(serviceName, "auth_method")
	s.currentUser = nil
	s.app.Event.Emit("auth:state-changed", map[string]any{"authenticated": false})
	return nil
}

// TryRestoreSession attempts to restore a previous session using stored refresh token.
// Returns (nil, nil) if no session exists.
func (s *AuthService) TryRestoreSession() (*domain.User, error) {
	// Check auth method first
	authMethod, _ := keyring.Get(serviceName, "auth_method")
	if authMethod == "azcli" {
		return s.SignInWithAzCli()
	}

	refreshToken, err := keyring.Get(serviceName, "refresh_token")
	if err != nil {
		// Check for PAT session
		pat, patErr := keyring.Get(serviceName, "pat")
		if patErr != nil {
			return nil, nil
		}
		return s.restoreFromPAT(pat)
	}

	tokenSource := s.oauthConfig.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: refreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		// Token expired or revoked — clean up
		_ = keyring.Delete(serviceName, "access_token")
		_ = keyring.Delete(serviceName, "refresh_token")
		_ = keyring.Delete(serviceName, "token_expiry")
		return nil, nil
	}

	if err := saveTokens(newToken); err != nil {
		return nil, fmt.Errorf("save refreshed tokens: %w", err)
	}

	user, err := fetchUserProfile(newToken.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("fetch profile after refresh: %w", err)
	}

	if err := s.upsertUser(user); err != nil {
		return nil, fmt.Errorf("store user: %w", err)
	}

	s.currentUser = user
	return user, nil
}

// GetCurrentUser returns the currently authenticated user, or nil if not signed in.
func (s *AuthService) GetCurrentUser() *domain.User {
	return s.currentUser
}

// IsAuthenticated returns whether a user is currently signed in.
func (s *AuthService) IsAuthenticated() bool {
	return s.currentUser != nil
}

// SignInWithPAT authenticates using a personal access token for development.
func (s *AuthService) SignInWithPAT(pat string) (*domain.User, error) {
	if pat == "" {
		return nil, fmt.Errorf("PAT cannot be empty")
	}

	if err := keyring.Set(serviceName, "pat", pat); err != nil {
		return nil, fmt.Errorf("store PAT: %w", err)
	}

	user := &domain.User{
		ID:          "pat-user",
		DisplayName: "PAT User",
		Email:       "",
		AvatarURL:   "",
	}

	if err := s.upsertUser(user); err != nil {
		return nil, fmt.Errorf("store user: %w", err)
	}

	s.currentUser = user
	s.app.Event.Emit("auth:state-changed", map[string]any{"authenticated": true})
	return user, nil
}

// SignInWithAzCli authenticates using the az CLI token for ADO access.
// Requires the user to have run `az login` first.
func (s *AuthService) SignInWithAzCli() (*domain.User, error) {
	provider := NewAzCliTokenProvider()
	token, err := provider.GetToken()
	if err != nil {
		return nil, fmt.Errorf("az CLI auth failed: %w", err)
	}

	// Try ADO profile API first, fall back to az account show
	user, err := fetchADOProfile(token)
	if err != nil {
		user, err = fetchAzAccountUser()
		if err != nil {
			return nil, fmt.Errorf("could not resolve user identity: %w", err)
		}
	}

	if err := s.upsertUser(user); err != nil {
		return nil, fmt.Errorf("store user: %w", err)
	}

	// Store auth method for session restore
	_ = keyring.Set(serviceName, "auth_method", "azcli")
	s.currentUser = user
	s.app.Event.Emit("auth:state-changed", map[string]any{"authenticated": true})
	return user, nil
}

// fetchADOProfile uses an ADO token to get the current user's profile.
func fetchADOProfile(token string) (*domain.User, error) {
	req, err := http.NewRequest("GET", "https://app.vssps.visualstudio.com/_apis/profile/profiles/me?api-version=7.0", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ADO profile request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ADO profile API returned %d: %s", resp.StatusCode, string(body))
	}

	var profile struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		EmailAddress string `json:"emailAddress"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode ADO profile: %w", err)
	}

	return &domain.User{
		ID:          profile.ID,
		DisplayName: profile.DisplayName,
		Email:       profile.EmailAddress,
		AvatarURL:   "",
	}, nil
}

// fetchAzAccountUser gets the signed-in user from `az account show`.
func fetchAzAccountUser() (*domain.User, error) {
	azBin := resolveAzPath()
	cmd := exec.Command(azBin, "account", "show", "--output", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("az account show failed: %w", err)
	}

	var acct struct {
		User struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"user"`
		ID string `json:"id"`
	}
	if err := json.Unmarshal(output, &acct); err != nil {
		return nil, fmt.Errorf("parse az account: %w", err)
	}

	return &domain.User{
		ID:          acct.ID,
		DisplayName: acct.User.Name,
		Email:       acct.User.Name,
	}, nil
}

func (s *AuthService) restoreFromPAT(pat string) (*domain.User, error) {
	_ = pat // PAT is stored; presence confirms session
	user := &domain.User{
		ID:          "pat-user",
		DisplayName: "PAT User",
		Email:       "",
		AvatarURL:   "",
	}
	s.currentUser = user
	return user, nil
}

func generateCodeVerifier() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func saveTokens(token *oauth2.Token) error {
	if err := keyring.Set(serviceName, "access_token", token.AccessToken); err != nil {
		return fmt.Errorf("save access_token: %w", err)
	}
	if token.RefreshToken != "" {
		if err := keyring.Set(serviceName, "refresh_token", token.RefreshToken); err != nil {
			return fmt.Errorf("save refresh_token: %w", err)
		}
	}
	if !token.Expiry.IsZero() {
		if err := keyring.Set(serviceName, "token_expiry", token.Expiry.Format(time.RFC3339)); err != nil {
			return fmt.Errorf("save token_expiry: %w", err)
		}
	}
	return nil
}

type graphMeResponse struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
}

func fetchUserProfile(accessToken string) (*domain.User, error) {
	req, err := http.NewRequest("GET", graphMeEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("graph API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("graph API returned %d: %s", resp.StatusCode, string(body))
	}

	var profile graphMeResponse
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}

	email := profile.Mail
	if email == "" {
		email = profile.UserPrincipalName
	}

	return &domain.User{
		ID:          profile.ID,
		DisplayName: profile.DisplayName,
		Email:       email,
		AvatarURL:   graphPhotoURL,
	}, nil
}

func (s *AuthService) upsertUser(user *domain.User) error {
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO users (id, display_name, email, avatar_url) VALUES (?, ?, ?, ?)`,
		user.ID, user.DisplayName, user.Email, user.AvatarURL,
	)
	return err
}
