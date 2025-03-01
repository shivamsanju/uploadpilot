import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { handleSocialLoginCallback } from '../../apis/auth';
import { AppLoader } from '../../components/Loader/AppLoader';

const SocialAuthCallbackHandlerPage = () => {
  const { provider } = useParams();
  useEffect(() => {
    handleSocialLoginCallback(provider);
  }, [provider]);

  return <AppLoader h="100vh" />;
};

export default SocialAuthCallbackHandlerPage;
