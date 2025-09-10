package swaggerhndlr

import (
	"net/http"
)

type SwaggerHandler struct {
}

var swaggerPath = "assets/docs/openapi/swagger.html"

func New() SwaggerHandler {
	return SwaggerHandler{}
}

func (h *SwaggerHandler) SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, swaggerPath)
}
