import { Box, Group, Title } from '@mantine/core';
import { IconTools } from '@tabler/icons-react';
import { useEffect } from 'react';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import ProcessorsList from './List';

const ProcessorsPage = () => {
  const setBreadcrumbs = useSetBreadcrumbs();

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Workspaces', path: '/' },
      { label: 'Processors' },
    ]);
  }, [setBreadcrumbs]);

  return (
    <Box mb={50}>
      <Group mb="xl">
        <IconTools size={24} />
        <Title order={3}>Processors</Title>
      </Group>
      <ProcessorsList />
    </Box>
  );
};

export default ProcessorsPage;
