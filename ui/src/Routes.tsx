import AuthPage from "./pages/Auth";
import AppLayout from "./components/Layout/AppLayout"
import { NotFoundImage } from "./components/NotFound";
import ProfilePage from "./pages/Profile";
import DashboardPage from "./pages/Dashboard";
import EmptyLayout from "./components/Layout/EmptyLayout";
import ErrorQueryDisplay from "./pages/Error";
import WorkspaceLandingPage from "./pages/Workspace";
import WorkspacesLayout from "./components/Layout/WorkspacesLayout";
import UploaderPreviewPage from "./pages/Workspace/Preview";
import ImportsPage from "./pages/Workspace/Imports/Imports";
import ConfigurationPage from "./pages/Workspace/Configuration";
import WorkspaceUsersPage from "./pages/Workspace/Users";
import WorkspaceHooksPage from "./pages/Workspace/Hooks";
import WorkspaceWebhooksPage from "./pages/Workspace/Webhooks";

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
        path: "/workspaces/:workspaceId/imports",
        layout: AppLayout,
        element: <ImportsPage />
    },
    {
        path: "/workspaces/:workspaceId/configuration",
        layout: AppLayout,
        element: <ConfigurationPage />
    },
    {
        path: "/workspaces/:workspaceId/users",
        layout: AppLayout,
        element: <WorkspaceUsersPage />
    },
    {
        path: "/workspaces/:workspaceId/hooks",
        layout: AppLayout,
        element: <WorkspaceHooksPage />
    },
    {
        path: "/workspaces/:workspaceId/webhooks",
        layout: AppLayout,
        element: <WorkspaceWebhooksPage />
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