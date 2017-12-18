package main

// StringService provides operations on strings.
import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/benmanns/goworker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/time/rate"
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

func loggingEndpointMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", "calling endpoint")
			defer logger.Log("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	tsvc := targetService{}
	triggerEndpoint := MakeTriggerHandlerEndpoint(tsvc)
	triggerEndpoint = loggingEndpointMiddleware(logger)(triggerEndpoint)
	triggerEndpoint = ratelimit.NewDelayingLimiter(rate.NewLimiter(2, 1))(triggerEndpoint)

	triggerHandler := httptransport.NewServer(
		triggerEndpoint,
		DecodeTriggerRequest,
		encodeResponse,
	)

	http.Handle("/trigger", triggerHandler)

	logger.Log("err", http.ListenAndServe(":8080", nil))
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
