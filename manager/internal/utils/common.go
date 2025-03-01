package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/msg"
)

func GetSessionFromCtx(ctx context.Context) (*dto.Session, error) {
	value := ctx.Value(dto.SessionCtxKey)

	if value == nil {
		return nil, errors.New(msg.ErrUserInfoNotFoundInRequest)
	}

	userInfo, ok := value.(dto.Session)
	if !ok {
		return nil, errors.New(msg.ErrInvalidUserInfoInRequest)
	}

	return &userInfo, nil
}

func GetStatusLabel(status int) string {
	switch {
	case status >= 100 && status < 300:
		return fmt.Sprintf("%d OK", status)
	case status >= 300 && status < 400:
		return fmt.Sprintf("%d Redirect", status)
	case status >= 400 && status < 500:
		return fmt.Sprintf("%d Client Error", status)
	case status >= 500:
		return fmt.Sprintf("%d Server Error", status)
	default:
		return fmt.Sprintf("%d Unknown", status)
	}
}

func GetAPIKeyFromReq(r *http.Request) string {
	return r.Header.Get("X-Api-Key")
}

func GetTenantIDFromReq(r *http.Request) string {
	return r.Header.Get("X-Tenant-Id")
}

func GetUserTenantsFromMetadata(metadata map[string]interface{}) (map[string]dto.TenantMetadata, error) {
	userTenants := make(map[string]dto.TenantMetadata)
	ut, ok := metadata[dto.UserMetadataTenantKey]
	if ok {
		if err := mapstructure.Decode(ut, &userTenants); err != nil {
			return nil, errors.New(msg.ErrInvalidTenantIDInRequest)
		}
	}

	return userTenants, nil
}

func GetUserAtrributesFromMetadata(metadata map[string]interface{}) (*dto.UserAttributes, error) {
	var userAttributes dto.UserAttributes
	ut, ok := metadata[dto.UserAttributesKey]
	if ok {
		if err := mapstructure.Decode(ut, &userAttributes); err != nil {
			return nil, errors.New(msg.ErrInvalidTenantIDInRequest)
		}
	}

	return &userAttributes, nil
}

func GetActiveTenantIDFromMetadata(metadata map[string]interface{}) *string {
	activeTenantID, ok := metadata[dto.ActiveTenantIDKey].(string)
	if !ok {
		return nil
	}

	return &activeTenantID
}
