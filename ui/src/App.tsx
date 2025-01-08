import "@mantine/core/styles.css";
import 'mantine-datatable/styles.layer.css';
import '@mantine/notifications/styles.css';
import { Routes, BrowserRouter as Router, Route } from "react-router-dom";
import { MantineProvider } from '@mantine/core';
import routes from "./Routes";
import { theme } from "./style/theme";
import { Notifications } from "@mantine/notifications";


function App() {
    return (
        <MantineProvider theme={theme}>
            <Notifications />
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
        </MantineProvider>
    );
}

export default App;
