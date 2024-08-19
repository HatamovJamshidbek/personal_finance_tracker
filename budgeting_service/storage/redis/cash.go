package redis

import (
	pbb "budgeting_service/genproto/budgeting_service"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRepo struct {
	Db *redis.Client
}

func NewRedisRepo(db *redis.Client) *RedisRepo {
	return &RedisRepo{
		Db: db,
	}
}

func (redisDb *RedisRepo) AddBalanceInCache(ctx context.Context, request *pbb.Account, expiration time.Duration) error {
	balanceKey := fmt.Sprintf("balance:%s", request.Id)

	err := redisDb.Db.HIncrByFloat(ctx, balanceKey, "balance", float64(request.Balance)).Err()
	if err != nil {
		return fmt.Errorf("failed to increment balance in cache: %v", err)
	}

	err = redisDb.Db.Expire(ctx, balanceKey, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for balance key: %v", err)
	}

	return nil
}
