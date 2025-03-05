package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/phuslu/log"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/msg"
	"github.com/uploadpilot/core/web/webutils"
)

type userHandler struct{}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) GetUserDetails(
	r *http.Request,
	params interface{},
	query interface{},
	body interface{},
) (*dto.UserDetailsResponse, int, error) {
	var us dto.UserDetailsResponse

	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userID := sessionContainer.GetUserID()
	if userID == "" {
		return nil, http.StatusUnauthorized, errors.New(msg.ErrUserInfoNotFoundInRequest)
	}

	tpusr, err := thirdparty.GetUserByID(userID)
	if err == nil && tpusr != nil {
		us.Email = tpusr.Email
	} else {
		epusr, err := emailpassword.GetUserByID(userID)
		if err == nil && epusr != nil {
			us.Email = epusr.Email
		} else {
			return nil, http.StatusUnauthorized, errors.New(msg.ErrUserInfoNotFoundInRequest)
		}
	}

	metadata, err := usermetadata.GetUserMetadata(userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to get user metadata")
		return nil, http.StatusUnauthorized, errors.New(msg.ErrInvalidUserInfoInRequest)
	}

	userAttr, err := webutils.GetUserAtrributesFromMetadata(metadata)
	if err != nil {
		log.Error().Err(err).Msg("unable to get user tenants")
		return nil, http.StatusUnauthorized, err
	}

	if userAttr != nil {
		us.Name = userAttr.Name
		us.Avatar = &userAttr.Avatar
		us.Theme = &userAttr.Theme
	}

	if us.Name == "" {
		us.Name = getNameFromEmail(us.Email)
	}

	tenants, err := webutils.GetUserTenantsFromMetadata(metadata)
	if err != nil {
		log.Error().Err(err).Msg("unable to get user tenants")
		return nil, http.StatusUnauthorized, err
	}
	tenantMap := make(map[string]string)
	for k, v := range tenants {
		tenantMap[k] = v.Name
	}
	us.Tenants = tenantMap

	us.ActiveTenant = webutils.GetActiveTenantIDFromMetadata(metadata)

	return &us, http.StatusOK, nil
}

func (h *userHandler) UpdateUserDetails(
	r *http.Request,
	params interface{},
	query interface{},
	body dto.UserAttributes,
) (string, int, error) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userID := sessionContainer.GetUserID()

	attributes := make(map[string]interface{})
	if err := mapstructure.Decode(body, &attributes); err != nil {
		return "", http.StatusUnprocessableEntity, err
	}

	if _, err := usermetadata.UpdateUserMetadata(userID, map[string]interface{}{
		dto.UserAttributesKey: attributes,
	}); err != nil {
		return "", http.StatusBadRequest, err
	}

	return "OK", http.StatusOK, nil
}

func getNameFromEmail(email string) string {
	atPos := strings.Index(email, "@")
	if atPos == -1 {
		return email
	}
	return strings.ToUpper(email[:1]) + email[1:atPos]
}
