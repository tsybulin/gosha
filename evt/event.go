package evt

import "fmt"

// State ...
type State struct {
	EntityID   string                 `json:"entity_id"`
	State      string                 `json:"state"`
	Attributes map[string]interface{} `json:"attributes"`
}

func (s *State) String() string {
	return fmt.Sprintf("State{entity_id: %v, state: %v, attributes: %v}", s.EntityID, s.State, s.Attributes)
}

// Data ...
type Data struct {
	EntityID     string `json:"entity_id"`
	OldState     *State `json:"old_state,omitempty"`
	NewState     *State `json:"new_state,omitempty"`
	CurrentState *State `json:"state,omitempty"`
}

func (d *Data) String() string {
	s := fmt.Sprintf("Data{entity_id: %v", d.EntityID)

	if d.OldState != nil {
		s = s + fmt.Sprintf(", old_state:%v", d.OldState.String())
	}

	if d.CurrentState != nil {
		s = s + fmt.Sprintf(", state:%v", d.CurrentState.String())
	}

	if d.NewState != nil {
		s = s + fmt.Sprintf(", new_state:%v", d.NewState.String())
	}

	s = s + "}"
	return s
}

// Event ...
type Event struct {
	EventType string `json:"event_type"`
	Data      Data   `json:"data"`
}

func (e *Event) String() string {
	return fmt.Sprintf("Event{event_type: %v, data: %v}", e.EventType, e.Data.String())
}

// Message ...
type Message struct {
	ID    int    `json:"id,omitempty"`
	Type  string `json:"type"`
	Event *Event `json:"event,omitempty"`
	Token string `json:"token,omitempty"`
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{id: %v, type: % v, event: %v}", m.ID, m.Type, m.Event.String())
}
