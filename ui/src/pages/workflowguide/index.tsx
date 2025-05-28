import { CodeHighlight } from '@mantine/code-highlight';
import { Box, Code, Group, Paper, Title } from '@mantine/core';
import { IconClipboardText } from '@tabler/icons-react';
import { useEffect } from 'react';
import Markdown from 'react-markdown';
import { useParams } from 'react-router-dom';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import { workflowGuideMd } from './guide';
import classes from './guide.module.css';

const components = {
  code: ({ node, className, children, ...props }: any) => {
    const isCodeBlock = node.position?.start.line !== node.position?.end.line;

    if (isCodeBlock) {
      return <CodeHighlight {...props} code={String(children).trim()} />;
    }

    return (
      <Code className={className} {...props}>
        {children}
      </Code>
    );
  },
};

const WorkflowBuilderGuidePage = () => {
  const { workspaceId } = useParams();
  const setBreadcrumbs = useSetBreadcrumbs();

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Workspaces', path: '/' },
      { label: 'Processors', path: `/workspace/${workspaceId}/processors` },
      { label: 'Builder guide' },
    ]);
  }, [setBreadcrumbs, workspaceId]);

  return (
    <Box mb={50}>
      <Group mb="xl">
        <IconClipboardText size={24} />
        <Title order={3}>Builder guide</Title>
      </Group>
      <Paper px="lg" withBorder className={classes.guideMd}>
        <Markdown components={components}>{workflowGuideMd}</Markdown>
      </Paper>
    </Box>
  );
};

export default WorkflowBuilderGuidePage;
