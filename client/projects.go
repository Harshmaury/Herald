package client

import (
	"context"
	"fmt"

	accord "github.com/Harshmaury/Accord/api"
)

// ProjectsClient provides typed access to the /projects API.
type ProjectsClient struct{ c *Client }

// Get returns a single project by ID.
func (p *ProjectsClient) Get(ctx context.Context, id string) (*accord.ProjectDTO, error) {
	var out accord.ProjectDTO
	if err := p.c.get(ctx, "/projects/"+id, &out); err != nil {
		return nil, fmt.Errorf("projects.get %s: %w", id, err)
	}
	return &out, nil
}

// List returns all registered projects.
func (p *ProjectsClient) List(ctx context.Context) ([]accord.ProjectDTO, error) {
	var out []accord.ProjectDTO
	if err := p.c.get(ctx, "/projects", &out); err != nil {
		return nil, fmt.Errorf("projects.list: %w", err)
	}
	return out, nil
}

// Start sets desired=running for all services in a project.
func (p *ProjectsClient) Start(ctx context.Context, id string) (*accord.ProjectActionData, error) {
	var out accord.ProjectActionData
	if err := p.c.post(ctx, "/projects/"+id+"/start", nil, &out); err != nil {
		return nil, fmt.Errorf("projects.start %s: %w", id, err)
	}
	return &out, nil
}

// Stop sets desired=stopped for all services in a project.
func (p *ProjectsClient) Stop(ctx context.Context, id string) (*accord.ProjectActionData, error) {
	var out accord.ProjectActionData
	if err := p.c.post(ctx, "/projects/"+id+"/stop", nil, &out); err != nil {
		return nil, fmt.Errorf("projects.stop %s: %w", id, err)
	}
	return &out, nil
}

// Exists returns true if the project is registered.
func (p *ProjectsClient) Exists(ctx context.Context, id string) (bool, error) {
	_, err := p.Get(ctx, id)
	if err == nil {
		return true, nil
	}
	if e, ok := err.(*accord.Error); ok && e.Code == accord.ErrNotFound {
		return false, nil
	}
	return false, err
}
