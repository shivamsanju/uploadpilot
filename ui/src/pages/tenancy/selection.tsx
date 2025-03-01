import {
  Button,
  Container,
  Group,
  Stack,
  Text,
  TextInput,
  UnstyledButton,
} from '@mantine/core';
import { IconPlus, IconSearch } from '@tabler/icons-react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSetActiveTenant } from '../../apis/tenant';
import { useGetUserDetails } from '../../apis/user';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { AppLoader } from '../../components/Loader/AppLoader';
import { TenantInfoCard } from '../../components/TenantInfo';
import { TENANT_ID_KEY } from '../../constants/tenancy';

const TenantSelectionPage = () => {
  const [search, setSearch] = useState('');
  const { user, isPending, error } = useGetUserDetails();
  const { mutateAsync } = useSetActiveTenant();

  const navigate = useNavigate();

  if (isPending) {
    return <AppLoader h="100vh" />;
  }

  if (error) {
    return <ErrorCard title={error.name} message={error.message} h="100vh" />;
  }

  const tenants = user?.tenants || {};

  const filteredTenants = Object.keys(tenants).filter(tenant =>
    tenants[tenant]?.toLowerCase()?.includes(search.toLowerCase()),
  );

  const selectTenant = async (id: string) => {
    try {
      await mutateAsync(id);
      localStorage.setItem(TENANT_ID_KEY, id);
      window.location.href = '/';
    } catch (error) {
      console.error(error);
    }
  };

  const addNewTenant = async () => {
    navigate('/register-tenant');
  };

  return (
    <Container size="sm" mt={30}>
      <Group justify="space-between">
        <Text size="lg" fw={500}>
          Select a Tenant
        </Text>
        <Button
          leftSection={<IconPlus size={16} />}
          onClick={addNewTenant}
          variant="subtle"
        >
          Register New Tenant
        </Button>
      </Group>

      <TextInput
        placeholder="Search tenants..."
        leftSection={<IconSearch size={16} />}
        mt="md"
        value={search}
        onChange={event => setSearch(event.currentTarget.value)}
      />

      <Stack p={0} m={0} mt="md" gap="sm">
        {filteredTenants.map((tenant, index) => (
          <UnstyledButton onClick={() => selectTenant(tenant)}>
            <TenantInfoCard
              key={index}
              tenantName={tenants[tenant]}
              tenantId={tenant}
            />
          </UnstyledButton>
        ))}
        {filteredTenants.length === 0 && (
          <Text size="md" mt="md" c="dimmed">
            No tenants found
          </Text>
        )}
      </Stack>
    </Container>
  );
};

export default TenantSelectionPage;
