import { Button, Group, Title } from '@mantine/core';
import { showNotification } from '@mantine/notifications';
import { IconDeviceFloppy, IconMenu4, IconRestore } from '@tabler/icons-react';
import { useUpdateProcessorWorkflowMutation } from '../../apis/processors';
import { showConfirmationPopup } from '../../components/Popups/ConfirmPopup';
import classes from './Builder.module.css';
import { validateWorkflowContent } from './editor/schema';

type Props = {
  workspaceId: string;
  processor: any;
  workflowContent: string;
  setWorkflowContent: React.Dispatch<React.SetStateAction<string>>;
};

export const EditorHeader: React.FC<Props> = ({
  workspaceId,
  processor,
  workflowContent,
  setWorkflowContent,
}) => {
  const { mutateAsync, isPending } = useUpdateProcessorWorkflowMutation();

  const saveWorkflow = async () => {
    try {
      const err = validateWorkflowContent(workflowContent || '');
      if (err !== null) {
        showNotification({
          color: 'red',
          title: 'Error',
          message: 'Invalid YAML: ' + err,
        });
        return;
      }
      await mutateAsync({
        workspaceId,
        processorId: processor?.id,
        workflow: workflowContent?.replace(/\t/g, '  '),
      });
    } catch (error: any) {
      //   setError(error?.response?.data?.message || error.message);
    }
  };

  const discardChanges = () => {
    showConfirmationPopup({
      message:
        'Are you sure you want to discard the changes? this is irreversible.',
      onOk: () => {
        setWorkflowContent(processor?.workflow || '');
      },
    });
  };

  return (
    <>
      <Group
        justify="space-between"
        align="center"
        p="xs"
        className={classes.editorHeader}
      >
        <Group gap="xs">
          <IconMenu4 size={18} />
          <Title order={4} opacity={0.8}>
            Steps
          </Title>
        </Group>

        <Group gap="md">
          <Button
            variant="subtle"
            leftSection={<IconRestore size={18} />}
            onClick={discardChanges}
            disabled={isPending}
          >
            Discard
          </Button>
          <Button
            variant="subtle"
            leftSection={<IconDeviceFloppy size={18} />}
            onClick={saveWorkflow}
            disabled={isPending}
          >
            Save
          </Button>
        </Group>
      </Group>
    </>
  );
};
