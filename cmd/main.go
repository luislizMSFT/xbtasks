package main

import (
	"embed"
	"log"

	"dev.azure.com/xbox/xb-tasks/internal/app"
	"dev.azure.com/xbox/xb-tasks/internal/auth"
	"dev.azure.com/xbox/xb-tasks/internal/db"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	database, err := db.Open(db.DefaultDBPath())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	taskService := app.NewTaskService(database)
	projectService := app.NewProjectService(database)
	depService := app.NewDependencyService(database)

	wailsApp := application.New(application.Options{
		Name:        "team-ado-tool",
		Description: "Unified team dashboard — tasks, ADO, PRs in one pane",
		Services: []application.Service{
			application.NewService(taskService),
			application.NewService(projectService),
			application.NewService(depService),
			application.NewService(auth.NewAuthService(database, wailsApp)),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// System tray
	tray := wailsApp.NewSystemTray()
	trayMenu := wailsApp.NewMenu()
	trayMenu.Add("Show").OnClick(func(ctx *application.Context) {
		for _, w := range wailsApp.GetWindows() {
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
		for _, w := range wailsApp.GetWindows() {
			w.Show()
			w.Focus()
		}
	})

	mainWindow := wailsApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Team ADO Tool",
		Width:  1200,
		Height: 800,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour:       application.NewRGB(27, 38, 54),
		URL:                    "/",
		ShouldClose: func(window *application.WebviewWindow) bool {
			window.Hide()
			return false
		},
	})
	_ = mainWindow

	// Restore session in background
	authService := auth.NewAuthService(database, wailsApp)
	go func() {
		if _, err := authService.TryRestoreSession(); err != nil {
			log.Printf("session restore: %v", err)
		}
	}()

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}
