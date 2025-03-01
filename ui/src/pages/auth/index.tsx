import { Image, Stack } from '@mantine/core';
import DarkLogo from '../../assets/images/full-logo-dark.png';
import { AuthenticationForm } from './Form';

const AuthPage = () => {
  return (
    <Stack align="center" justify="center" h="100vh" w="100%">
      <Image src={DarkLogo} alt="logo" h={60} w={200} />
      <AuthenticationForm />
    </Stack>
  );
};

export default AuthPage;
