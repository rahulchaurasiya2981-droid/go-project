package logger

const (
	// Application lifecycle messages.
	MsgAppStart    = "Application starting"
	MsgAppReady    = "Application ready"
	MsgAppShutdown = "Application shutting down"

	// Environment loading messages.
	MsgEnvLoadStart = "Starting environment load"
	MsgEnvLoaded    = "Environment loaded successfully"
	MsgEnvLoadError = "Failed to load environment file"

	// Configuration loading messages.
	MsgConfigLoadStart = "Starting configuration load"
	MsgConfigLoaded    = "Configuration loaded successfully"
	MsgConfigLoadError = "Failed to load configuration"

	// Database messages.
	MsgDatabaseConnectionStart = "Connecting to database"
	MsgDatabaseConnected       = "Database connection established successfully"
	MsgDatabaseConnectionError = "Database connection could not be established"
	MsgDatabaseHealthCheck     = "Running database health check"
	MsgDatabaseHealthOK        = "Database health check passed"
	MsgDatabaseHealthFailed    = "Database health check failed"

	// Server messages.
	MsgServerStart   = "Starting HTTP server"
	MsgServerStarted = "HTTP server listening"
	MsgServerStop    = "Stopping HTTP server"
	MsgServerStopped = "HTTP server stopped"

	// HTTP messages.
	MsgHTTPRequest      = "HTTP request received"
	MsgHTTPResponse     = "HTTP response sent"
	MsgHTTPRequestError = "HTTP request failed"

	// CRUD messages.
	MsgCRUDCreate        = "Creating resource"
	MsgCRUDRead          = "Reading resource"
	MsgCRUDUpdate        = "Updating resource"
	MsgCRUDDelete        = "Deleting resource"
	MsgCRUDReadAll       = "Reading all resources"
	MsgCRUDCreateSuccess = "Resource created successfully"
	MsgCRUDReadSuccess   = "Resource retrieved successfully"
	MsgCRUDUpdateSuccess = "Resource updated successfully"
	MsgCRUDDeleteSuccess = "Resource deleted successfully"
	MsgCRUDCreateError   = "Failed to create resource"
	MsgCRUDReadError     = "Failed to read resource"
	MsgCRUDUpdateError   = "Failed to update resource"
	MsgCRUDDeleteError   = "Failed to delete resource"
)
