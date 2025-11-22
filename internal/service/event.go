package service

import (
	"context"

	"github.com/Vighnesh-V-H/sync/internal/repositories"
	"github.com/google/uuid"
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
	Payload map[string]any `json:"payload"`
}

type AddEventResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	EventID string `json:"event_id"`
}

func (s *EventService) AddEvent(ctx context.Context, apiKey string, req AddEventRequest) (*AddEventResponse, error) {
	// Generate unique event ID
	eventID := uuid.New().String()

	s.logger.Debug().
		Str("event_id", eventID).
		Str("api_key", apiKey).
		Int("payload_size", len(req.Payload)).
		Msg("Processing event addition")

	err := s.repo.AddEvent(ctx, apiKey, eventID, req.Payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("event_id", eventID).
			Str("api_key", apiKey).
			Msg("Failed to add event to repository")
		return nil, err
	}

	s.logger.Info().
		Str("event_id", eventID).
		Str("api_key", apiKey).
		Msg("Event persisted successfully")

	return &AddEventResponse{
		Success: true,
		Message: "Event added successfully",
		EventID: eventID,
	}, nil
}
