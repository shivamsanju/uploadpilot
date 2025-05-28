import { Box, Button, Group, Title } from '@mantine/core';
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

  useEffect(() => {
    if (!isPending && !error && workspaces?.length === 0) {
      navigate('/workspace/new');
    }
  }, [workspaces, navigate, isPending, error]);

  return (
    <Box mr="xl">
      <ErrorLoadingWrapper error={error} isPending={isPending}>
        <Group align="center" gap="xs" h="10%" mb={30} justify="space-between">
          <Title order={3}>Workspaces</Title>
        </Group>
        <Button
          leftSection={<IconPlus size={15} />}
          variant="subtle"
          onClick={() => navigate('/workspace/new')}
          mb="md"
        >
          Create
        </Button>
        <Box mb="lg" className={classes.wscontainer}>
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
                  totalUploads={2400}
                  storageUsed={3000}
                  processors={3}
                />
              </Box>
            ))}
        </Box>
      </ErrorLoadingWrapper>
    </Box>
  );
};

export default WorkspaceLandingPage;
