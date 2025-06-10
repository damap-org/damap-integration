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
	GetProjects(ctx context.Context) ([]domain.Project, error)
	GetStaffForProject(ctx context.Context, projectId string) ([]domain.Person, error)
}

type client struct {
	lowLevelClient cli.ClientWithResponsesInterface
}

func (c client) GetProjects(ctx context.Context) ([]domain.Project, error) {
	resp, err := c.lowLevelClient.GetProjectsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return *resp.JSON200, nil
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
		if resp.JSON502 != nil {
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
