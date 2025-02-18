import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";
import { Routes, BrowserRouter as Router, Route } from "react-router-dom";
import { MantineProvider } from "@mantine/core";
import routes from "./Routes";
import { myAppTheme } from "./style/theme";
import { Notifications } from "@mantine/notifications";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ModalsProvider } from "@mantine/modals";
const queryClient = new QueryClient();
function App() {
  return (
    <MantineProvider theme={myAppTheme}>
      <ModalsProvider>
        <Notifications position="bottom-right" transitionDuration={500} />
        <QueryClientProvider client={queryClient}>
          <Router>
            <Routes>
              {routes.map((route) => (
                <Route
                  key={route.path}
                  path={route.path}
                  element={<route.layout>{route.element}</route.layout>}
                />
              ))}
            </Routes>
          </Router>
        </QueryClientProvider>
      </ModalsProvider>
    </MantineProvider>
  );
}

export default App;
