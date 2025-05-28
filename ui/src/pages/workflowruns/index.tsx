import { Box, Group, Title } from '@mantine/core';
import { IconBolt } from '@tabler/icons-react';
import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import ProcessorRunsList from './List';

const WorkflowRunsPage = () => {
  const { workspaceId } = useParams();
  const setBreadcrumbs = useSetBreadcrumbs();

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Workspaces', path: '/' },
      { label: 'Processors', path: `/workspace/${workspaceId}/processors` },
      { label: 'Runs' },
    ]);
  }, [setBreadcrumbs, workspaceId]);

  return (
    <Box mb={50} mr="md">
      <Group mb="xl">
        <IconBolt size={24} />
        <Title order={3}>Runs</Title>
      </Group>
      <ProcessorRunsList />
    </Box>
  );
};

export default WorkflowRunsPage;
