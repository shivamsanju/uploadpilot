import SuperTokens from 'supertokens-auth-react';
import EmailPassword from 'supertokens-auth-react/recipe/emailpassword';
import Session from 'supertokens-auth-react/recipe/session';
import ThirdParty, {
  Github,
  Google,
} from 'supertokens-auth-react/recipe/thirdparty';

export const InitSupertokens = () => {
  SuperTokens.init({
    appInfo: {
      // learn more about this on https://supertokens.com/docs/references/app-info
      appName: 'Upload Pilot',
      apiDomain: 'http://localhost:8080',
      websiteDomain: 'http://localhost:3000',
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
      Session.init(),
    ],
  });
};
