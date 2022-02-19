package router

import "xfschainbrowser/router/home"

type RouterGroup struct {
	HomeRouter home.HomeRouterGroup
}

var RouterGroupApp = new(RouterGroup)
