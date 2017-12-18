package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/mitchellh/mapstructure"
)

func MakeTriggerHandlerEndpoint(svc TargetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TriggerRequest)
		err := svc.MakeJob(ctx, req.Status, req.TicketID)

		if err != nil {
			return TriggerResponse{Message: "failed", Err: err.Error()}, nil
		}

		return TriggerResponse{Message: "ok", Err: ""}, nil
	}
}

func DecodeTriggerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request TriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeCommandRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CommandRequest
	r.ParseForm()

	formData := make(map[string]string)

	for k, v := range r.PostForm {
		formData[k] = strings.Join(v, "")
	}

	if err := mapstructure.Decode(formData, &request); err != nil {
		fmt.Println(request)
		return nil, err
	}

	return request, nil
}

func MakeCreateTicketEndpoint(svc EventService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EventRequest)
		v, err := svc.CreateTicket(ctx, req)
		if err != nil {
			return CreateTicketResponse{Message: v, Err: err.Error()}, nil
		}
		return CreateTicketResponse{Message: v, Err: ""}, nil
	}
}

func DecodeCreateTicketRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request EventRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func MakeHelpEndpoint(svc DoesNotQueueService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CommandRequest)
		slackResponse, err := svc.SendResponse(ctx, req.Command)

		if err != nil {
			return nil, err
		}

		return slackResponse, nil
	}
}
