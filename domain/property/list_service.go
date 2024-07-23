package property

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type ListRequest struct {
	// we might add some filters here.
	// or can also add pagination, etc.
}

type ListResponse struct {
	Message    string         `json:"message"`
	Properties []PropertyData `json:"properties"`
}

type Lister interface {
	List(context.Context, *ListRequest) (*ListResponse, error)
}

func NewLister(repo PropertyRepository) *ListerService {
	return &ListerService{
		repo: repo,
	}
}

type ListerService struct {
	repo PropertyRepository
}

func (h *ListerService) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	properties, err := h.repo.ListAll(ctx)
	if err != nil {
		// log the error in this layer.
		log.Println(err.Error())
		return nil, errors.New("failed to get properties data")
	}

	// set full address:
	for i := range properties {
		properties[i].FullAddress = FullAddress(properties[i])
	}

	return &ListResponse{
		Message:    fmt.Sprintf("Got %d properties data.", len(properties)),
		Properties: properties,
	}, nil
}
