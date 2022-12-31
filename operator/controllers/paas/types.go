package paas

const (
	eventOutOfWindow = "OutOfWindow"
	eventUpdated     = "Updated"
)

const (
	ConditionInWindow          = "InMaintenanceWindow"
	ConditionInWindowReasonOut = "OutOfWindow"
	ConditionInWindowReasonIn  = "InWindow"

	ConditionVersionReady            = "VersionReady"
	ConditionVersionReadyReasonFound = "Found"
	ConditionVersionReadyReasonErr   = "Error"
)
