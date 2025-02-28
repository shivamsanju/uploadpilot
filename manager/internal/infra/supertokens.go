package infra

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Provider struct {
	Key          string
	ClientID     string
	ClientSecret string
}

type SuperTokensOptions struct {
	ConnectionURI   string
	APIKey          string
	AppName         string
	APIBasePath     string
	WebsiteBasePath string
	APIDomain       string
	WebsiteDomain   string
	Providers       []Provider
}

func InitSuperTokens(opts *SuperTokensOptions) error {
	apiBasePath := opts.APIBasePath
	websiteBasePath := opts.WebsiteBasePath

	providers := make([]tpmodels.ProviderInput, len(opts.Providers))
	for i, provider := range opts.Providers {
		providers[i] = tpmodels.ProviderInput{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId: provider.Key,
				Clients: []tpmodels.ProviderClientConfig{
					{
						ClientID:     provider.ClientID,
						ClientSecret: provider.ClientSecret,
					},
				},
			},
		}
	}

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: opts.ConnectionURI,
			APIKey:        opts.APIKey,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         opts.AppName,
			APIDomain:       opts.APIDomain,
			WebsiteDomain:   opts.WebsiteDomain,
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			dashboard.Init(nil),
			usermetadata.Init(nil),
			emailpassword.Init(nil),
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: providers,
				},
				// Override: &tpmodels.OverrideStruct{
				// 	Functions: func(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
				// 		ogSignInUp := *originalImplementation.SignInUp

				// 		(*originalImplementation.SignInUp) = func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens map[string]interface{}, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId string, userContext *map[string]interface{}) (tpmodels.SignInUpResponse, error) {
				// 			existingUsers, err := thirdparty.GetUsersByEmail(tenantId, email)
				// 			if err != nil {
				// 				return tpmodels.SignInUpResponse{}, err
				// 			}

				// 			emailPasswordExistingUser, err := passwordless.GetUserByEmail(tenantId, email)
				// 			if err != nil {
				// 				return tpmodels.SignInUpResponse{}, err
				// 			}

				// 			if emailPasswordExistingUser != nil {
				// 				return tpmodels.SignInUpResponse{}, errors.New("cannot sign up as email already exists")
				// 			}

				// 			if len(existingUsers) == 0 {
				// 				// this means this email is new so we allow sign up
				// 				return ogSignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, tenantId, userContext)
				// 			}

				// 			isSignIn := false
				// 			for _, user := range existingUsers {
				// 				if user.ThirdParty.ID == thirdPartyID && user.ThirdParty.UserID == thirdPartyUserID {
				// 					// this means we are trying to sign in with the same social login. So we allow it
				// 					isSignIn = true
				// 				}
				// 			}
				// 			if isSignIn {
				// 				return ogSignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, tenantId, userContext)
				// 			}
				// 			return tpmodels.SignInUpResponse{}, errors.New("cannot sign up as email already exists")
				// 		}

				// 		return originalImplementation
				// 	},

				// 	APIs: func(originalImplementation tpmodels.APIInterface) tpmodels.APIInterface {
				// 		originalSignInUpPOST := *originalImplementation.SignInUpPOST

				// 		(*originalImplementation.SignInUpPOST) = func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, tenantId string, options tpmodels.APIOptions, userContext *map[string]interface{}) (tpmodels.SignInUpPOSTResponse, error) {

				// 			resp, err := originalSignInUpPOST(provider, input, tenantId, options, userContext)

				// 			if err != nil && err.Error() == "cannot sign up as email already exists" {
				// 				// this error was thrown from our function override above.
				// 				// so we send a useful message to the user
				// 				return tpmodels.SignInUpPOSTResponse{
				// 					GeneralError: &supertokens.GeneralErrorResponse{
				// 						Message: "Seems like you already have an account with another method. Please use that instead.",
				// 					},
				// 				}, nil
				// 			}

				// 			return resp, err
				// 		}

				// 		return originalImplementation
				// 	},
				// },
			}),
			session.Init(nil),
		},
	})

	if err != nil {
		return err
	}

	return nil
}
