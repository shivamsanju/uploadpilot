import AuthWrapper from "../AuthWrapper/AuthWrapper";

const EmptyAuthLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return <AuthWrapper>{children}</AuthWrapper>;
};

export default EmptyAuthLayout;
