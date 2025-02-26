import {
  ActionIcon,
  Box,
  Breadcrumbs,
  Group,
  Text,
  Title,
} from '@mantine/core';
import { IconChevronLeft } from '@tabler/icons-react';
import { NavLink, useNavigate } from 'react-router-dom';
import ApiKeyList from './List';

const ApiKeyPage = () => {
  const navigate = useNavigate();

  return (
    <Box>
      <Breadcrumbs separator=">">
        <NavLink to="/" className="bredcrumb-link">
          <Text>Workspaces</Text>
        </NavLink>
        <Text>API Keys</Text>
      </Breadcrumbs>
      <Group mt="xs" mb="xl">
        <ActionIcon
          variant="default"
          radius="xl"
          size="sm"
          onClick={() => navigate(`/`)}
        >
          <IconChevronLeft size={16} />
        </ActionIcon>
        <Group align="center" gap="xs" h="10%">
          <Title order={3}>API Keys</Title>
        </Group>
      </Group>

      <ApiKeyList />
    </Box>
  );
};

export default ApiKeyPage;
