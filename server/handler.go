package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tuwien-csd/damap-integration/backend"
	"github.com/tuwien-csd/damap-integration/domain"
	api "github.com/tuwien-csd/damap-integration/server/internal"
)

// NewHandler creates a net/http-compatible handler with the provided backend.
func NewHandler(backend backend.ConfiguredBackend) http.Handler {
	apiHandler := api.NewStrictHandler(&handler{
		backend: backend,
	}, nil)
	return api.HandlerWithOptions(apiHandler, api.StdHTTPServerOptions{
		BaseURL:    "/api",
		BaseRouter: http.NewServeMux(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			sendError := func(data any, status int) {
				body, err := json.Marshal(data)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				w.WriteHeader(status)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(body)
			}

			var badRequestError *domain.BadRequestError
			var internalError *domain.InternalError
			var badGatewayError *domain.BadGatewayError
			switch {
			case errors.As(err, &badRequestError):
				sendError(badRequestError, 400)
			case errors.As(err, &internalError):
				sendError(internalError, 500)
			case errors.As(err, &badGatewayError):
				sendError(badGatewayError, 502)
			default:
				w.WriteHeader(500)
			}
		},
	})
}

type handler struct {
	backend backend.ConfiguredBackend
}

func (h handler) GetProjects(_ context.Context, _ api.GetProjectsRequestObject) (api.GetProjectsResponseObject, error) {
	projects, err := h.backend.GetProjects()
	if err != nil {
		return nil, err
	}
	return api.GetProjects200JSONResponse(projects), nil
}

func (h handler) GetStaffForProject(_ context.Context, request api.GetStaffForProjectRequestObject) (api.GetStaffForProjectResponseObject, error) {
	people, err := h.backend.GetPeople(request.ProjectId)
	if err != nil {
		return nil, err
	}
	return api.GetStaffForProject200JSONResponse(people), nil
}
