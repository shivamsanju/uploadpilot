import Editor, { Monaco, OnMount } from '@monaco-editor/react';
import React, { useState } from 'react';

import {
  Box,
  Group,
  LoadingOverlay,
  Text,
  Title,
  Tooltip,
  useMantineColorScheme,
} from '@mantine/core';
import { IconAlertCircle, IconCircleCheck } from '@tabler/icons-react';
import { useUpdateProcessorWorkflowMutation } from '../../../apis/processors';
import { DiscardButton } from '../../../components/Buttons/DiscardButton';
import { SaveButton } from '../../../components/Buttons/SaveButton';
import { showConfirmationPopup } from '../../../components/Popups/ConfirmPopup';
import { validateYaml } from './schema';

type Props = {
  workspaceId: string;
  processor: any;
  editor: any;
  setEditor: any;
};

export const WorkflowYamlEditor: React.FC<Props> = ({
  processor,
  workspaceId,
  setEditor,
  editor,
}) => {
  const err = validateYaml(processor?.workflow || '');
  const [yamlContent, setYamlContent] = useState<string>(
    processor?.workflow || '',
  );
  const [error, setError] = useState<string | null>(err);
  const [initialLoad, setInitialLoad] = useState(true);
  const { colorScheme } = useMantineColorScheme();

  const { mutateAsync, isPending } = useUpdateProcessorWorkflowMutation();

  const editorMount: OnMount = editorL => {
    setEditor(editorL);
  };

  const handleEditorDidMount = (monaco: Monaco) => {
    monaco.editor.defineTheme('myCustomThemeDark', {
      base: 'vs-dark',
      inherit: true,
      rules: [{ token: 'comment', fontStyle: 'italic' }],
      colors: {
        'editor.background': '#141414',
      },
    });
    setInitialLoad(false);
  };

  const saveYaml = async () => {
    try {
      await mutateAsync({
        workspaceId,
        processorId: processor?.id,
        workflow: yamlContent?.replace(/\t/g, '  '),
      });
    } catch (error: any) {
      setError(error?.response?.data?.message || error.message);
    }
  };

  const discardChanges = () => {
    showConfirmationPopup({
      message:
        'Are you sure you want to discard the changes? this is irreversible.',
      onOk: () => {
        setYamlContent(processor?.workflow || '');
      },
    });
  };

  return (
    <Box>
      <LoadingOverlay
        visible={isPending || initialLoad}
        overlayProps={{ backgroundOpacity: 0 }}
        zIndex={1000}
      />
      <Group justify="space-between" align="center" p="xs">
        <Box w="65%">
          <Title order={4} opacity={0.8}>
            Steps
          </Title>
          <Group
            align="center"
            gap={2}
            c={error ? 'red' : 'dimmed'}
            p={0}
            pt={2}
            wrap="nowrap"
          >
            <Box w="12">
              {error ? (
                <IconAlertCircle size="12" />
              ) : (
                <IconCircleCheck size="12" />
              )}
            </Box>
            <Tooltip
              multiline
              w={500}
              maw="90vw"
              label={error || 'Everything looks good'}
              color={error ? 'red' : 'dimmed'}
            >
              <Text size="xs" lineClamp={1}>
                {error || 'Everything looks good'}
              </Text>
            </Tooltip>
          </Group>
        </Box>
        <Group gap="md">
          <DiscardButton onClick={discardChanges} />
          <SaveButton onClick={saveYaml} />
        </Group>
      </Group>

      <Editor
        loading={false}
        beforeMount={handleEditorDidMount}
        onMount={editorMount}
        theme={colorScheme === 'dark' ? 'myCustomThemeDark' : 'vs'}
        language="yaml"
        height="70vh"
        defaultLanguage="yaml"
        value={yamlContent}
        onChange={(value: any) => {
          if (typeof value === 'string') {
            setYamlContent(value);
            const err = validateYaml(value);
            setError(err);
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
    </Box>
  );
};
