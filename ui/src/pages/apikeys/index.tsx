import { Box, Group, Title } from '@mantine/core';
import { useEffect } from 'react';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import ApiKeyList from './List';

const ApiKeyPage = () => {
  const setBreadcrumbs = useSetBreadcrumbs();

  useEffect(() => {
    setBreadcrumbs([{ label: 'Workspaces', path: '/' }, { label: 'API Keys' }]);
  }, [setBreadcrumbs]);

  return (
    <Box>
      <Group align="center" gap="xs" h="10%" mb="md">
        <Title order={3}>API Keys</Title>
      </Group>
      <ApiKeyList />
    </Box>
  );
};

export default ApiKeyPage;
