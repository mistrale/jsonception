package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

func init() {
	revel.OnAppStart(InitDB)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}

func (c App) Index() revel.Result {
	return c.Render()
}
