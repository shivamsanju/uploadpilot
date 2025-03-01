import { Box, Grid, Group, Paper, Title } from '@mantine/core';
import { IconFileStack } from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useGetProcessor } from '../../apis/processors';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../components/Overlay';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import { BlockSearch } from './blocksearch';
import { WorkflowYamlEditor } from './editor';
import { EditorHeader } from './Header';

const WorkflowBuilderPage = () => {
  const { workspaceId, processorId } = useParams();
  const setBreadcrumbs = useSetBreadcrumbs();
  const [workflowContent, setWorkflowContent] = useState<string>('');

  const { isPending, error, processor } = useGetProcessor(
    workspaceId as string,
    processorId as string,
  );

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Workspaces', path: '/' },
      { label: 'Processors', path: `/workspace/${workspaceId}/processors` },
      { label: 'Workflow' },
    ]);
  }, [setBreadcrumbs, workspaceId]);

  useEffect(() => {
    setWorkflowContent(processor?.workflow || '# Start your workflow here');
  }, [processor]);

  if (error) {
    return <ErrorCard title={error.name} message={error.message} h="70vh" />;
  }

  return (
    <Box mb={50}>
      <ContainerOverlay visible={isPending} />
      <Group mb="xl">
        <IconFileStack size={24} />
        <Title order={3}>Workflow</Title>
      </Group>
      <Paper withBorder bg="transparent">
        <EditorHeader
          processor={processor}
          workspaceId={workspaceId || ''}
          workflowContent={workflowContent}
          setWorkflowContent={setWorkflowContent}
        />
        <Grid>
          <Grid.Col span={{ base: 12, lg: 3 }}>
            <BlockSearch />
          </Grid.Col>
          <Grid.Col span={{ base: 12, lg: 9 }}>
            {processor && (
              <WorkflowYamlEditor
                workflowContent={workflowContent}
                setWorkflowContent={setWorkflowContent}
              />
            )}
          </Grid.Col>
        </Grid>
      </Paper>
    </Box>
  );
};

export default WorkflowBuilderPage;
