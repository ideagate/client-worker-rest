package job

import (
	"context"

	entityContext "github.com/bayu-aditya/ideagate/backend/core/model/entity/context"
	entityDataSource "github.com/bayu-aditya/ideagate/backend/core/model/entity/datasource"
	pbEndpoint "github.com/bayu-aditya/ideagate/backend/model/gen-go/core/endpoint"
	"github.com/gin-gonic/gin"
)

type StartInput struct {
	Ctx        context.Context
	GinCtx     *gin.Context
	DataCtx    *entityContext.ContextData
	DataSource *entityDataSource.DataSource
	Endpoint   *pbEndpoint.Endpoint
	Step       *pbEndpoint.Step
}

type StartOutput struct {
	NextStepIds []string
	Data        any
}
