package main

// type loggingMiddleware struct {
// 	logger log.Logger
// 	next   EventService
// }
//
// func loggingMiddleware(logger log.Logger) Middleware {
// 	return func(next endpoint.Endpoint) endpoint.Endpoint {
// 		return func(ctx context.Context, request interface{}) (interface{}, error) {
// 			logger.Log("msg", "calling endpoint")
// 			defer logger.Log("msg", "called endpoint")
// 			return next(ctx, request)
// 		}
// 	}
// }
