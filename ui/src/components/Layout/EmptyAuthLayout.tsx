import TenancyManager from '../Tenancy/TenancyManager';

const EmptyAuthLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return <TenancyManager>{children}</TenancyManager>;
};

export default EmptyAuthLayout;
