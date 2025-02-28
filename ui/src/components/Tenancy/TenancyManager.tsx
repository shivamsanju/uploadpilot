import { useNavigate } from 'react-router-dom';
import { useGetUserDetails } from '../../apis/user';
import { TENANT_ID_KEY } from '../../constants/tenancy';
import { AppLoader } from '../Loader/AppLoader';
import { TenantRegistrationForm } from './RegistrationForm';

const TenancyManager: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const navigate = useNavigate();

  const { user, isPending, error } = useGetUserDetails();

  if (isPending) {
    return <AppLoader h="100vh" />;
  }

  if (error) {
    navigate('/auth');
  }

  if (!user?.tenants || Object.keys(user.tenants).length === 0) {
    return <TenantRegistrationForm />;
  }

  if (!user?.activeTenant) {
    navigate('/tenants');
    return <></>;
  }

  if (user?.activeTenant) {
    localStorage.setItem(TENANT_ID_KEY, user.activeTenant);
  }

  return <>{children}</>;
};

export default TenancyManager;
