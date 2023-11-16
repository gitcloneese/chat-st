package scene

import (
	v8 "github.com/go-redis/redis/v8"
)

type Scene struct {
	client *v8.Client
}

func New(r *v8.Client) *Scene {
	return &Scene{
		client: r,
	}
}


