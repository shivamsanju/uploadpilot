import {
  Anchor,
  Box,
  Button,
  Group,
  Image,
  Modal,
  ScrollArea,
  Stack,
  Text,
} from '@mantine/core';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import GithubIcon from '../../assets/icons/github.svg';
import GoogleIcon from '../../assets/icons/google.svg';
import LoginSvg from '../../assets/images/login.svg';
import { Lines } from '../../components/Lines';
import { AppLoader } from '../../components/Loader/AppLoader';
import { Logo2 } from '../../components/Logo/Logo2';
import PrivacyPolicy from '../../components/Policy/PrivacyPolicy';
import TermsOfService from '../../components/Policy/TermsOfService';
import { getApiDomain } from '../../utils/config';
import classes from './Auth.module.css';
import TokenCatcher from './TokenCatcher';

const AuthPage = () => {
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [modalContent, setModalContent] = useState<'terms' | 'privacy'>(
    'terms',
  );
  const navigate = useNavigate();

  const handleLogin = (provider: string) => {
    window.location.href = getApiDomain() + `/auth/${provider}/authorize`;
  };

  useEffect(() => {
    const token = localStorage.getItem('uploadpilottoken');
    if (token) {
      navigate('/', { replace: true });
    }
    setLoading(false);
  }, [navigate]);

  const openModal = (content: 'terms' | 'privacy') => {
    setModalContent(content);
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
  };

  if (loading) {
    return <AppLoader h="100vh" />;
  }

  return (
    <ScrollArea h="100vh" w="100vw" offsetScrollbars={false} scrollbarSize={6}>
      <Group justify="center" align="center" h="100vh">
        <TokenCatcher />
        <Lines />
        <Box w={{ base: '100vw', sm: 600 }} className={classes.form}>
          <Stack gap="xs" mb="60" align="center">
            <Logo2 enableOnClick={false} height="90px" width="320px" />
          </Stack>
          <Image src={LoginSvg} alt="login" height={300} width={100} />
          <Stack>
            <Text size="xs" ta="center" c="dimmed" mt="sm">
              By continuing, you agree to our{' '}
              <Anchor
                onClick={() => openModal('terms')}
                style={{ cursor: 'pointer' }}
              >
                Terms of Service
              </Anchor>{' '}
              and acknowledge you have read our{' '}
              <Anchor
                onClick={() => openModal('privacy')}
                style={{ cursor: 'pointer' }}
              >
                Privacy Policy
              </Anchor>
            </Text>
            <Button
              c="dark"
              variant="white"
              leftSection={<Image src={GoogleIcon} width={20} height={20} />}
              onClick={() => handleLogin('google')}
              size="sm"
            >
              Google
            </Button>
            <Text ta="center" c="dimmed">
              or
            </Text>
            <Button
              c="dark"
              variant="white"
              leftSection={<Image src={GithubIcon} width={25} height={25} />}
              onClick={() => handleLogin('github')}
              size="sm"
            >
              Github
            </Button>
          </Stack>
        </Box>

        {/* Modal for Terms of Service or Privacy Policy */}
        <Modal
          opened={modalOpen}
          onClose={closeModal}
          title={
            modalContent === 'terms' ? 'Terms of Service' : 'Privacy Policy'
          }
          size="xl"
          fullScreen
          scrollAreaComponent={props => (
            <ScrollArea.Autosize {...props} scrollbarSize={10} />
          )}
        >
          {modalContent === 'terms' ? <TermsOfService /> : <PrivacyPolicy />}
        </Modal>
      </Group>
    </ScrollArea>
  );
};

export default AuthPage;
