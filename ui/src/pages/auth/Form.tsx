import {
  Anchor,
  Button,
  Checkbox,
  Divider,
  Group,
  Image,
  Modal,
  Paper,
  PaperProps,
  PasswordInput,
  ScrollArea,
  Stack,
  Text,
  TextInput,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { upperFirst, useToggle, useViewportSize } from '@mantine/hooks';
import { useState } from 'react';
import {
  handleEPSignIn,
  handleEPSignup,
  handleSocialLogin,
} from '../../apis/auth';
import GithubIcon from '../../assets/icons/github.svg';
import GoogleIcon from '../../assets/icons/google.svg';
import PrivacyPolicy from '../../components/Policy/PrivacyPolicy';
import TermsOfService from '../../components/Policy/TermsOfService';

export const AuthenticationForm = (props: PaperProps) => {
  const { width } = useViewportSize();
  const [type, toggle] = useToggle(['login', 'register']);
  const [modalOpen, setModalOpen] = useState(false);
  const [modalContent, setModalContent] = useState<'terms' | 'privacy'>(
    'terms',
  );

  const form = useForm({
    initialValues: {
      email: '',
      name: '',
      password: '',
      terms: true,
    },

    validate: {
      terms: val => (val ? null : 'You must accept terms and conditions'),
      email: val => (/^\S+@\S+$/.test(val) ? null : 'Invalid email'),
      password: val =>
        val.length <= 6
          ? 'Password should include at least 6 characters'
          : null,
    },
  });

  const handleLoginSignup = (values: any) => {
    if (type === 'login') {
      handleEPSignIn(values.email, values.password);
    } else {
      handleEPSignup(values.name, values.email, values.password);
    }
  };

  const handleModalOpen = (content: 'terms' | 'privacy') => {
    setModalContent(content);
    setModalOpen(true);
  };

  return (
    <Paper
      radius="md"
      p="xl"
      withBorder={width > 768}
      w={width < 768 ? '100%' : '400px'}
      {...props}
    >
      <Text size="lg" fw={500}>
        Welcome to UploadPilot, {type} with
      </Text>

      <Group grow mb="md" mt="md">
        <Button
          variant="default"
          radius="xl"
          leftSection={<Image src={GoogleIcon} width={20} height={20} />}
          onClick={() => handleSocialLogin('google')}
          size="sm"
        >
          Google
        </Button>
        <Button
          variant="default"
          radius="xl"
          leftSection={<Image src={GithubIcon} width={20} height={20} />}
          onClick={() => handleSocialLogin('github')}
          size="sm"
        >
          Github
        </Button>
      </Group>

      <Divider label="Or continue with email" labelPosition="center" my="lg" />

      <form onSubmit={form.onSubmit(handleLoginSignup)}>
        <Stack>
          {type === 'register' && (
            <TextInput
              label="Name"
              placeholder="Your name"
              value={form.values.name}
              onChange={event =>
                form.setFieldValue('name', event.currentTarget.value)
              }
              radius="md"
            />
          )}

          <TextInput
            required
            label="Email"
            placeholder="hello@mantine.dev"
            value={form.values.email}
            onChange={event =>
              form.setFieldValue('email', event.currentTarget.value)
            }
            error={form.errors.email && 'Invalid email'}
            radius="md"
          />

          <PasswordInput
            required
            label="Password"
            placeholder="Your password"
            value={form.values.password}
            onChange={event =>
              form.setFieldValue('password', event.currentTarget.value)
            }
            error={
              form.errors.password &&
              'Password should include at least 6 characters'
            }
            radius="md"
          />

          {type === 'register' && (
            <Group align="center">
              <Checkbox
                checked={form.values.terms}
                onChange={event =>
                  form.setFieldValue('terms', event.currentTarget.checked)
                }
              />
              <Text size="xs" ta="center" c="dimmed" mt="2 ">
                I accept the{' '}
                <Anchor
                  onClick={() => handleModalOpen('terms')}
                  style={{ cursor: 'pointer' }}
                >
                  Terms of Service
                </Anchor>{' '}
                and{' '}
                <Anchor
                  onClick={() => handleModalOpen('privacy')}
                  style={{ cursor: 'pointer' }}
                >
                  Privacy Policy
                </Anchor>
              </Text>
            </Group>
          )}
        </Stack>

        <Group justify="space-between" mt="xl">
          <Anchor
            component="button"
            type="button"
            c="dimmed"
            onClick={() => toggle()}
            size="xs"
          >
            {type === 'register'
              ? 'Already have an account? Login'
              : "Don't have an account? Register"}
          </Anchor>
          <Button type="submit" radius="xl">
            {upperFirst(type)}
          </Button>
        </Group>
      </form>
      <Modal
        opened={modalOpen}
        onClose={() => setModalOpen(false)}
        title={modalContent === 'terms' ? 'Terms of Service' : 'Privacy Policy'}
        size="xl"
        fullScreen
        scrollAreaComponent={props => (
          <ScrollArea.Autosize {...props} scrollbarSize={10} />
        )}
      >
        {modalContent === 'terms' ? <TermsOfService /> : <PrivacyPolicy />}
      </Modal>
    </Paper>
  );
};
