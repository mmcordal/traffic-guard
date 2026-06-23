package router

import (
	"traffic-guarder/internal/handler"
	"traffic-guarder/internal/infrastructure/app"
	"traffic-guarder/internal/infrastructure/router"
	"traffic-guarder/internal/repository"
	"traffic-guarder/internal/service"
)

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (Router) RegisterRouter(a *app.App) {

	app := a.FiberApp
	db := a.DB
	bc := a.Bc
	br := a.Br
	as := a.As

	// Repositories
	tr := repository.NewTrafficRepository(db)

	// Services
	bs := service.NewBucketService(br, bc)
	ts := service.NewTrafficService(tr, bs)

	_ = as // UNUTMA

	// Handlers
	th := handler.NewTrafficHandler(ts)

	v1 := app.Group("/api/v1")

	log := v1.Group("/traffic-log")
	router.Put(log, "/", th.CreateTrafficLog)

}
