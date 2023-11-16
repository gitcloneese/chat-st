package model

import (
	"fmt"
)

// account token out date second
const AccountTokenTime = 1800

// token format
const AccountTokenFmt = "%s:%d:%d:%d:%d"

// redis key
const RedisKeyAccountToken = "account:token:{%s}"

func GetAccountTokenKey(accountId string) string {
	return fmt.Sprintf(RedisKeyAccountToken, accountId)
}
