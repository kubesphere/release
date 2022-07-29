package notes

type Area string

const (
	// Observability
	AreaAlerting      Area = "alerting"
	AreaAuditing      Area = "auditing"
	AreaBilling       Area = "billing"
	AreaLogging       Area = "logging"
	AreaMetering      Area = "metering"
	AreaMonitoring    Area = "monitoring"
	AreaNotification  Area = "notification"
	AreaObservability Area = "observability"
	//
	AreaObservabilityAlias = "Observability"

	// from /kind api-change
	AreaAPIChange      Area = "api-change"
	AreaAPIChangeAlias Area = "API Changes"

	// APP Store
	AreaAppManagement Area = "app-management"
	AreaApps          Area = "apps"
	AreaAppAlias      Area = "APP Store"

	// devops
	AreaDevOps      Area = "devops"
	AreaDevOpsAlias Area = "DevOps"

	// Development & Testing
	AreaE2ETestFramework Area = "e2e-test-framework"
	AreaInfra            Area = "infra"
	AreaTesting          Area = "testing"
	AreaDevTestAlias     Area = "Development & Testing"

	// Edge
	AreaEdge      Area = "edge"
	AreaEdgeAlias Area = "KubeEdge Intergration"

	// Multi-tenancy & Multi-cluster
	AreaIAM                  Area = "iam"
	AreaMulticluster         Area = "multicluster"
	AreaIAMMulticlusterAlias Area = "Multi-tenancy & Multi-cluster"

	// Service Mesh
	AreaMicroService      Area = "microservice"
	AreaMicroServiceAlias Area = "Service Mesh"

	// Network
	AreaNetWork      Area = "networking"
	AreaNetWorkAlias Area = "Network"
	// Storage
	AreaStorage      Area = "storage"
	AreaStorageAlias Area = "Storage"

	// User Experience
	AreaUI      Area = "ui"
	AreaUIAlias Area = "User Experience"

	// Security
	AreaSecurity      Area = "security"
	AreaSecurityAlias Area = "Security"

	AreaUncategorized Area = "Uncategorized"
)
