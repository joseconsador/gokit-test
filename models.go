package main

// TriggerRequest gago
type TriggerRequest struct {
	Status    string `json:"status"`
	TicketID  int    `json:"id"`
	UpdatedAt string `json:"updated_at"`
}

// TriggerResponse gago
type TriggerResponse struct {
	Message string `json:"message"`
	Err     string `json:"error"`
}

// TicketComment struct
type (
	TicketComment struct {
		Body string `json:"body"`
	}

	// Ticket struct
	Ticket struct {
		Subject string        `json:"subject"`
		Comment TicketComment `json:"comment"`
	}

	// WeirdTicketWrapper struct
	WeirdTicketWrapper struct {
		Ticket Ticket `json:"ticket"`
	}
)

// Event Gago
type Event struct {
	Channel string `json:"channel"`
	Ts      string `json:"ts"`
	Text    string `json:"text"`
	EventTs string `json:"event_ts"`
	Type    string `json:"type"`
	User    string `json:"user"`
}

// EventRequest gago
type EventRequest struct {
	TeamID      string   `json:"team_id"`
	Event       Event    `json:"event"`
	APIAppID    string   `json:"api_app_id"`
	AuthedUsers []string `json:"authed_users"`
	EventTime   int      `json:"event_time"`
	Token       string   `json:"token"`
	Type        string   `json:"type"`
	EventID     string   `json:"event_id"`
}

// CreateTicketResponse gago
type CreateTicketResponse struct {
	Message string `json:"message"`
	Err     string `json:"error"`
}
