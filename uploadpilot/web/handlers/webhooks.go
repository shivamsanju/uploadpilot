package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"github.com/uploadpilot/uploadpilot/internal/webhook"
)

type webhooksHandler struct {
	webhookSvc *webhook.WebhookService
}

func NewWebhooksHandler() *webhooksHandler {
	return &webhooksHandler{
		webhookSvc: webhook.NewWebhookService(),
	}
}

func (wh *webhooksHandler) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	webhooks, err := wh.webhookSvc.GetAllWebhooksInWorkspace(r.Context(), workspaceID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, webhooks)
}

func (wh *webhooksHandler) GetWebhookDetailsByID(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	webhookID := chi.URLParam(r, "webhookId")

	webhook, err := wh.webhookSvc.GetWebhook(r.Context(), workspaceID, webhookID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, webhook)
}

func (wh *webhooksHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	webhook := &models.Webhook{}
	if err := render.DecodeJSON(r.Body, webhook); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if err := wh.webhookSvc.CreateWebhook(r.Context(), workspaceID, webhook); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, webhook.ID)
}

func (wh *webhooksHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	webhookID := chi.URLParam(r, "webhookId")

	if err := wh.webhookSvc.DeleteWebhook(r.Context(), workspaceID, webhookID); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (wh *webhooksHandler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	webhookID := chi.URLParam(r, "webhookId")

	webhook := &models.Webhook{}
	if err := render.DecodeJSON(r.Body, webhook); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if err := wh.webhookSvc.UpdateWebhook(r.Context(), workspaceID, webhookID, webhook); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (wh *webhooksHandler) PatchWebhook(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	webhookID := chi.URLParam(r, "webhookId")

	patchReq := &dto.PatchWebhookRequest{}
	if err := render.DecodeJSON(r.Body, patchReq); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if err := wh.webhookSvc.PatchWebhook(r.Context(), workspaceID, webhookID, patchReq); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)
}
