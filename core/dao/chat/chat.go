package chat

import (
	v8 "github.com/go-redis/redis/v8"
)

type Chat struct {
	client *v8.Client
}

func New(r *v8.Client) *Chat {
	return &Chat {
		client: r,
	}
}
