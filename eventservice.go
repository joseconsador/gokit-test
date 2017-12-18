package main

import (
	"context"

	"github.com/benmanns/goworker"
)

// TargetService Gago
type TargetService interface {
	MakeJob(context.Context, string, string) error
}

// StringService Gago
type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}

type targetService struct{}

func (targetService) MakeJob(_ context.Context, status string, ticketID string) error {
	return goworker.Enqueue(&goworker.Job{
		Queue: "myqueue",
		Payload: goworker.Payload{
			Class: "MyClass",
			Args:  []interface{}{status, ticketID},
		},
	})
}

// EventService gago
type EventService interface {
	GetEvent(context.Context, EventRequest) (string, error)
	CreateTicket(context.Context, EventRequest) (string, error)
}

type eventService struct{}

func (eventService) GetEvent(_ context.Context, e EventRequest) (string, error) {
	return e.Type, nil
}
