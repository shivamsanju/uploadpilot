import { AppLoader } from "../Loader/AppLoader"
import { ErrorCard } from "../ErrorCard/ErrorCard"

type ErrorLoadingWrapperProps = {
    children: React.ReactNode
    isPending: boolean
    error: Error | null
    h?: string
    w?: string
}
export const ErrorLoadingWrapper: React.FC<ErrorLoadingWrapperProps> = ({ children, isPending, error, h = "50vh", w }) => {
    return (<>
        {isPending ?
            <AppLoader h={h} /> : error ?
                <ErrorCard title={error.name} message={error.message} h={h} /> : children}
    </>)
}