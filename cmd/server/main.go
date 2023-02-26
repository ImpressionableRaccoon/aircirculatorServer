package main

import (
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/aircirculatorServer/configs"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/handlers"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/middlewares"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/routers"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/service"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/workers"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := configs.NewConfig()

	st, err := storage.NewPsqlStorage(cfg)
	if err != nil {
		panic(err)
	}

	s := service.NewService(st, cfg)

	h := handlers.NewHandler(s, cfg)

	m := middlewares.NewMiddlewares(s, h, cfg)

	r := routers.NewRouter(h, m)

	w := workers.NewWorkers(s, cfg)
	defer w.Stop()

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
