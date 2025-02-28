import { useEffect } from 'react';
import { handleSocialLoginCallback } from '../../apis/auth';
import { AppLoader } from '../../components/Loader/AppLoader';

const SocialAuthCallbackHandlerPage = () => {
  useEffect(() => {
    handleSocialLoginCallback();
  }, []);

  return <AppLoader h="100vh" />;
};

export default SocialAuthCallbackHandlerPage;
