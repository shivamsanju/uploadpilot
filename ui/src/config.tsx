import ThirdParty, { Google, Github, Apple, Twitter } from "supertokens-auth-react/recipe/thirdparty";
import { ThirdPartyPreBuiltUI } from "supertokens-auth-react/recipe/thirdparty/prebuiltui";
import Session from "supertokens-auth-react/recipe/session";

export function getApiDomain() {
    const apiUrl = process.env.REACT_APP_BACKEND_URL || `http://localhost:8000`;
    return apiUrl;
}

export function getWebsiteDomain() {
    const websiteUrl = process.env.REACT_APP_WEBSITE_URL || `http://localhost:3000`;
    return websiteUrl;
}

export function getAppName() {
    return process.env.REACT_APP_APP_NAME || "Code Monk";
}

export const SuperTokensConfig = {
    appInfo: {
        appName: getAppName(),
        apiDomain: getApiDomain(),
        websiteDomain: getWebsiteDomain(),
        apiBasePath: "/auth",
        websiteBasePath: "/auth"
    },
    recipeList: [
        ThirdParty.init({
            signInAndUpFeature: {
                providers: [Github.init(), Google.init(), Apple.init(), Twitter.init()],
            },
        }),
        Session.init(),
    ],
};

export const recipeDetails = {
    docsLink: "https://supertokens.com/docs/thirdpartyemailpassword/introduction",
};

export const PreBuiltUIList = [ThirdPartyPreBuiltUI];

export const ComponentWrapper = (props: { children: JSX.Element }): JSX.Element => {
    return props.children;
};
