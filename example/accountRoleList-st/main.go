package main

import (
	"x-server/example/accountRoleList-st/tools"
)

// -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=100
func main() {
	tools.PreparePlatformAccount()
	tools.AccountRoleList()
	tools.GetLoginToken()
}
