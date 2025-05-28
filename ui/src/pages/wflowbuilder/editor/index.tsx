import {
  Alert,
  Group,
  Text,
  Tooltip,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import Editor, { Monaco } from '@monaco-editor/react';
import { IconAlertCircle, IconCircleCheck } from '@tabler/icons-react';
import React, { useEffect, useState } from 'react';
import { validateWorkflowContent } from './schema';

type Props = {
  workflowContent: string;
  setWorkflowContent: React.Dispatch<React.SetStateAction<string>>;
};

export const WorkflowYamlEditor: React.FC<Props> = ({
  workflowContent,
  setWorkflowContent,
}) => {
  const [error, setError] = useState<string | null>('');
  const [initialLoad, setInitialLoad] = useState(true);
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();

  const handleEditorDidMount = (monaco: Monaco) => {
    monaco.editor.defineTheme('myCustomThemeDark', {
      base: 'vs-dark',
      inherit: true,
      rules: [{ token: 'comment', fontStyle: 'italic' }],
      colors: {
        'editor.background': theme.colors.dark[9],
      },
    });
    setInitialLoad(false);
  };

  useEffect(() => {
    setError(validateWorkflowContent(workflowContent || ''));
  }, [workflowContent]);

  return (
    <>
      <Editor
        loading={initialLoad}
        beforeMount={handleEditorDidMount}
        theme={colorScheme === 'dark' ? 'myCustomThemeDark' : 'vs'}
        language="yaml"
        height="70vh"
        defaultLanguage="yaml"
        value={workflowContent}
        onChange={(value: any) => {
          if (typeof value === 'string') {
            setWorkflowContent(value);
          }
        }}
        options={{
          minimap: { enabled: false },
          scrollBeyondLastLine: false,
          renderLineHighlight: 'none',
          padding: {
            top: 10,
            bottom: 50,
          },
          rulers: [],
        }}
      />
      <Alert color={error ? 'red' : 'green'} radius={0} p="xs">
        <Tooltip.Floating
          label={error || 'Workflow content is valid'}
          multiline
          maw="700"
        >
          <Group
            p={0}
            m={0}
            align="center"
            gap="xs"
            c={error ? 'red' : 'green'}
            wrap="nowrap"
          >
            {error ? (
              <IconAlertCircle size={16} stroke={1.5} />
            ) : (
              <IconCircleCheck size={16} stroke={1.5} />
            )}

            <Text lineClamp={1} w="95%">
              {error || 'Workflow content is valid'}
            </Text>
          </Group>
        </Tooltip.Floating>
      </Alert>
    </>
  );
};
