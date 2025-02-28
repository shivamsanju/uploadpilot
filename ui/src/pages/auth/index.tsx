import { Stack } from '@mantine/core';
import { AuthenticationForm } from './Form';

const AuthPage = () => {
  return (
    <Stack align="center" justify="center" h="100vh" w="100%">
      <AuthenticationForm />
    </Stack>
  );
};

export default AuthPage;
