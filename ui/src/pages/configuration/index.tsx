import { Box, Group, Title } from '@mantine/core';
import { IconAdjustments } from '@tabler/icons-react';
import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useGetUploaderConfig } from '../../apis/uploader';
import { ErrorLoadingWrapper } from '../../components/ErrorLoadingWrapper';
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

  return (
    <Box mb={50}>
      <ErrorLoadingWrapper error={error} isPending={isPending}>
        <Group mb="xl">
          <IconAdjustments size={24} />
          <Title order={3}>Configuration</Title>
        </Group>
        <Box>
          <UploaderConfigForm config={config} />
        </Box>
      </ErrorLoadingWrapper>
    </Box>
  );
};

export default ConfigurationPage;
