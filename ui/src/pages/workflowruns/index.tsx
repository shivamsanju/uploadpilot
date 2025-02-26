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
import { AppLoader } from '../../components/Loader/AppLoader';
import WorkflowRunsList from './List';

const WorkflowRunsPage = () => {
  const { workspaceId, processorId } = useParams();
  const navigate = useNavigate();

  if (!workspaceId) {
    return <AppLoader h="70vh" />;
  }

  return (
    <Box>
      <Breadcrumbs separator=">">
        <NavLink to="/" className="bredcrumb-link">
          <Text>Workspaces</Text>
        </NavLink>
        <NavLink
          to={`/workspace/${workspaceId}/processors`}
          className="bredcrumb-link"
        >
          <Text>Processors</Text>
        </NavLink>
        <Text>{processorId}</Text>
      </Breadcrumbs>
      <Group mt="xs" mb="xl">
        <ActionIcon
          variant="default"
          radius="xl"
          size="sm"
          onClick={() => navigate(`/workspace/${workspaceId}/processors`)}
        >
          <IconChevronLeft size={16} />
        </ActionIcon>
        <Title order={3}>Workflow runs</Title>
      </Group>
      <WorkflowRunsList />
    </Box>
  );
};

export default WorkflowRunsPage;
