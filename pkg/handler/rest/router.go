package rest

import (
	"net/http"

	"github.com/hossein1376/spallet/pkg/application/service"
	"github.com/hossein1376/spallet/pkg/handler/rest/swaggerhndlr"
	"github.com/hossein1376/spallet/pkg/handler/rest/usershndlr"
	"github.com/hossein1376/spallet/pkg/handler/rest/wallethndlr"
)

func routes(services *service.Services) *http.ServeMux {
	r := http.NewServeMux()

	swagger := swaggerhndlr.New()
	wallet := wallethndlr.New(services)
	user := usershndlr.New(services)

	r.Handle("GET /swagger/", withDefauls(swagger.SwaggerHandler))

	r.Handle(
		"POST /wallets/{user_id}/topup", withDefauls(wallet.TopUpHandler),
	)
	r.Handle(
		"POST /wallets/{user_id}/withdraw", withDefauls(wallet.WithdrawHandler),
	)
	r.Handle(
		"GET /wallets/{user_id}/balance", withDefauls(wallet.BalanceHandler),
	)
	r.Handle(
		"GET /wallets/{user_id}/transactions", withDefauls(wallet.HistoryHandler),
	)

	r.Handle("POST /users", withDefauls((user.CreateUserHandler)))

	return r
}

func withDefauls(handler http.HandlerFunc) http.Handler {
	return withMiddlewares(
		handler, requestIDMiddleware, loggerMiddleware, recoverMiddleware,
	)
}
