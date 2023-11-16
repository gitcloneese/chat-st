package scene

import (
	"context"
	"errors"
	"fmt"

	"xy3-proto/pkg/log"
	"x-server/core/dao/model"
)

// VerifyToken .
func (s *Scene) VerifyToken(playerID int64, stoken string) (err error) {
	if playerID == 0 || stoken == "" {
		return errors.New("params wrong")
	}

	aKey := fmt.Sprintf(model.AccessTokenKey, playerID)
	dtoken, err := s.client.Get(context.TODO(), aKey).Result()
	if err != nil {
		log.Error("GetValue token fail, err: %s", err.Error())
		return err
	}

	if stoken != dtoken {
		return errors.New("token not match")
	}

	return nil
}
