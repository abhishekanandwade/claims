package main

import (
	"context"
	"gitlab.cept.gov.in/pli/claims-api/bootstrap"

	bootstrapper "gitlab.cept.gov.in/it-2.0-common/n-api-bootstrapper"
)

func main() {
	app := bootstrapper.New().Options(
		// Add your FX modules here
		bootstrap.FxHandler, // Register all handlers
		bootstrap.FxRepo,    // Register all repositories
		// bootstrap.Fxvalidator, // Optional: custom validators
	)
	app.WithContext(context.Background()).Run()
}
