import { useGetUserDetails } from '../../apis/user';
import { ErrorLoadingWrapper } from '../../components/ErrorLoadingWrapper';
import { TenantRegistrationForm } from '../../components/Tenancy/RegistrationForm';

export const TenantRegistrationPage = () => {
  const { user, isPending, error } = useGetUserDetails();

  return (
    <ErrorLoadingWrapper isPending={isPending} error={error}>
      <TenantRegistrationForm
        enableCancel={user?.tenants && Object.keys(user.tenants).length > 0}
      />
    </ErrorLoadingWrapper>
  );
};
