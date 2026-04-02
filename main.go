package main

import (
	"embed"
	"log"

	"dev.azure.com/xbox/xb-tasks/internal/app"
	"dev.azure.com/xbox/xb-tasks/internal/auth"
	"dev.azure.com/xbox/xb-tasks/internal/config"
	"dev.azure.com/xbox/xb-tasks/internal/db"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	database, err := db.Open(config.DBPath())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	configService := config.NewConfigService()
	taskService := app.NewTaskService(database)
	projectService := app.NewProjectService(database)
	depService := app.NewDependencyService(database)

	wailsApp := application.New(application.Options{
		Name:        "team-ado-tool",
		Description: "Unified team dashboard — tasks, ADO, PRs in one pane",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// Register services after app creation so authService can reference wailsApp
	authService := auth.NewAuthService(database, wailsApp)
	wailsApp.RegisterService(application.NewService(configService))
	wailsApp.RegisterService(application.NewService(taskService))
	wailsApp.RegisterService(application.NewService(projectService))
	wailsApp.RegisterService(application.NewService(depService))
	wailsApp.RegisterService(application.NewService(authService))

	// System tray
	tray := wailsApp.SystemTray.New()
	trayMenu := wailsApp.NewMenu()
	trayMenu.Add("Show").OnClick(func(ctx *application.Context) {
		if w, ok := wailsApp.Window.Get("main"); ok {
			w.Show()
			w.Focus()
		}
	})
	trayMenu.AddSeparator()
	trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		wailsApp.Quit()
	})
	tray.SetMenu(trayMenu)
	tray.OnClick(func() {
		if w, ok := wailsApp.Window.Get("main"); ok {
			w.Show()
			w.Focus()
		}
	})

	mainWindow := wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "main",
		Title:  "Team ADO Tool",
		Width:  config.WindowWidth(),
		Height: config.WindowHeight(),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 32,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// Intercept window close to hide instead of quit
	mainWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		e.Cancel()
		mainWindow.Hide()
	})

	// Restore session in background
	go func() {
		if _, err := authService.TryRestoreSession(); err != nil {
			log.Printf("session restore: %v", err)
		}
	}()

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}
