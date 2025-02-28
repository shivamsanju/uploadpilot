import { useNavigate } from 'react-router-dom';
import { useGetSession } from '../../apis/user';
import { ErrorCard } from '../ErrorCard/ErrorCard';
import { AppLoader } from '../Loader/AppLoader';

const SesionManager: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const navigate = useNavigate();

  const { session, isPending, error } = useGetSession();

  if (isPending) {
    return <AppLoader h="100vh" />;
  }

  if (error) {
    return <ErrorCard title={error.name} message={error.message} h="100vh" />;
  }

  if (!session) {
    navigate('/auth');
    return <></>;
  }

  return <>{children}</>;
};

export default SesionManager;
