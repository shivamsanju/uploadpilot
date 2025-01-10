import AuthPage from "./pages/Auth";
import AppLayout from "./components/Layout/AppLayout"
import { NotFoundImage } from "./components/NotFound";
import WorkflowsPage from "./pages/Workflow";
import NewWorkflowPage from "./pages/Workflow/New";
import WorkflowDetailsPage from "./pages/Workflow/Details";
import ProfilePage from "./pages/Profile";
import DashboardPage from "./pages/Dashboard";
import NewConnectorsPage from "./pages/Storage/Connectors/New";
import ConnectorsPage from "./pages/Storage/Connectors";
import EmptyLayout from "./components/Layout/EmptyLayout";
import ImportPoliciesPage from "./pages/ImportPolicy";
import NewImportPolicyPage from "./pages/ImportPolicy/New";
import ImportPolicyDetailsPage from "./pages/ImportPolicy/Details";
import DataStoresPage from "./pages/Storage/DataStores";

type Route = {
    path: string
    layout: React.FC<{ children: React.ReactNode; }>
    element: React.ReactNode
}

const routes: Route[] = [
    {
        path: "/auth",
        layout: EmptyLayout,
        element: <AuthPage />
    },
    {
        path: "/",
        layout: AppLayout,
        element: <DashboardPage />
    },
    {
        path: "/workflows",
        layout: AppLayout,
        element: <WorkflowsPage />
    },
    {
        path: "/workflows/new",
        layout: AppLayout,
        element: <NewWorkflowPage />
    },
    {
        path: "/workflows/:workflowId",
        layout: AppLayout,
        element: <WorkflowDetailsPage />
    },
    {
        path: "/importPolicies",
        layout: AppLayout,
        element: <ImportPoliciesPage />
    },
    {
        path: "/importPolicies/new",
        layout: AppLayout,
        element: <NewImportPolicyPage />
    },
    {
        path: "/importPolicies/:importPolicyId",
        layout: AppLayout,
        element: <ImportPolicyDetailsPage />
    },
    {
        path: "/storage/datastores",
        layout: AppLayout,
        element: <DataStoresPage />
    },
    {
        path: "/storage/connectors",
        layout: AppLayout,
        element: <ConnectorsPage />
    },
    {
        path: "/storage/connectors/new",
        layout: AppLayout,
        element: <NewConnectorsPage />
    },
    {
        path: "/profile",
        layout: AppLayout,
        element: <ProfilePage />
    },
    {
        path: "*",
        element: <NotFoundImage />,
        layout: AppLayout
    },
]

export default routes