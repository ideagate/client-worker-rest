package main

import (
	"context"

	"github.com/bayu-aditya/ideagate/backend/client/worker-rest/adapter/controller"
	"github.com/bayu-aditya/ideagate/backend/client/worker-rest/config"
	"github.com/bayu-aditya/ideagate/backend/client/worker-rest/usecase/handler"
	"github.com/bayu-aditya/ideagate/backend/core/utils/log"
	"github.com/gin-gonic/gin"
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
