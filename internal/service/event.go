package service

import (
	"context"

	"github.com/Vighnesh-V-H/sync/internal/repositories"
	"github.com/rs/zerolog"
)

type EventService struct {
	repo   *repositories.EventRepository
	logger zerolog.Logger
}

func NewEventService(repo *repositories.EventRepository, logger zerolog.Logger) *EventService {
	return &EventService{
		repo:   repo,
		logger: logger.With().Str("service", "event").Logger(),
	}
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
	s.logger.Debug().
		Str("event_id", req.ID).
		Str("api_key", req.ApiKey).
		Int("payload_size", len(req.Payload)).
		Msg("Processing event addition")

	err := s.repo.AddEvent(ctx, req.ApiKey, req.ID, req.Payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("event_id", req.ID).
			Str("api_key", req.ApiKey).
			Msg("Failed to add event to repository")
		return nil, err
	}

	s.logger.Info().
		Str("event_id", req.ID).
		Str("api_key", req.ApiKey).
		Msg("Event persisted successfully")

	return &AddEventResponse{
		Success: true,
		Message: "Event added successfully",
	}, nil
}
