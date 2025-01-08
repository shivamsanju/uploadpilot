import WorkflowsPage from "./pages/workflows";
import NewWorkflowPage from "./pages/workflows/newWorkflow";
import WorkflowDetailsPage from "./pages/workflowDetails";
import ProfilePage from "./pages/profile";
import HomeLayout from "./components/Layout"
import { NotFoundImage } from "./components/NotFound";
import DashboardPage from "./pages/dashboard";
import NewConnectorsPage from "./pages/storage/newConnector";
import ConnectorsPage from "./pages/storage/connectors";

type Route = {
    path: string
    layout: React.FC<{ children: React.ReactNode; }>
    element: React.ReactNode
}

const routes: Route[] = [
    {
        path: "/",
        layout: HomeLayout,
        element: <DashboardPage />
    },
    {
        path: "/workflows",
        layout: HomeLayout,
        element: <WorkflowsPage />
    },
    {
        path: "/workflows/new",
        layout: HomeLayout,
        element: <NewWorkflowPage />
    },
    {
        path: "/storage/units",
        layout: HomeLayout,
        element: <WorkflowsPage />
    },
    {
        path: "/storage/connectors",
        layout: HomeLayout,
        element: <ConnectorsPage />
    },
    {
        path: "/storage/connectors/new",
        layout: HomeLayout,
        element: <NewConnectorsPage />
    },
    {
        path: "/workflows/:workflowId",
        layout: HomeLayout,
        element: <WorkflowDetailsPage />
    },
    {
        path: "/profile",
        layout: HomeLayout,
        element: <ProfilePage />
    },
    {
        path: "*",
        element: <NotFoundImage />,
        layout: HomeLayout
    },
]

export default routes