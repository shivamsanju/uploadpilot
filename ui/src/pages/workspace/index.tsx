import { Group, Paper, Title } from '@mantine/core';
import { IconPlus } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { useGetWorkspaces } from '../../apis/workspace';
import { ErrorLoadingWrapper } from '../../components/ErrorLoadingWrapper';
import { WorkspaceCard } from '../../components/WorkspaceCard';
import classes from './Workspace.module.css';

const WorkspaceLandingPage = () => {
  const { isPending, error, workspaces } = useGetWorkspaces();
  const navigate = useNavigate();

  return (
    <ErrorLoadingWrapper error={error} isPending={isPending}>
      <Group align="center" gap="xs" h="10%">
        <Title order={3}>Workspaces</Title>
      </Group>
      <Group mb="50" align="center" mt="lg">
        <Paper
          withBorder
          h="200"
          w={{ base: '100%', md: '400' }}
          className={classes.wsItemAdd}
          onClick={() => navigate('/workspace/new')}
        >
          <Group justify="center" h="100%">
            <IconPlus size={30} stroke={2} color="gray" />
          </Group>
        </Paper>
        {workspaces?.length > 0 &&
          workspaces.map((workspace: any) => (
            <WorkspaceCard
              id={workspace.id}
              name={workspace.name}
              description={workspace.description}
              uploads={4031}
              storage={24.2}
              tags={workspace.tags || []}
            />
          ))}
      </Group>
    </ErrorLoadingWrapper>
  );
};

export default WorkspaceLandingPage;
