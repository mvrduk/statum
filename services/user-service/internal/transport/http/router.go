package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mvrduk/statum/services/user-service/internal/transport/http/gen"
	"github.com/mvrduk/statum/services/user-service/internal/transport/http/handler"
)

func New(userHandler *handler.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/swagger", handler.SwaggerUI)
	r.Get("/openapi.yaml", handler.OpenAPISpec)
	return gen.HandlerFromMux(userHandler, r)
}
