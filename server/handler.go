package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/damap-org/damap-integration/backend"
	"github.com/damap-org/damap-integration/domain"
	api "github.com/damap-org/damap-integration/server/internal"
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

func (h handler) GetProject(_ context.Context, request api.GetProjectRequestObject) (api.GetProjectResponseObject, error) {
	project, err := h.backend.GetProject(request.ProjectId)
	if err != nil {
		return nil, err
	}
	return api.GetProject200JSONResponse(*project), nil
}

func (h handler) GetStaffForProject(_ context.Context, request api.GetStaffForProjectRequestObject) (api.GetStaffForProjectResponseObject, error) {
	people, err := h.backend.GetStaffForProject(request.ProjectId)
	if err != nil {
		return nil, err
	}
	return api.GetStaffForProject200JSONResponse(people), nil
}

func (h handler) GetLeaderForProject(_ context.Context, request api.GetLeaderForProjectRequestObject) (api.GetLeaderForProjectResponseObject, error) {
	leader, err := h.backend.GetLeaderForProject(request.ProjectId)
	if err != nil {
		return nil, err
	}
	return api.GetLeaderForProject200JSONResponse(*leader), nil
}

func (h handler) GetRecommendedProjects(_ context.Context, request api.GetRecommendedProjectsRequestObject) (api.GetRecommendedProjectsResponseObject, error) {
	projects, err := h.backend.GetRecommendedProjects(request.Params.Query)
	if err != nil {
		return nil, err
	}
	return api.GetRecommendedProjects200JSONResponse(projects), nil
}

func (h handler) SearchProjects(_ context.Context, request api.SearchProjectsRequestObject) (api.SearchProjectsResponseObject, error) {
	projects, err := h.backend.SearchProjects(request.Params.Query)
	if err != nil {
		return nil, err
	}
	return api.SearchProjects200JSONResponse(projects), nil
}

func (h handler) GetProjectSupplements(_ context.Context, request api.GetProjectSupplementsRequestObject) (api.GetProjectSupplementsResponseObject, error) {
	supplements, err := h.backend.GetProjectSupplements(request.ProjectId)
	if err != nil {
		return nil, err
	}
	return api.GetProjectSupplements200JSONResponse(*supplements), nil
}

func (h handler) GetPerson(_ context.Context, request api.GetPersonRequestObject) (api.GetPersonResponseObject, error) {
	person, err := h.backend.GetPerson(request.PersonId)
	if err != nil {
		return nil, err
	}
	return api.GetPerson200JSONResponse(*person), nil
}

func (h handler) SearchPersons(_ context.Context, request api.SearchPersonsRequestObject) (api.SearchPersonsResponseObject, error) {
	persons, err := h.backend.SearchPersons(request.Params.Query)
	if err != nil {
		return nil, err
	}
	return api.SearchPersons200JSONResponse(persons), nil
}
