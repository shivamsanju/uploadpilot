import "@mantine/core/styles.css";
import '@mantine/notifications/styles.css';
import { Routes, BrowserRouter as Router, Route } from "react-router-dom";
import { MantineProvider } from '@mantine/core';
import routes from "./Routes";
import { theme } from "./style/theme";
import { Notifications } from "@mantine/notifications";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";


const queryClient = new QueryClient()
function App() {

    return (
        <MantineProvider theme={theme}>
            <Notifications position="bottom-center" />
            <QueryClientProvider client={queryClient}>
                <Router>
                    <Routes>
                        {routes.map((route) => (
                            <Route
                                key={route.path}
                                path={route.path}
                                element={
                                    <route.layout>
                                        {route.element}
                                    </route.layout>
                                }
                            />
                        ))}
                    </Routes>
                </Router>
            </QueryClientProvider>
        </MantineProvider >
    );
}

export default App;
