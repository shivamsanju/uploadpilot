import '@mantine/core/styles.css';
import '@mantine/dates/styles.css';
import '@mantine/notifications/styles.css';

import { MantineProvider } from '@mantine/core';
import { ModalsProvider } from '@mantine/modals';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { BreadcrumbsProvider } from './context/BreadcrumbContext';
import routes from './Routes';
import { myAppTheme } from './style/theme';
import { InitSupertokens } from './utils/supertokens';
const queryClient = new QueryClient();

InitSupertokens();

function App() {
  return (
    <MantineProvider theme={myAppTheme} forceColorScheme="dark">
      <ModalsProvider>
        <BreadcrumbsProvider>
          <Notifications position="bottom-right" transitionDuration={500} />
          <QueryClientProvider client={queryClient}>
            <Router>
              <Routes>
                {routes.map(route => (
                  <Route
                    key={route.path}
                    path={route.path}
                    element={<route.layout>{route.element}</route.layout>}
                  />
                ))}
              </Routes>
            </Router>
          </QueryClientProvider>
        </BreadcrumbsProvider>
      </ModalsProvider>
    </MantineProvider>
  );
}

export default App;
