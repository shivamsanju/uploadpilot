import SuperTokens from 'supertokens-auth-react';
import EmailPassword from 'supertokens-auth-react/recipe/emailpassword';
import Session from 'supertokens-auth-react/recipe/session';
import ThirdParty, {
  Github,
  Google,
} from 'supertokens-auth-react/recipe/thirdparty';
import { getApiDomain, getWebsiteDomain } from './config';

export const InitSupertokens = () => {
  SuperTokens.init({
    appInfo: {
      // learn more about this on https://supertokens.com/docs/references/app-info
      appName: 'Upload Pilot',
      apiDomain: getApiDomain(),
      websiteDomain: getWebsiteDomain(),
      apiBasePath: '/auth',
      websiteBasePath: '/auth',
    },
    recipeList: [
      EmailPassword.init(),
      ThirdParty.init({
        signInAndUpFeature: {
          providers: [Google.init(), Github.init()],
        },
      }),
      Session.init({
        tokenTransferMethod: 'header',
      }),
    ],
  });
};
