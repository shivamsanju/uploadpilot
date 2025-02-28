import {
  Box,
  Button,
  Group,
  Image,
  Select,
  SimpleGrid,
  Stack,
  Textarea,
  TextInput,
  Title,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useNavigate } from 'react-router-dom';
import { useOnboardTenant } from '../../apis/tenant';
import onboardingSvg from '../../assets/images/onboarding.svg';
import { TenantOnboardingRequest } from '../../types/tenant';
import { FullScreenOverlay } from '../Overlay';

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
      companyName: value => (value ? null : 'Company name is required'),
      role: value => (value ? null : 'Role is required'),
    },
    validateInputOnChange: true,
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
    <Box>
      <Group justify="center">
        <Title order={2} mt="xl" mb="xl">
          Tenant Registration
        </Title>
      </Group>
      <FullScreenOverlay visible={isPending} />
      <SimpleGrid cols={{ base: 1, md: 2 }} px="xl" pb="xl" spacing="100">
        <Stack align="flex-end" justify="center" visibleFrom="md">
          <Image src={onboardingSvg} h={500} w={500} />
        </Stack>
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <Stack gap="sm" mx="xl" p={0}>
            <TextInput
              maw="400"
              label="Tenant Name"
              placeholder="Create a unique name for your tenant"
              {...form.getInputProps('name')}
              required
            />

            <TextInput
              maw="400"
              label="Contact Email"
              placeholder="We will use this email to contact you"
              {...form.getInputProps('contactEmail')}
              required
            />

            <TextInput
              maw="400"
              label="Phone (With country code)"
              {...form.getInputProps('phone')}
              placeholder="+1234567890"
              required
            />

            <TextInput
              maw="400"
              label="Company Name"
              placeholder="Your company name"
              {...form.getInputProps('companyName')}
              required
            />

            <Select
              maw="400"
              label="Role"
              placeholder="Your role in the company"
              {...form.getInputProps('role')}
              required
              data={[
                { value: 'c-suite', label: 'C-Suite (CEO, CTO, etc.)' },
                { value: 'vp', label: 'VP (Vice President)' },
                { value: 'developer', label: 'Developer' },
              ]}
            />

            <TextInput
              maw="400"
              label="Industry"
              placeholder="Your industry"
              {...form.getInputProps('industry')}
            />

            <Textarea
              maw="400"
              label="Address"
              placeholder="Your company address"
              {...form.getInputProps('address')}
            />

            <Group mt="md" gap="sm">
              {enableCancel && (
                <Button variant="default" onClick={() => navigate('/tenants')}>
                  Back
                </Button>
              )}
              <Button type="submit">Submit</Button>
            </Group>
          </Stack>
        </form>
      </SimpleGrid>
    </Box>
  );
};
