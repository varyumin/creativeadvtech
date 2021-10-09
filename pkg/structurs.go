package pkg

import (
	"github.com/go-redis/redis"
	"net"
)

type ServerSettings struct {
	Host string
	Port string
}

type RedisSettings struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type BodyJson struct {
	Count int `json:"count"`
}

func (r *RedisSettings) Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(r.Host, r.Port),
		Password: r.Password,
		DB:       r.DB,
	})
}

func (b *BodyJson) Flush() {
	b.Count = 0
}
