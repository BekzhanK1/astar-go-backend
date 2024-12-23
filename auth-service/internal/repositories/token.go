package repositories

import (
	"auth-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

func (r *TokenRepository) SaveRefreshToken(userID uint, token models.RefreshToken, ttl time.Duration) error {
	key := fmt.Sprintf("refresh_token:%d", userID)

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to serialize token: %w", err)
	}

	return r.client.Set(context.Background(), key, data, ttl).Err()
}

func (r *TokenRepository) GetRefreshToken(userID uint) (*models.RefreshToken, error) {
	key := fmt.Sprintf("refresh_token:%d", userID)

	data, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("token not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch token: %w", err)
	}

	var token models.RefreshToken
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to deserialize token: %w", err)
	}

	return &token, nil
}

func (r *TokenRepository) DeleteRefreshToken(userID uint) error {
	key := fmt.Sprintf("refresh_token:%d", userID)
	return r.client.Del(context.Background(), key).Err()
}

func (r *TokenRepository) ValidateRefreshToken(userID uint, tokenHash string) (bool, error) {
	token, err := r.GetRefreshToken(userID)
	if err != nil {
		return false, err
	}
	return token.TokenHash == tokenHash, nil
}
