import AuthPage from "./pages/auth";
import AppLayout from "./components/Layout/AppLayout";
import { NotFoundImage } from "./components/NotFound";
import DashboardPage from "./pages/analytics";
import EmptyLayout from "./components/Layout/EmptyLayout";
import ErrorQueryDisplay from "./pages/error";
import WorkspaceLandingPage from "./pages/workspace";
import WorkspacesLayout from "./components/Layout/WorkspacesLayout";
import UploaderPreviewPage from "./pages/getstarted";
import UploadsPage from "./pages/uploads/Uploads";
import ConfigurationPage from "./pages/configuration";
import WorkspaceUsersPage from "./pages/users";
import ToolsPage from "./pages/tools";
import ProcEditorPage from "./pages/taskeditor";
import ProcessorPage from "./pages/processors";
import EmptyAuthLayout from "./components/Layout/EmptyAuthLayout";

type Route = {
  path: string;
  layout: React.FC<{ children: React.ReactNode }>;
  element: React.ReactNode;
  children?: Route[];
};

const routes: Route[] = [
  {
    path: "/auth",
    layout: EmptyLayout,
    element: <AuthPage />,
  },
  {
    path: "/error",
    layout: EmptyLayout,
    element: <ErrorQueryDisplay />,
  },
  {
    path: "/",
    layout: WorkspacesLayout,
    element: <WorkspaceLandingPage />,
  },
  {
    path: "/workspaces/:workspaceId",
    layout: AppLayout,
    element: <UploaderPreviewPage />,
  },
  {
    path: "/workspaces/:workspaceId/uploads",
    layout: AppLayout,
    element: <UploadsPage />,
  },
  {
    path: "/workspaces/:workspaceId/configuration",
    layout: AppLayout,
    element: <ConfigurationPage />,
  },
  {
    path: "/workspaces/:workspaceId/processors",
    layout: AppLayout,
    element: <ProcessorPage />,
  },
  {
    path: "/workspaces/:workspaceId/processors/:processorId",
    layout: EmptyAuthLayout,
    element: <ProcEditorPage />,
  },
  {
    path: "/workspaces/:workspaceId/users",
    layout: AppLayout,
    element: <WorkspaceUsersPage />,
  },
  {
    path: "/workspaces/:workspaceId/tools",
    layout: AppLayout,
    element: <ToolsPage />,
  },
  {
    path: "/workspaces/:workspaceId/analytics",
    layout: AppLayout,
    element: <DashboardPage />,
  },
  {
    path: "*",
    element: <NotFoundImage />,
    layout: WorkspacesLayout,
  },
];

export default routes;
