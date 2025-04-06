package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ideagate/client-worker-rest/adapter/controller"
	"github.com/ideagate/client-worker-rest/config"
	"github.com/ideagate/client-worker-rest/usecase/handler"
	"github.com/ideagate/core/utils/log"
)

func main() {
	ctx := context.Background()
	router := gin.Default()

	if err := config.Init(); err != nil {
		log.Fatal("init config: %v", err)
	}

	adapterController, err := controller.New()
	if err != nil {
		log.Fatal("init controller adapter: %v", err)
	}

	usecaseHandler := handler.New(adapterController)
	if err := usecaseHandler.GenerateEndpoint(ctx, router); err != nil {
		log.Fatal("generate endpoint: %v", err)
	}

	if err := router.Run(); err != nil {
		log.Fatal("run router: %v", err)
	}
}
