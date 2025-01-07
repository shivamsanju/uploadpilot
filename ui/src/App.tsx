import "@mantine/core/styles.css";
import '@mantine/notifications/styles.css';
import SuperTokens, { SuperTokensWrapper } from "supertokens-auth-react";
import { getSuperTokensRoutesForReactRouterDom } from "supertokens-auth-react/ui";
import { SessionAuth } from "supertokens-auth-react/recipe/session";
import { Routes, BrowserRouter as Router, Route } from "react-router-dom";
import { PreBuiltUIList, SuperTokensConfig, ComponentWrapper } from "./config";
import { MantineProvider } from '@mantine/core';
import routes from "./Routes";
import { theme } from "./style/theme";
import { Notifications } from "@mantine/notifications";

SuperTokens.init(SuperTokensConfig);

function App() {
    return (
        <SuperTokensWrapper>
            <ComponentWrapper>
                <MantineProvider theme={theme}>
                    <Notifications />
                    <Router>
                        <Routes>
                            {getSuperTokensRoutesForReactRouterDom(require("react-router-dom"), PreBuiltUIList)}
                            {routes.map((route) => (
                                <Route
                                    key={route.path}
                                    path={route.path}
                                    element={
                                        <SessionAuth>
                                            <route.layout>
                                                {route.element}
                                            </route.layout>
                                        </SessionAuth>
                                    }
                                />
                            ))}
                        </Routes>
                    </Router>
                </MantineProvider>
            </ComponentWrapper>
        </SuperTokensWrapper>
    );
}

export default App;
