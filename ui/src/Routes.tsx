import AuthPage from "./pages/Auth";
import AppLayout from "./components/Layout/AppLayout"
import { NotFoundImage } from "./components/NotFound";
import UploaderPage from "./pages/Uploader";
import CreateNewUploaderPage from "./pages/Uploader/New";
import UploaderDetailsPage from "./pages/Uploader/Details";
import ProfilePage from "./pages/Profile";
import DashboardPage from "./pages/Dashboard";
import NewConnectorsPage from "./pages/Storage/Connectors/New";
import ConnectorsPage from "./pages/Storage/Connectors";
import EmptyLayout from "./components/Layout/EmptyLayout";
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
        path: "/uploaders",
        layout: AppLayout,
        element: <UploaderPage />
    },
    {
        path: "/uploaders/new",
        layout: AppLayout,
        element: <CreateNewUploaderPage />
    },
    {
        path: "/uploaders/:uploaderId",
        layout: AppLayout,
        element: <UploaderDetailsPage />
    },
    {
        path: "/uploaders/:uploaderId/:tabValue",
        layout: AppLayout,
        element: <UploaderDetailsPage />
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