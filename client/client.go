package client

import (
	"context"
	"fmt"
	"net/http"

	cli "github.com/damap-org/damap-integration/client/internal"
	"github.com/damap-org/damap-integration/domain"
)

func New(endpoint string) (Client, error) {
	lowLevelClient, err := cli.NewClientWithResponses(endpoint)
	if err != nil {
		return nil, err
	}
	return &client{lowLevelClient}, nil
}

type Client interface {
	GetProject(ctx context.Context, projectId string) (*domain.Project, error)
	GetStaffForProject(ctx context.Context, projectId string) ([]domain.Person, error)
	GetLeaderForProject(ctx context.Context, projectId string) (*domain.Person, error)
	GetRecommendedProjects(ctx context.Context, params *domain.GetRecommendedProjectsParams) ([]domain.Project, error)
	SearchProjects(ctx context.Context, params *domain.SearchProjectsParams) ([]domain.Project, error)
	GetProjectSupplements(ctx context.Context, projectId string) (*domain.ProjectSupplement, error)

	GetPerson(ctx context.Context, personId string) (*domain.Person, error)
	SearchPersons(ctx context.Context, params *domain.SearchPersonsParams) ([]domain.Person, error)
}

type client struct {
	lowLevelClient cli.ClientWithResponsesInterface
}

func (c client) GetProject(ctx context.Context, projectId string) (*domain.Project, error) {
	resp, err := c.lowLevelClient.GetProjectWithResponse(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	case http.StatusBadRequest:
		if resp.JSON400 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON400)
		}
		return nil, fmt.Errorf("request failed with 400 status code")
	case http.StatusInternalServerError:
		if resp.JSON500 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON500)
		}
		return nil, fmt.Errorf("request failed with 500 status code")
	case http.StatusBadGateway:
		if resp.JSON502 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON502)
		}
		return nil, fmt.Errorf("request failed with 502 status code")
	default:
		return nil, fmt.Errorf("request failed with unexpected status %d", resp.StatusCode())
	}
}

func (c client) GetStaffForProject(ctx context.Context, projectId string) ([]domain.Person, error) {
	resp, err := c.lowLevelClient.GetStaffForProjectWithResponse(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return *resp.JSON200, nil
	case http.StatusBadRequest:
		if resp.JSON400 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON400)
		}
		return nil, fmt.Errorf("request failed with 400 status code")
	case http.StatusInternalServerError:
		if resp.JSON500 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON500)
		}
		return nil, fmt.Errorf("request failed with 500 status code")
	case http.StatusBadGateway:
		if resp.JSON502 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON502)
		}
		return nil, fmt.Errorf("request failed with 502 status code")
	default:
		return nil, fmt.Errorf("request failed with unexpected status %d", resp.StatusCode())
	}
}

func (c client) GetLeaderForProject(ctx context.Context, projectId string) (*domain.Person, error) {
	resp, err := c.lowLevelClient.GetLeaderForProjectWithResponse(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	case http.StatusBadRequest:
		if resp.JSON400 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON400)
		}
		return nil, fmt.Errorf("request failed with 400 status code")
	case http.StatusInternalServerError:
		if resp.JSON500 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON500)
		}
		return nil, fmt.Errorf("request failed with 500 status code")
	case http.StatusBadGateway:
		if resp.JSON502 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON502)
		}
		return nil, fmt.Errorf("request failed with 502 status code")
	default:
		return nil, fmt.Errorf("request failed with unexpected status %d", resp.StatusCode())
	}
}

func (c client) GetRecommendedProjects(ctx context.Context, query *domain.GetRecommendedProjectsParams) ([]domain.Project, error) {
	resp, err := c.lowLevelClient.GetRecommendedProjectsWithResponse(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return *resp.JSON200, nil
	case http.StatusBadRequest:
		if resp.JSON400 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON400)
		}
		return nil, fmt.Errorf("request failed with 400 status code")
	case http.StatusInternalServerError:
		if resp.JSON500 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON500)
		}
		return nil, fmt.Errorf("request failed with 500 status code")
	case http.StatusBadGateway:
		if resp.JSON502 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON502)
		}
		return nil, fmt.Errorf("request failed with 502 status code")
	default:
		return nil, fmt.Errorf("request failed with unexpected status %d", resp.StatusCode())
	}
}

func (c client) SearchProjects(ctx context.Context, query *domain.SearchProjectsParams) ([]domain.Project, error) {
	resp, err := c.lowLevelClient.SearchProjectsWithResponse(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return *resp.JSON200, nil
	case http.StatusBadRequest:
		if resp.JSON400 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON400)
		}
		return nil, fmt.Errorf("request failed with 400 status code")
	case http.StatusInternalServerError:
		if resp.JSON500 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON500)
		}
		return nil, fmt.Errorf("request failed with 500 status code")
	case http.StatusBadGateway:
		if resp.JSON502 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON502)
		}
		return nil, fmt.Errorf("request failed with 502 status code")
	default:
		return nil, fmt.Errorf("request failed with unexpected status %d", resp.StatusCode())
	}
}

func (c client) GetProjectSupplements(ctx context.Context, projectId string) (*domain.ProjectSupplement, error) {
	resp, err := c.lowLevelClient.GetProjectSupplementsWithResponse(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	case http.StatusBadRequest:
		if resp.JSON400 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON400)
		}
		return nil, fmt.Errorf("request failed with 400 status code")
	case http.StatusInternalServerError:
		if resp.JSON500 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON500)
		}
		return nil, fmt.Errorf("request failed with 500 status code")
	case http.StatusBadGateway:
		if resp.JSON502 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON502)
		}
		return nil, fmt.Errorf("request failed with 502 status code")
	default:
		return nil, fmt.Errorf("request failed with unexpected status %d", resp.StatusCode())
	}
}

func (c client) GetPerson(ctx context.Context, personId string) (*domain.Person, error) {
	resp, err := c.lowLevelClient.GetPersonWithResponse(ctx, personId)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	case http.StatusBadRequest:
		if resp.JSON400 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON400)
		}
		return nil, fmt.Errorf("request failed with 400 status code")
	case http.StatusInternalServerError:
		if resp.JSON500 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON500)
		}
		return nil, fmt.Errorf("request failed with 500 status code")
	case http.StatusBadGateway:
		if resp.JSON502 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON502)
		}
		return nil, fmt.Errorf("request failed with 502 status code")
	default:
		return nil, fmt.Errorf("request failed with unexpected status %d", resp.StatusCode())
	}
}

func (c client) SearchPersons(ctx context.Context, query *domain.SearchPersonsParams) ([]domain.Person, error) {
	resp, err := c.lowLevelClient.SearchPersonsWithResponse(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return *resp.JSON200, nil
	case http.StatusBadRequest:
		if resp.JSON400 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON400)
		}
		return nil, fmt.Errorf("request failed with 400 status code")
	case http.StatusInternalServerError:
		if resp.JSON500 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON500)
		}
		return nil, fmt.Errorf("request failed with 500 status code")
	case http.StatusBadGateway:
		if resp.JSON502 != nil {
			return nil, fmt.Errorf("request failed (%w)", resp.JSON502)
		}
		return nil, fmt.Errorf("request failed with 502 status code")
	default:
		return nil, fmt.Errorf("request failed with unexpected status %d", resp.StatusCode())
	}
}
