import {
  doesEmailExist,
  signIn,
  signUp,
} from 'supertokens-web-js/recipe/emailpassword';
import Session from 'supertokens-web-js/recipe/session';
import {
  getAuthorisationURLWithQueryParamsAndSetState,
  signInAndUp,
} from 'supertokens-web-js/recipe/thirdparty';
import { getWebsiteDomain } from '../utils/config';

export const handleEPSignup = async (
  name: string,
  email: string,
  password: string,
) => {
  try {
    let response = await signUp({
      formFields: [
        {
          id: 'email',
          value: email,
        },
        {
          id: 'password',
          value: password,
        },
      ],
    });

    if (response.status === 'FIELD_ERROR') {
      // one of the input formFields failed validation
      response.formFields.forEach(formField => {
        if (formField.id === 'email') {
          // Email validation failed (for example incorrect email syntax),
          // or the email is not unique.
          window.alert(formField.error);
        } else if (formField.id === 'password') {
          // Password validation failed.
          // Maybe it didn't match the password strength
          window.alert(formField.error);
        }
      });
    } else if (response.status === 'SIGN_UP_NOT_ALLOWED') {
      // the reason string is a user friendly message
      // about what went wrong. It can also contain a support code which users
      // can tell you so you know why their sign up was not allowed.
      window.alert(response.reason);
    } else {
      // sign up successful. The session tokens are automatically handled by
      // the frontend SDK.
      window.location.href = '/';
    }
  } catch (err: any) {
    if (err.isSuperTokensGeneralError === true) {
      // this may be a custom error message sent from the API by you.
      window.alert(err.message);
    } else {
      window.alert('Oops! Something went wrong.');
    }
  }
};

export const checkEmailExists = async (email: string) => {
  try {
    let response = await doesEmailExist({
      email,
    });

    if (response.doesExist) {
      window.alert('Email already exists. Please sign in instead');
    }
  } catch (err: any) {
    if (err.isSuperTokensGeneralError === true) {
      // this may be a custom error message sent from the API by you.
      window.alert(err.message);
    } else {
      window.alert('Oops! Something went wrong.');
    }
  }
};

export const handleEPSignIn = async (email: string, password: string) => {
  try {
    let response = await signIn({
      formFields: [
        {
          id: 'email',
          value: email,
        },
        {
          id: 'password',
          value: password,
        },
      ],
    });

    if (response.status === 'FIELD_ERROR') {
      response.formFields.forEach(formField => {
        if (formField.id === 'email') {
          // Email validation failed (for example incorrect email syntax).
          window.alert(formField.error);
        }
      });
    } else if (response.status === 'WRONG_CREDENTIALS_ERROR') {
      window.alert('Email password combination is incorrect.');
    } else if (response.status === 'SIGN_IN_NOT_ALLOWED') {
      // the reason string is a user friendly message
      // about what went wrong. It can also contain a support code which users
      // can tell you so you know why their sign in was not allowed.
      window.alert(response.reason);
    } else {
      // sign in successful. The session tokens are automatically handled by
      // the frontend SDK.
      window.location.href = '/';
    }
  } catch (err: any) {
    if (err.isSuperTokensGeneralError === true) {
      // this may be a custom error message sent from the API by you.
      window.alert(err.message);
    } else {
      window.alert('Oops! Something went wrong.');
    }
  }
};

export const doesSessionExist = async () => {
  const exists = await Session.doesSessionExist();
  return exists;
};

export const handleSignout = async () => {
  await Session.signOut();
  window.location.href = '/auth';
};

export const handleSocialLogin = async (provider: string) => {
  try {
    const authUrl = await getAuthorisationURLWithQueryParamsAndSetState({
      thirdPartyId: provider,

      // This is where Google should redirect the user back after login or error.
      // This URL goes on the Google's dashboard as well.
      frontendRedirectURI: getWebsiteDomain() + `/auth/callback/${provider}`,
    });

    /*
        Example value of authUrl: https://accounts.google.com/o/oauth2/v2/auth/oauthchooseaccount?scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&access_type=offline&include_granted_scopes=true&response_type=code&client_id=1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com&state=5a489996a28cafc83ddff&redirect_uri=https%3A%2F%2Fsupertokens.io%2Fdev%2Foauth%2Fredirect-to-app&flowName=GeneralOAuthFlow
        */

    // we redirect the user to google for auth.
    window.location.assign(authUrl);
  } catch (err: any) {
    console.log(err);
    if (err.isSuperTokensGeneralError === true) {
      // this may be a custom error message sent from the API by you.
      window.alert(err.message);
    } else {
      window.alert('Oops! Something went wrong.');
    }
  }
};

export const handleSocialLoginCallback = async (provider?: string) => {
  try {
    const response = await signInAndUp();
    console.log({ response });
    if (response.status === 'OK') {
      if (
        response.createdNewRecipeUser &&
        response.user.loginMethods.length === 1
      ) {
        // sign up successful
      } else {
        // sign in successful
      }
      window.location.assign('/');
    } else if (response.status === 'SIGN_IN_UP_NOT_ALLOWED') {
      // the reason string is a user friendly message
      // about what went wrong. It can also contain a support code which users
      // can tell you so you know why their sign in / up was not allowed.
      window.alert(response.reason);
      window.location.assign('/auth'); // redirect back to login page
    } else {
      // SuperTokens requires that the third party provider
      // gives an email for the user. If that's not the case, sign up / in
      // will fail.

      // As a hack to solve this, you can override the backend functions to create a fake email for the user.

      window.alert(
        `No email provided by ${provider} login. Please use another form of login`,
      );
      window.location.assign('/auth'); // redirect back to login page
    }
  } catch (err: any) {
    console.log(err);
    if (err.isSuperTokensGeneralError === true) {
      // this may be a custom error message sent from the API by you.
      window.alert(err.message);
    } else {
      window.alert('Oops! Something went wrong.');
    }
    window.location.assign('/auth'); // redirect back to login page
  }
};
