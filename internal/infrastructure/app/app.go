package app

import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"traffic-guarder/internal/cron"
	"traffic-guarder/internal/infrastructure/cache"
	"traffic-guarder/internal/infrastructure/config"
	"traffic-guarder/internal/infrastructure/database"
	"traffic-guarder/internal/infrastructure/errorsx"
	"traffic-guarder/internal/repository"
	"traffic-guarder/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/uptrace/bun"
)

type App struct {
	FiberApp *fiber.App
	DB       *bun.DB
	Cfg      *config.Config
	Bc       cache.BucketCache
	As       service.AnomalyService
	Br       repository.BucketRepository
}

type IRouter interface {
	RegisterRouter(app *App)
}

func New(router IRouter) *App {
	cfg, err := config.Setup()
	if err != nil {
		panic(err)
	}

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: errorsx.ErrorHandler,
	})

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,http://127.0.0.1:5173", // http://localhost:5173 || http://127.0.0.1:5173 || http://---IP---:5173
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Currency",
	}))

	db := database.New(cfg.Database)

	redisClient := cache.NewRedisClient(
		cfg.Redis.Host + ":" + cfg.Redis.Port,
	)

	bc := cache.NewBucketCache(redisClient, cfg.Analyze)

	// repositories
	br := repository.NewBucketRepository(db)
	dc := repository.NewDomainCheck(db)
	ar := repository.NewAnomalyRepository(db)

	as := service.NewAnomalyService(ar, dc, bc, br, cfg.Analyze)

	app := &App{
		FiberApp: fiberApp,
		DB:       db,
		Cfg:      cfg,
		Bc:       bc,
		As:       as,
		Br:       br,
	}

	router.RegisterRouter(app)

	return app
}

func (a *App) Start() {
	cron.Start(a.As, a.Cfg.Analyze)

	go func() {

		err := a.FiberApp.Listen(fmt.Sprintf(":%v", a.Cfg.Server.Port))
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}
