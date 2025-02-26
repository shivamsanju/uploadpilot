import AppLayout from './components/Layout/AppLayout';
import EmptyLayout from './components/Layout/EmptyLayout';
import ProcessorLayout from './components/Layout/ProcessorLayout';
import WorkspacesLayout from './components/Layout/WorkspaceLayout';
import { NotFoundImage } from './components/NotFound';
import DashboardPage from './pages/analytics';
import ApiKeyPage from './pages/apikeys';
import AuthPage from './pages/auth';
import BillingsPage from './pages/billing';
import ConfigurationPage from './pages/configuration';
import ErrorQueryDisplay from './pages/error';
import { GetStartedPage } from './pages/getstarted';
import ProcessorPage from './pages/processors';
import NewprocessorPage from './pages/processors/Add';
import ProcessorSettingsPage from './pages/processors/settings';
import ToolsPage from './pages/tools';
import UploadsPage from './pages/uploads';
import WorkspaceUsersPage from './pages/users';
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
    path: '/error',
    layout: EmptyLayout,
    element: <ErrorQueryDisplay />,
  },
  {
    path: '/',
    layout: WorkspacesLayout,
    element: <WorkspaceLandingPage />,
  },
  {
    path: '/billing',
    layout: WorkspacesLayout,
    element: <BillingsPage />,
  },
  {
    path: '/workspace/:workspaceId',
    layout: AppLayout,
    element: <GetStartedPage />,
  },
  {
    path: '/workspace/:workspaceId/uploads',
    layout: AppLayout,
    element: <UploadsPage />,
  },
  {
    path: '/workspace/:workspaceId/configuration',
    layout: AppLayout,
    element: <ConfigurationPage />,
  },
  {
    path: '/workspace/:workspaceId/processors',
    layout: AppLayout,
    element: <ProcessorPage />,
  },
  {
    path: '/workspace/:workspaceId/processors/new',
    layout: AppLayout,
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
    path: '/workspace/:workspaceId/users',
    layout: AppLayout,
    element: <WorkspaceUsersPage />,
  },
  {
    path: '/workspace/:workspaceId/apikeys',
    layout: AppLayout,
    element: <ApiKeyPage />,
  },
  {
    path: '/workspace/:workspaceId/tools',
    layout: AppLayout,
    element: <ToolsPage />,
  },
  {
    path: '/workspace/:workspaceId/analytics',
    layout: AppLayout,
    element: <DashboardPage />,
  },
  {
    path: '*',
    element: <NotFoundImage />,
    layout: WorkspacesLayout,
  },
];

export default routes;
