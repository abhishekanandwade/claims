package bootstrap

import (
	"go.uber.org/fx"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	handler "gitlab.cept.gov.in/pli/claims-api/handler"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// FxRepo module provides all repository implementations
var FxRepo = fx.Module(
	"Repomodule",
	fx.Provide(
		// Claim repositories
		repo.NewClaimRepository,
		repo.NewClaimDocumentRepository,
		repo.NewClaimPaymentRepository,
		repo.NewClaimHistoryRepository,
		repo.NewClaimCommunicationRepository,

		// Investigation repositories
		repo.NewInvestigationRepository,
		repo.NewInvestigationProgressRepository,

		// AML repository
		repo.NewAMLAlertRepository,

		// Appeal repository
		repo.NewAppealRepository,

		// Ombudsman repository
		repo.NewOmbudsmanComplaintRepository,

		// Policy bond tracking repository
		repo.NewPolicyBondTrackingRepository,

		// Free look cancellation repository
		repo.NewFreeLookCancellationRepository,

		// Document checklist repository
		repo.NewDocumentChecklistRepository,

		// SLA tracking repository
		repo.NewSLATrackingRepository,

		// Add more repository constructors here as needed
		// repo.New{Resource}Repository,
	),
)

// FxHandler module provides all HTTP handlers
var FxHandler = fx.Module(
	"Handlermodule",
	fx.Provide(
		// Each handler must be annotated to implement serverHandler.Handler interface

		// Claim handlers
		fx.Annotate(
			handler.NewClaimHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Investigation handlers
		fx.Annotate(
			handler.NewInvestigationHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Maturity claim handlers
		fx.Annotate(
			handler.NewMaturityClaimHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Survival benefit handlers
		fx.Annotate(
			handler.NewSurvivalBenefitHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// AML handlers
		fx.Annotate(
			handler.NewAMLHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Banking handlers
		fx.Annotate(
			handler.NewBankingHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Free look handlers
		fx.Annotate(
			handler.NewFreeLookHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Appeal handlers
		fx.Annotate(
			handler.NewAppealHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Ombudsman handlers
		fx.Annotate(
			handler.NewOmbudsmanHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Notification handlers
		fx.Annotate(
			handler.NewNotificationHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Policy service handlers
		fx.Annotate(
			handler.NewPolicyServiceHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Validation service handlers
		fx.Annotate(
			handler.NewValidationServiceHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Lookup handlers
		fx.Annotate(
			handler.NewLookupHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Report handlers
		fx.Annotate(
			handler.NewReportHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Workflow handlers
		fx.Annotate(
			handler.NewWorkflowHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Status and tracking handlers
		fx.Annotate(
			handler.NewStatusHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Add more handler constructors here as needed
		// fx.Annotate(
		//     handler.New{Resource}Handler,
		//     fx.As(new(serverHandler.Handler)),
		//     fx.ResultTags(serverHandler.ServerControllersGroupTag),
		// ),
	),
)

// Optional: Custom validator module (if using custom validators)
// var Fxvalidator = fx.Module(
//     "Validatormodule",
//     fx.Provide(
//         // Add custom validator providers here
//     ),
// )
