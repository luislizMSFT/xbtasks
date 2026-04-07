package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"dev.azure.com/xbox/xb-tasks/domain"
	"github.com/spf13/viper"
)

const (
	AppName    = "xb-tasks"
	ConfigName = "config"
	ConfigType = "yaml"
)

// Init initializes Viper with defaults, config file, and env overrides.
// Call once at startup before accessing any config values.
func Init() error {
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)
	viper.AddConfigPath(ConfigDir())

	setDefaults()

	viper.SetEnvPrefix("XBT")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// No config file yet — write defaults
			if err := ensureConfigDir(); err != nil {
				return fmt.Errorf("create config dir: %w", err)
			}
			return viper.SafeWriteConfig()
		}
		return fmt.Errorf("read config: %w", err)
	}
	return nil
}

func setDefaults() {
	viper.SetDefault("db.path", filepath.Join(DataDir(), "data.db"))
	viper.SetDefault("theme", "system")
	viper.SetDefault("window.width", 1440)
	viper.SetDefault("window.height", 900)
	viper.SetDefault("ado.organization", "")
	viper.SetDefault("ado.project", "")
	viper.SetDefault("ado.pat_keychain_key", "xbt-ado-pat")
	viper.SetDefault("ado.orgs", []map[string]any{})
	viper.SetDefault("sync.interval_minutes", 15)
	viper.SetDefault("log.level", "info")
}

// ConfigDir returns the OS-appropriate config directory.
//
//	macOS:   ~/Library/Application Support/xb-tasks
//	Windows: %APPDATA%\xb-tasks
//	Linux:   ~/.config/xb-tasks
func ConfigDir() string {
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", AppName)
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), AppName)
	default:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, AppName)
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", AppName)
	}
}

// DataDir returns the OS-appropriate data directory.
//
//	macOS:   ~/Library/Application Support/xb-tasks
//	Windows: %LOCALAPPDATA%\xb-tasks
//	Linux:   ~/.local/share/xb-tasks
func DataDir() string {
	switch runtime.GOOS {
	case "darwin":
		return ConfigDir() // macOS keeps data with config
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), AppName)
	default:
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			return filepath.Join(xdg, AppName)
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".local", "share", AppName)
	}
}

func ensureConfigDir() error {
	return os.MkdirAll(ConfigDir(), 0755)
}

// Convenience accessors

func DBPath() string            { return viper.GetString("db.path") }
func Theme() string             { return viper.GetString("theme") }
func WindowWidth() int          { return viper.GetInt("window.width") }
func WindowHeight() int         { return viper.GetInt("window.height") }
func ADOOrganization() string   { return viper.GetString("ado.organization") }
func ADOProject() string        { return viper.GetString("ado.project") }
func SyncIntervalMinutes() int  { return viper.GetInt("sync.interval_minutes") }
func LogLevel() string          { return viper.GetString("log.level") }

// Set writes a config key and persists to disk.
func Set(key string, value any) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}

// GetOrgProjects returns all configured org/project pairs.
// Falls back to legacy single-org config if ado.orgs is empty.
func GetOrgProjects() []domain.OrgProject {
	var orgs []domain.OrgProject
	viper.UnmarshalKey("ado.orgs", &orgs)
	if len(orgs) == 0 {
		org := viper.GetString("ado.organization")
		proj := viper.GetString("ado.project")
		if org != "" && proj != "" {
			orgs = []domain.OrgProject{{Org: org, Projects: []string{proj}}}
		}
	}
	return orgs
}

// SetOrgProjects writes the org/project list to config.
func SetOrgProjects(orgs []domain.OrgProject) error {
	viper.Set("ado.orgs", orgs)
	return viper.WriteConfig()
}
