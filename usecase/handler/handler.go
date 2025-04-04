package handler

import (
	"context"

	adapterController "github.com/bayu-aditya/ideagate/backend/client/worker-rest/adapter/controller"
	"github.com/bayu-aditya/ideagate/backend/core/utils/errors"
	pbEndpoint "github.com/bayu-aditya/ideagate/backend/model/gen-go/core/endpoint"
	"github.com/gin-gonic/gin"
)

type IHandlerUsecase interface {
	GenerateEndpoint(ctx context.Context, router *gin.Engine) error
}

func New(adapterController adapterController.IControllerAdapter) IHandlerUsecase {
	return &handler{
		prefix:            "handler",
		adapterController: adapterController,
	}
}

type handler struct {
	prefix            string
	adapterController adapterController.IControllerAdapter
}

func (h *handler) GenerateEndpoint(ctx context.Context, router *gin.Engine) error {
	prefix := h.prefix + ".GenerateEndpoint"

	resultListEndpoint, err := h.adapterController.GetListEndpoint(ctx)
	if err != nil {
		return errors.Wrap(prefix, err, "get list endpoint")
	}

	for _, endpointPb := range resultListEndpoint.GetEndpoints() {
		setting := endpointPb.GetSettingRest()
		if setting == nil {
			continue
		}

		router.Handle(setting.GetMethod(), setting.GetPath(), h.handler(h.adapterController, endpointPb))
	}

	return nil
}

func (h *handler) handler(adapterController adapterController.IControllerAdapter, endpoint *pbEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get workflow from controller
		respWorkflow, err := adapterController.GetWorkflow(c.Request.Context(), endpoint.GetId())
		if err != nil {
			return
		}

		mgr, _ := newManager(c, endpoint, respWorkflow.GetWorkflow())
		mgr.RunHandler()
	}
}
