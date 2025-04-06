package job

import (
	"context"

	"github.com/gin-gonic/gin"
	entityContext "github.com/ideagate/core/model/entity/context"
	entityDataSource "github.com/ideagate/core/model/entity/datasource"
	pbEndpoint "github.com/ideagate/model/gen-go/core/endpoint"
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
