import CodeBasePage from "./app/codebases";
import CodeMapPage from "./app/codebases/details";
import HomeLayout from "./app/layouts"
import { NotFoundImage } from "./components/NotFound";
import ProfilePage from "./components/Profile";

type Route = {
    path: string
    layout: React.FC<{ children: React.ReactNode; }>
    element: React.ReactNode
}

const routes: Route[] = [
    {
        path: "/",
        layout: HomeLayout,
        element: <CodeBasePage />
    },
    {
        path: "/codebases/:codebaseId",
        layout: HomeLayout,
        element: <CodeMapPage />
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