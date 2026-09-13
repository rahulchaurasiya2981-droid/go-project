package logger

const (
	// Application lifecycle actions.
	ActionAppStart    = "APP_START"
	ActionAppReady    = "APP_READY"
	ActionAppShutdown = "APP_SHUTDOWN"

	// Environment loading actions.
	ActionEnvLoadStart = "ENV_LOAD_START"
	ActionEnvLoaded    = "ENV_LOADED"
	ActionEnvLoadError = "ENV_LOAD_ERROR"

	// Configuration loading actions.
	ActionConfigLoadStart = "CONFIG_LOAD_START"
	ActionConfigLoaded    = "CONFIG_LOADED"
	ActionConfigLoadError = "CONFIG_LOAD_ERROR"

	// Database actions.
	ActionDatabaseConnectionStart = "DB_CONNECTION_START"
	ActionDatabaseConnected       = "DB_CONNECTED"
	ActionDatabaseConnectionError = "DB_CONNECTION_ERROR"
	ActionDatabaseHealthCheck     = "DB_HEALTH_CHECK"
	ActionDatabaseHealthOK        = "DB_HEALTH_OK"
	ActionDatabaseHealthFailed    = "DB_HEALTH_FAILED"

	// Server actions.
	ActionServerStart   = "SERVER_START"
	ActionServerStarted = "SERVER_STARTED"
	ActionServerStop    = "SERVER_STOP"
	ActionServerStopped = "SERVER_STOPPED"

	// HTTP actions.
	ActionHTTPRequest      = "HTTP_REQUEST"
	ActionHTTPResponse     = "HTTP_RESPONSE"
	ActionHTTPRequestError = "HTTP_REQUEST_ERROR"

	// CRUD actions.
	ActionCRUDCreate        = "CRUD_CREATE"
	ActionCRUDRead          = "CRUD_READ"
	ActionCRUDUpdate        = "CRUD_UPDATE"
	ActionCRUDDelete        = "CRUD_DELETE"
	ActionCRUDReadAll       = "CRUD_READ_ALL"
	ActionCRUDCreateSuccess = "CRUD_CREATE_SUCCESS"
	ActionCRUDReadSuccess   = "CRUD_READ_SUCCESS"
	ActionCRUDUpdateSuccess = "CRUD_UPDATE_SUCCESS"
	ActionCRUDDeleteSuccess = "CRUD_DELETE_SUCCESS"
	ActionCRUDCreateError   = "CRUD_CREATE_ERROR"
	ActionCRUDReadError     = "CRUD_READ_ERROR"
	ActionCRUDUpdateError   = "CRUD_UPDATE_ERROR"
	ActionCRUDDeleteError   = "CRUD_DELETE_ERROR"
)
