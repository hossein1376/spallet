package usershndlr

import (
	"net/http"

	"github.com/hossein1376/spallet/pkg/handler/rest/serde"
	"github.com/hossein1376/spallet/pkg/service"
)

type UsersHandler struct {
	services *service.Services
}

func New(services *service.Services) *UsersHandler {
	return &UsersHandler{services: services}
}

func (h UsersHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateUserRequest
	err := serde.ReadJson(r, &req)
	if err != nil {
		resp := serde.Response{Message: err.Error()}
		serde.WriteJson(ctx, w, http.StatusBadRequest, resp)
		return
	}

	user, err := h.services.Users.CreateUserService(ctx, req.Username)
	if err != nil {
		status, resp := serde.ExtractFromErr(ctx, err)
		serde.WriteJson(ctx, w, status, resp)
		return
	}

	serde.WriteJson(ctx, w, http.StatusCreated, serde.Response{Data: user})
}
