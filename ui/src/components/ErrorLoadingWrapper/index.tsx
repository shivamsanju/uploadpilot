import { ErrorCard } from '../ErrorCard/ErrorCard';
import { AppLoader } from '../Loader/AppLoader';

type ErrorLoadingWrapperProps = {
  children: React.ReactNode;
  isPending: boolean;
  error: Error | null;
  h?: string;
  w?: string;
};
export const ErrorLoadingWrapper: React.FC<ErrorLoadingWrapperProps> = ({
  children,
  isPending,
  error,
  h = '70vh',
  w,
}) => {
  return (
    <>
      {isPending ? (
        <AppLoader h={h} />
      ) : error ? (
        <ErrorCard title={error.name} message={error.message} h={h} />
      ) : (
        children
      )}
    </>
  );
};
