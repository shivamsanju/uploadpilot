package handlers

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"github.com/uploadpilot/uploadpilot/web/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type webhooksHandler struct {
	webhookRepo db.WebhookRepo
}

func NewWebhooksHandler() *webhooksHandler {
	return &webhooksHandler{
		webhookRepo: db.NewWebhookRepo(),
	}
}

func (wh *webhooksHandler) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	wsID := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	webhooks, err := wh.webhookRepo.GetWebhooks(r.Context(), workspaceId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, webhooks)
}

func (wh *webhooksHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	wsID := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	webhook := &models.Webhook{}
	if err := render.DecodeJSON(r.Body, webhook); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	webhook.Enabled = true

	if err := validate.Struct(webhook); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}

	webhook, err = wh.webhookRepo.CreateWebhook(r.Context(), workspaceId, webhook)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, webhook.ID)
}

func (wh *webhooksHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	wsID := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	whID := chi.URLParam(r, "webhookId")
	webhookId, err := primitive.ObjectIDFromHex(whID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	err = wh.webhookRepo.DeleteWebhook(r.Context(), workspaceId, webhookId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)
}

func (wh *webhooksHandler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("email")
	wsID := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	whID := chi.URLParam(r, "webhookId")
	webhookId, err := primitive.ObjectIDFromHex(whID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	webhook := &models.Webhook{}
	webhook.ID = webhookId
	if err := render.DecodeJSON(r.Body, webhook); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	webhook, err = wh.webhookRepo.UpdateWebhook(r.Context(), workspaceId, webhook, userEmail)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, webhook.ID)
}

func (wh *webhooksHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	wsID := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	whID := chi.URLParam(r, "webhookId")
	webhookId, err := primitive.ObjectIDFromHex(whID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	webhook, err := wh.webhookRepo.GetWebhook(r.Context(), workspaceId, webhookId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, webhook)
}

func (wh *webhooksHandler) PatchWebhook(w http.ResponseWriter, r *http.Request) {
	wsID := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	whID := chi.URLParam(r, "webhookId")
	webhookId, err := primitive.ObjectIDFromHex(whID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	patchReq := &dto.PatchWebhookRequest{}
	if err := render.DecodeJSON(r.Body, patchReq); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	if reflect.TypeOf(patchReq.Enabled).Kind() != reflect.Bool {
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("invalid type for enabled: %T", patchReq.Enabled))
		return
	}

	err = wh.webhookRepo.PatchWebhook(r.Context(), workspaceId, webhookId, &bson.M{
		"enabled": patchReq.Enabled,
	})
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)
}
