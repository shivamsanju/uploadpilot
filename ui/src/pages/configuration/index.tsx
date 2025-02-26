import {
  ActionIcon,
  Box,
  Breadcrumbs,
  Group,
  Text,
  Title,
} from '@mantine/core';
import { IconChevronLeft } from '@tabler/icons-react';
import { NavLink, useNavigate, useParams } from 'react-router-dom';
import { useGetUploaderConfig } from '../../apis/uploader';
import { ErrorLoadingWrapper } from '../../components/ErrorLoadingWrapper';
import UploaderConfigForm from './Config';

const ConfigurationPage = () => {
  const { workspaceId } = useParams();
  const navigate = useNavigate();

  let { isPending, error, config } = useGetUploaderConfig(
    workspaceId as string,
  );

  if (!isPending && !error && !config) {
    error = new Error('No config found for this workspace');
  }

  return (
    <ErrorLoadingWrapper error={error} isPending={isPending}>
      <Breadcrumbs separator=">">
        <NavLink to="/" className="bredcrumb-link">
          <Text>Workspaces</Text>
        </NavLink>
        <Text>Configuration</Text>
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
        <Title order={3}>Configuration</Title>
      </Group>

      <Box mb={50}>
        <UploaderConfigForm config={config} />
      </Box>
    </ErrorLoadingWrapper>
  );
};

export default ConfigurationPage;
