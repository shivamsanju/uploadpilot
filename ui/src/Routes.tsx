import AuthPage from "./pages/Auth";
import AppLayout from "./components/Layout/AppLayout"
import { NotFoundImage } from "./components/NotFound";
import ProfilePage from "./pages/Profile";
import DashboardPage from "./pages/Dashboard";
import EmptyLayout from "./components/Layout/EmptyLayout";
import ErrorQueryDisplay from "./pages/Error";
import WorkspaceLandingPage from "./pages/Workspace";
import WorkspacesLayout from "./components/Layout/WorkspacesLayout";
import UploaderPreviewPage from "./pages/GetStarted";
import UploadsPage from "./pages/Uploads/Uploads";
import ConfigurationPage from "./pages/Configuration";
import WorkspaceUsersPage from "./pages/Users";
import ToolsPage from "./pages/Tools";
import WebhooksPage from "./pages/Webhooks";
import ProcEditorPage from "./pages/ProcEditor";
import ProcessorPage from "./pages/Processors";
import EmptyAuthLayout from "./components/Layout/EmptyAuthLayout";

type Route = {
    path: string
    layout: React.FC<{ children: React.ReactNode; }>
    element: React.ReactNode
    children?: Route[]
}

const routes: Route[] = [
    {
        path: "/auth",
        layout: EmptyLayout,
        element: <AuthPage />
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
        path: "/profile",
        layout: WorkspacesLayout,
        element: <ProfilePage />
    },
    {
        path: "/workspaces/:workspaceId",
        layout: AppLayout,
        element: <UploaderPreviewPage />
    },
    {
        path: "/workspaces/:workspaceId/uploads",
        layout: AppLayout,
        element: <UploadsPage />
    },
    {
        path: "/workspaces/:workspaceId/configuration",
        layout: AppLayout,
        element: <ConfigurationPage />
    },
    {
        path: "/workspaces/:workspaceId/processors",
        layout: AppLayout,
        element: <ProcessorPage />
    },
    {
        path: "/workspaces/:workspaceId/processors/:processorId",
        layout: EmptyAuthLayout,
        element: <ProcEditorPage />
    },
    {
        path: "/workspaces/:workspaceId/users",
        layout: AppLayout,
        element: <WorkspaceUsersPage />
    },
    {
        path: "/workspaces/:workspaceId/tools",
        layout: AppLayout,
        element: <ToolsPage />
    },
    {
        path: "/workspaces/:workspaceId/webhooks",
        layout: AppLayout,
        element: <WebhooksPage />
    },
    {
        path: "/workspaces/:workspaceId/analytics",
        layout: AppLayout,
        element: <DashboardPage />
    },
    {
        path: "*",
        element: <NotFoundImage />,
        layout: WorkspacesLayout
    },
]

export default routes