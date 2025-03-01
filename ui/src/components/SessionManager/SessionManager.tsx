import { useLocation, useNavigate } from 'react-router-dom';
import { useGetSession } from '../../apis/user';
import { ErrorCard } from '../ErrorCard/ErrorCard';
import { AppLoader } from '../Loader/AppLoader';

export const SessionManager: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const { session, isPending, error } = useGetSession();

  if (isPending) {
    return <AppLoader h="100vh" />;
  }

  if (error) {
    return <ErrorCard title={error.name} message={error.message} h="100vh" />;
  }

  if (!session && !pathname.startsWith('/auth')) {
    navigate('/auth');
    return <></>;
  }

  if (session && pathname.startsWith('/auth')) {
    navigate('/');
    return <></>;
  }

  return <>{children}</>;
};
