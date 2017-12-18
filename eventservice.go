package main

import (
	"context"
	"fmt"

	"github.com/benmanns/goworker"
)

type (
	// SlackMessage gago eh
	SlackMessage struct {
		Text        string                   `json:"text"`
		Attachments []SlackMessageAttachment `json:"attachments"`
	}

	// SlackMessageAttachment gago e
	SlackMessageAttachment struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}
)

type DoesNotQueueService interface {
	SendResponse(context.Context, string) (SlackMessage, error)
}

type helloService struct{}

func (helloService) SendResponse(_ context.Context, command string) (SlackMessage, error) {
	return SlackMessage{Text: fmt.Sprintf("Hello ahole you typed `%s`", command)}, nil
}

// TargetService Gago
type TargetService interface {
	MakeJob(context.Context, string, string) error
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
