import { Box, Group, Paper, Title } from '@mantine/core';
import { IconAdjustments } from '@tabler/icons-react';
import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useGetUploaderConfig } from '../../apis/uploader';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import UploaderConfigForm from './Config';

const ConfigurationPage = () => {
  const { workspaceId } = useParams();
  const setBreadcrumbs = useSetBreadcrumbs();

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Workspaces', path: '/' },
      { label: 'Configuration' },
    ]);
  }, [setBreadcrumbs]);

  let { isPending, error, config } = useGetUploaderConfig(
    workspaceId as string,
  );

  if (!isPending && !error && !config) {
    error = new Error('No config found for this workspace');
  }

  if (error) {
    return <ErrorCard title={error.name} message={error.message} h="70vh" />;
  }

  return (
    <Box mb={50}>
      <Group mb="xl">
        <IconAdjustments size={24} />
        <Title order={3}>Configuration</Title>
      </Group>
      <Paper withBorder p="md">
        <UploaderConfigForm config={config} isPending={isPending} />
      </Paper>
    </Box>
  );
};

export default ConfigurationPage;
