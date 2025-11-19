package repositories

import (
	"context"
	"encoding/json"

	"github.com/Vighnesh-V-H/sync/internal/db"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type EventRepository struct {
	db    *db.DB
	redis *redis.Client
	log   zerolog.Logger
}

func NewEventRepository(db *db.DB, redisClient *redis.Client, log zerolog.Logger) *EventRepository {
	return &EventRepository{
		db:    db,
		redis: redisClient,
		log:   log.With().Str("repository", "event").Logger(),
	}
}

func (r *EventRepository) AddEvent(ctx context.Context, apiKey string, id string, payload map[string]interface{}) error {
	r.log.Debug().
		Str("api_key", apiKey).
		Str("event_id", id).
		Msg("Receiving event for processing")

	eventData := map[string]interface{}{
		"api_key": apiKey,
		"id":      id,
		"payload": payload,
	}

	payloadJSON, err := json.Marshal(eventData)
	if err != nil {
		r.log.Error().Err(err).
			Str("event_id", id).
			Msg("Failed to marshal event payload")
		return err
	}

	queueKey := "events:queue"
	err = r.redis.RPush(ctx, queueKey, payloadJSON).Err()
	if err != nil {
		r.log.Error().Err(err).
			Str("event_id", id).
			Str("queue_key", queueKey).
			Msg("Failed to push event to Redis queue")
		return err
	}

	r.log.Info().
		Str("api_key", apiKey).
		Str("event_id", id).
		Str("queue_key", queueKey).
		Int("payload_size", len(payloadJSON)).
		Msg("Event queued successfully in Redis")

	return nil
}
