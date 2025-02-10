package web

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/momentum/internal/dto"
	"github.com/uploadpilot/uploadpilot/momentum/internal/utils"
)

type handler struct {
}

func Newhandler() *handler {
	return &handler{}
}

func (h *handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "momentum is healthy")
}

func (h *handler) TriggerWorkflow(w http.ResponseWriter, r *http.Request) {
	req := &dto.TriggerworkflowReq{}
	if err := render.DecodeJSON(r.Body, req); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	infra.Log.Infof("Triggering workflow: %+v", req)

	render.JSON(w, r, &dto.TriggerworkflowResp{WorkflowID: "test"})
}
