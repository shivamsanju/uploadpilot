import {
  Box,
  Button,
  Group,
  Stack,
  Textarea,
  TextInput,
  Title,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useNavigate } from 'react-router-dom';
import { useOnboardTenant } from '../../apis/tenant';
import { TenantOnboardingRequest } from '../../types/tenant';
import { ContainerOverlay } from '../Overlay';

export const TenantRegistrationForm = ({
  enableCancel,
}: {
  enableCancel?: boolean;
}) => {
  const form = useForm<TenantOnboardingRequest>({
    initialValues: {
      name: '',
      contactEmail: '',
      phone: '',
      address: '',
      industry: '',
      companyName: '',
      role: '',
    },
    validate: {
      name: value => {
        if (!value.trim()) {
          return 'Tenant name is required';
        }
        if (value.trim().length > 25 || value.trim().length < 3) {
          return 'Tenant name must be between 3 and 25 characters';
        }

        if (!/^[a-zA-Z0-9 ]+$/.test(value)) {
          return 'Tenant name can only contain letters and numbers';
        }

        return null;
      },
      contactEmail: value => (/^\S+@\S+$/.test(value) ? null : 'Invalid email'),
      phone: value =>
        value && !/^\+?\d{7,15}$/.test(value) ? 'Invalid phone number' : null,
    },
    validateInputOnBlur: true,
  });

  const navigate = useNavigate();
  const { mutateAsync, isPending } = useOnboardTenant();

  const handleSubmit = async (values: TenantOnboardingRequest) => {
    try {
      await mutateAsync(values);
      window.location.href = '/';
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <>
      <ContainerOverlay visible={isPending} />

      <Stack align="center">
        <Box>
          <Title order={2} mb="xl">
            Tenant Registration
          </Title>
          <form onSubmit={form.onSubmit(handleSubmit)}>
            <Stack gap="xl" w={{ base: '100%', md: '500' }}>
              <TextInput
                label="Tenant Name"
                placeholder="Create a unique name for your tenant"
                {...form.getInputProps('name')}
                required
              />

              <TextInput
                label="Contact Email"
                placeholder="We will use this email to contact you"
                {...form.getInputProps('contactEmail')}
                required
              />

              <TextInput
                label="Phone (With country code)"
                {...form.getInputProps('phone')}
                placeholder="+1234567890"
                required
              />

              <TextInput
                label="Industry"
                placeholder="Your industry"
                {...form.getInputProps('industry')}
              />

              <Textarea
                label="Address"
                placeholder="Your company address"
                {...form.getInputProps('address')}
              />

              <Group mt="md" gap="sm">
                {enableCancel && (
                  <Button
                    variant="outline"
                    onClick={() => navigate('/tenants')}
                  >
                    Back
                  </Button>
                )}
                <Button type="submit" variant="white">
                  Submit
                </Button>
              </Group>
            </Stack>
          </form>
        </Box>
      </Stack>
    </>
  );
};
