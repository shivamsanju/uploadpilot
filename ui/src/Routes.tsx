import EmptyLayout from './components/Layout/EmptyLayout';
import HeaderAuthNoTenancyLayout from './components/Layout/HeaderAuthNoTenancyLayout';
import HomeLayout from './components/Layout/HomeLayout';
import ProcessorLayout from './components/Layout/ProcessorLayout';
import WorkspaceLayout from './components/Layout/WorkspaceLayout';
import { NotFoundImage } from './components/NotFound';
import DashboardPage from './pages/analytics';
import ApiKeyPage from './pages/apikeys';
import CreateApiKeyPage from './pages/apikeys/add';
import AuthPage from './pages/auth';
import SocialAuthCallbackHandlerPage from './pages/auth/callback';
import BillingsPage from './pages/billing';
import ConfigurationPage from './pages/configuration';
import ErrorQueryDisplay from './pages/error';
import { GetStartedPage } from './pages/getstarted';
import ProcessorPage from './pages/processors';
import NewprocessorPage from './pages/processors/Add';
import ProcessorSettingsPage from './pages/processors/settings';
import { TenantRegistrationPage } from './pages/tenancy/registration';
import TenantSelectionPage from './pages/tenancy/selection';
import UploadsPage from './pages/uploads';
import WorkflowBuilderPage from './pages/wflowbuilder';
import WorkflowRunsPage from './pages/workflowruns';
import WorkspaceLandingPage from './pages/workspace';

type Route = {
  path: string;
  layout: React.FC<{ children: React.ReactNode }>;
  element: React.ReactNode;
  children?: Route[];
};

const routes: Route[] = [
  {
    path: '/auth',
    layout: EmptyLayout,
    element: <AuthPage />,
  },
  {
    path: '/auth/callback/social',
    layout: EmptyLayout,
    element: <SocialAuthCallbackHandlerPage />,
  },
  {
    path: '/register-tenant',
    layout: HeaderAuthNoTenancyLayout,
    element: <TenantRegistrationPage />,
  },
  {
    path: '/tenants',
    layout: HeaderAuthNoTenancyLayout,
    element: <TenantSelectionPage />,
  },
  {
    path: '/error',
    layout: EmptyLayout,
    element: <ErrorQueryDisplay />,
  },
  {
    path: '/',
    layout: HomeLayout,
    element: <WorkspaceLandingPage />,
  },
  {
    path: '/billing',
    layout: HomeLayout,
    element: <BillingsPage />,
  },
  {
    path: '/api-keys',
    layout: HomeLayout,
    element: <ApiKeyPage />,
  },
  {
    path: '/api-keys/new',
    layout: HomeLayout,
    element: <CreateApiKeyPage />,
  },
  {
    path: '/workspace/:workspaceId',
    layout: WorkspaceLayout,
    element: <GetStartedPage />,
  },
  {
    path: '/workspace/:workspaceId/uploads',
    layout: WorkspaceLayout,
    element: <UploadsPage />,
  },
  {
    path: '/workspace/:workspaceId/configuration',
    layout: WorkspaceLayout,
    element: <ConfigurationPage />,
  },
  {
    path: '/workspace/:workspaceId/processors',
    layout: WorkspaceLayout,
    element: <ProcessorPage />,
  },
  {
    path: '/workspace/:workspaceId/processors/new',
    layout: WorkspaceLayout,
    element: <NewprocessorPage />,
  },
  {
    path: '/workspace/:workspaceId/processors/:processorId/workflow',
    layout: ProcessorLayout,
    element: <WorkflowBuilderPage />,
  },
  {
    path: '/workspace/:workspaceId/processors/:processorId/runs',
    layout: ProcessorLayout,
    element: <WorkflowRunsPage />,
  },
  {
    path: '/workspace/:workspaceId/processors/:processorId/settings',
    layout: ProcessorLayout,
    element: <ProcessorSettingsPage />,
  },

  {
    path: '/workspace/:workspaceId/analytics',
    layout: WorkspaceLayout,
    element: <DashboardPage />,
  },
  {
    path: '*',
    element: <NotFoundImage />,
    layout: HomeLayout,
  },
];

export default routes;
