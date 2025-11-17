package repositories

import (
	"context"
	"encoding/json"

	"github.com/Vighnesh-V-H/sync/internal/db"
	"github.com/rs/zerolog"
)

type EventRepository struct {
	db  *db.DB
	log zerolog.Logger
}

func NewEventRepository(db *db.DB, log zerolog.Logger) *EventRepository {
	return &EventRepository{
		db:  db,
		log: log,
	}
}

func (r *EventRepository) AddEvent(ctx context.Context, apiKey string, id string, payload map[string]interface{}) error {
	
	r.log.Info().
		Str("api_key", apiKey).
		Str("id", id).
		Interface("payload", payload).
		Msg("Event received")

	
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		r.log.Error().Err(err).Msg("Failed to marshal payload")
		return err
	}

	r.log.Info().
		Str("api_key", apiKey).
		Str("id", id).
		Str("payload_json", string(payloadJSON)).
		Msg("Event processed successfully")

	return nil
}
