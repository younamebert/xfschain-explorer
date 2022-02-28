package router

import (
	"xfschainbrowser/router/accounts"
	"xfschainbrowser/router/blocks"
	"xfschainbrowser/router/home"
	"xfschainbrowser/router/transfer"
)

type RouterGroup struct {
	HomeRouter     home.HomeRouterGroup
	BlocksRouter   blocks.BlocksRouterGroup
	TxsRouter      transfer.TxsRouterGroup
	AccountsRouter accounts.AccountsRouterGroup
}

var RouterGroupApp = new(RouterGroup)
