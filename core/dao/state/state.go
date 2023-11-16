package state

import (
	v8 "github.com/go-redis/redis/v8"
)

type State struct {
	client *v8.Client
}

func New(r *v8.Client) *State {
	return &State {
		client: r,
	}
}
