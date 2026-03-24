package plugin

// StartRequest is sent from shield to the plugin via stdin.
type StartRequest struct {
	Action string       `json:"action"` // "start" or "stop"
	Config PluginConfig `json:"config,omitempty"`
}

// PluginConfig contains the database connection details.
type PluginConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user,omitempty"`
	Pass     string `json:"pass,omitempty"`
	Database string `json:"database,omitempty"`
	ReadOnly bool   `json:"readonly,omitempty"`
}

// StartResponse is returned from the plugin via stdout.
type StartResponse struct {
	Status  string `json:"status"`            // "ready" or "error"
	WebPort int    `json:"web_port,omitempty"` // local port the plugin's web UI is listening on
	Name    string `json:"name,omitempty"`     // display name, e.g. "MySQL Web Client"
	Version string `json:"version,omitempty"`
	Message string `json:"message,omitempty"` // error message when status="error"
}
