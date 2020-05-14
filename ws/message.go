package ws

type serviceData struct {
	EntityID string `json:"entity_id"`
}

// Command ...
type Command struct {
	ID          int         `json:"id"`
	Type        string      `json:"type"`
	Domain      string      `json:"domain"`
	Command     string      `json:"service"`
	ServiceData serviceData `json:"service_data"`
}

// Result ...
type Result struct {
	ID      int    `json:"id,omitempty"`
	Type    string `json:"type"`
	Success string `json:"success"`
}
