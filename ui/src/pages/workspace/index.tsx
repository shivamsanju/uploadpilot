import { Box, Group, Paper, Title } from '@mantine/core';
import { IconPlus } from '@tabler/icons-react';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useGetWorkspaces } from '../../apis/workspace';
import { ErrorLoadingWrapper } from '../../components/ErrorLoadingWrapper';
import { WorkspaceCard } from '../../components/WorkspaceCard';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import classes from './Workspace.module.css';

const WorkspaceLandingPage = () => {
  const { isPending, error, workspaces } = useGetWorkspaces();
  const navigate = useNavigate();

  const setBreadcrumbs = useSetBreadcrumbs();

  useEffect(() => {
    setBreadcrumbs([{ label: '' }]);
  }, [setBreadcrumbs]);

  return (
    <ErrorLoadingWrapper error={error} isPending={isPending}>
      <Group align="center" gap="xs" h="10%" mb={30}>
        <Title order={3}>Workspaces</Title>
      </Group>
      <Box mb="lg" className={classes.wscontainer}>
        <Paper
          withBorder
          h="80"
          className={classes.wsItemAdd}
          onClick={() => navigate('/workspace/new')}
        >
          <IconPlus size={30} stroke={2} color="gray" />
        </Paper>
        {workspaces?.length > 0 &&
          workspaces.map((workspace: any) => (
            <Box className={classes.wsitem} key={workspace.id}>
              <WorkspaceCard
                id={workspace.id}
                name={workspace.name}
                description={workspace.description}
                uploads={4031}
                storage={24.2}
                tags={workspace.tags || []}
              />
            </Box>
          ))}
      </Box>
    </ErrorLoadingWrapper>
  );
};

export default WorkspaceLandingPage;
