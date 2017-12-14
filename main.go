package main

// StringService provides operations on strings.
import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/benmanns/goworker"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

// ErrEmpty is returned when input string is empty
var ErrEmpty = errors.New("Empty string")

func init() {
	settings := goworker.WorkerSettings{
		URI:            "redis://localhost:6379/",
		Connections:    100,
		Queues:         []string{"myqueue", "delimited", "queues"},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    2,
		Namespace:      "resque:",
		Interval:       5.0,
	}
	goworker.SetSettings(settings)
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	esvc := eventService{}
	tsvc := targetService{}

	createTicketHandlerEndpoint := MakeCreateTicketEndpoint(esvc)
	//createTicketHandlerEndpoint = loggingMiddleware(log.With(logger, "method", "createTicket"))(createTicketHandlerEndpoint)

	createTicketHandler := httptransport.NewServer(
		createTicketHandlerEndpoint,
		DecodeCreateTicketRequest,
		encodeResponse,
	)

	triggerHandler := httptransport.NewServer(
		MakeTriggerHandlerEndpoint(tsvc),
		DecodeTriggerRequest,
		encodeResponse,
	)

	http.Handle("/createTicket", createTicketHandler)
	http.Handle("/trigger", triggerHandler)

	logger.Log("err", http.ListenAndServe(":8080", nil))
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
