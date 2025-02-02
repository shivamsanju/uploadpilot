import { useNavigate } from "react-router-dom";
import { useGetSession } from "../../apis/user";
import { AppLoader } from "../Loader/AppLoader";

const AuthWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isPending, error, session } = useGetSession();
  const navigate = useNavigate();

  if (isPending) {
    return <AppLoader h="100vh" />;
  }

  if (error || !session) {
    navigate("/auth", { replace: true });
  }

  return <>{children}</>;
};

export default AuthWrapper;
