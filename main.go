package main

import (
	"embed"
	"log"

	"dev.azure.com/microsoft/Xbox/xb-tasks/internal/app"
	"dev.azure.com/microsoft/Xbox/xb-tasks/internal/db"
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

	wailsApp := application.New(application.Options{
		Name:        "team-ado-tool",
		Description: "Unified team dashboard — tasks, ADO, PRs in one pane",
		Services: []application.Service{
			application.NewService(taskService),
			application.NewService(projectService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Team ADO Tool",
		Width:  1200,
		Height: 800,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}
