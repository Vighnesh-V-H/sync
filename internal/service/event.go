package service

import (
	"context"

	"github.com/Vighnesh-V-H/sync/internal/repositories"
)

type EventService struct {
	repo *repositories.EventRepository
}

func NewEventService(repo *repositories.EventRepository) *EventService {
	return &EventService{repo: repo}
}

type AddEventRequest struct {
	ApiKey  string                 `json:"api_key"`
	ID      string                 `json:"id"`
	Payload map[string]interface{} `json:"payload"`
}

type AddEventResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (s *EventService) AddEvent(ctx context.Context, req AddEventRequest) (*AddEventResponse, error) {
	err := s.repo.AddEvent(ctx, req.ApiKey, req.ID, req.Payload)
	if err != nil {
		return nil, err
	}

	return &AddEventResponse{
		Success: true,
		Message: "Event added successfully",
	}, nil
}
