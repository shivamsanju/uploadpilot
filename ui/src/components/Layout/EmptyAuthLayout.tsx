import SessionManager from '../SessionManager/SessionManager';
import TenancyManager from '../Tenancy/TenancyManager';

const EmptyAuthLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return (
    <SessionManager>
      <TenancyManager>{children}</TenancyManager>
    </SessionManager>
  );
};

export default EmptyAuthLayout;
