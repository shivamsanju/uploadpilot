import {
  Button,
  Group,
  NumberInput,
  SimpleGrid,
  Stack,
  Switch,
  TagsInput,
  Text,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { showNotification } from '@mantine/notifications';
import { useParams } from 'react-router-dom';
import { useUpdateUploaderConfigMutation } from '../../apis/uploader';
import { ContainerOverlay } from '../../components/Overlay';
import { WorkspaceConfig } from '../../types/uploader';
import classes from './Form.module.css';

type NewUploaderConfigProps = {
  config: WorkspaceConfig;
};

const UploaderConfigForm: React.FC<NewUploaderConfigProps> = ({ config }) => {
  const { workspaceId } = useParams();

  const { mutateAsync, isPending } = useUpdateUploaderConfigMutation();

  const form = useForm<WorkspaceConfig>({
    initialValues: {
      ...config,
      allowedContentTypes: config?.allowedContentTypes || [],
      requiredMetadataFields: config?.requiredMetadataFields || [],
      allowedOrigins: config?.allowedOrigins || [],
    },
  });

  const handleEditAndSaveButton = async () => {
    if (!workspaceId) {
      showNotification({
        color: 'red',
        title: 'Error',
        message: 'Workspace ID is not available',
      });
      return;
    }
    if (form.isDirty()) {
      mutateAsync({
        workspaceId: workspaceId,
        config: {
          ...form.values,
          minFileSize: form.values.minFileSize || 0,
          maxFileSize: form.values.maxFileSize || 0,
        },
      }).catch(error => {
        console.log(error);
      });
    }
    form.resetDirty();
  };

  const handleResetButton = () => {
    form.reset();
    form.resetDirty();
  };

  return (
    <form
      onSubmit={form.onSubmit(handleEditAndSaveButton)}
      onReset={handleResetButton}
    >
      <ContainerOverlay visible={isPending} />
      <SimpleGrid cols={{ base: 1, xl: 2 }} spacing="50">
        <Stack gap="xl">
          <NumberInput
            required
            label="Max upload url validity (seconds)"
            placeholder="Enter maximum upload url validity in seconds"
            {...form.getInputProps('maxUploadURLLifetimeSecs')}
            min={0}
          />
          {/* Max file size */}
          <NumberInput
            label="Max file size"
            placeholder="Enter maximum file size in bytes"
            {...form.getInputProps('maxFileSize')}
            min={0}
          />
          {/* Min file size */}
          <NumberInput
            label="Min file size"
            placeholder="Enter minimum file size in bytes"
            {...form.getInputProps('minFileSize')}
            min={0}
          />
        </Stack>
        <Stack gap="xl">
          <TagsInput
            label="Allowed content types"
            placeholder="Comma separated content types"
            {...form.getInputProps('allowedContentTypes')}
            min={0}
          />
          {/*Allowed origins */}
          <TagsInput
            label="Allowed origins"
            placeholder="Comma separated origins"
            {...form.getInputProps('allowedOrigins')}
            min={0}
          />
          <TagsInput
            label="Required metadata fields"
            placeholder="Comma separated fields"
            {...form.getInputProps('requiredMetadataFields')}
            min={0}
          />
        </Stack>
      </SimpleGrid>
      <SimpleGrid cols={{ sm: 1, lg: 2 }} spacing="50" mt="xl">
        <Stack gap="xl">
          <Group justify="space-between" className={classes.item}>
            <div>
              <Text fw="500">Allow pause and resume</Text>
              <Text c="dimmed">
                Toggle to allow pause and resume in the uploader
              </Text>
            </div>
            <Switch
              className={classes.cusomSwitch}
              onLabel="ON"
              offLabel="OFF"
              checked={form.values.allowPauseAndResume}
              onChange={e =>
                form.setFieldValue('allowPauseAndResume', e.target.checked)
              }
            />
          </Group>

          <Group justify="space-between" className={classes.item}>
            <div>
              <Text fw="500">Enable image editing</Text>
              <Text c="dimmed">
                Toggle to enable image editing in the uploader ui
              </Text>
            </div>
            <Switch
              className={classes.cusomSwitch}
              onLabel="ON"
              offLabel="OFF"
              checked={form.values.enableImageEditing}
              onChange={e =>
                form.setFieldValue('enableImageEditing', e.target.checked)
              }
            />
          </Group>
        </Stack>
        <Stack gap="xl">
          <Group justify="space-between" className={classes.item}>
            <div>
              <Text fw="500">Use compression</Text>
              <Text c="dimmed">
                Toggle to enable compression while uploading files
              </Text>
            </div>
            <Switch
              className={classes.cusomSwitch}
              onLabel="ON"
              offLabel="OFF"
              checked={form.values.useCompression}
              onChange={e =>
                form.setFieldValue('useCompression', e.target.checked)
              }
            />
          </Group>

          <Group justify="space-between" className={classes.item}>
            <div>
              <Text fw="500">Use fault tolerant mode</Text>
              <Text c="dimmed">
                Fault tolerant mode allows to recover from browser crashes
              </Text>
            </div>
            <Switch
              className={classes.cusomSwitch}
              onLabel="ON"
              offLabel="OFF"
              checked={form.values.useFaultTolerantMode}
              onChange={e =>
                form.setFieldValue('useFaultTolerantMode', e.target.checked)
              }
            />
          </Group>
        </Stack>
      </SimpleGrid>

      <Group gap="md" mt="50">
        <Button variant="outline" type="reset">
          Reset
        </Button>
        <Button type="submit">Save</Button>
      </Group>
    </form>
  );
};

export default UploaderConfigForm;
