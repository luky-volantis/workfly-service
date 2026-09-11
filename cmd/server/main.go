package main

import (
	"log"

	"Slackluky/workfly/internal/service"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// @title           Workfly API
// @version         1.0
// @description     Go + PocketBase + Temporal backend API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@workfly.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		log.Println("Server is starting on http://0.0.0.0:8080...")

		// Run schema setup and data seeding on boot
		if err := service.SeedData(app); err != nil {
			log.Printf("Warning: SeedData returned error: %v\n", err)
		}
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
