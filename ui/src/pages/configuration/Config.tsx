import {
  Group,
  MultiSelect,
  NumberInput,
  Paper,
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
import { useGetAllAllowedSources } from '../../apis/workspace';
import { DiscardButton } from '../../components/Buttons/DiscardButton';
import { SaveButton } from '../../components/Buttons/SaveButton';
import { ContainerOverlay } from '../../components/Overlay';
import { WorkspaceConfig } from '../../types/uploader';
import classes from './Form.module.css';

type NewUploaderConfigProps = {
  config: WorkspaceConfig;
};

const UploaderConfigForm: React.FC<NewUploaderConfigProps> = ({ config }) => {
  const { workspaceId } = useParams();
  const { isPending, allowedSources } = useGetAllAllowedSources(
    workspaceId || '',
  );
  const { mutateAsync, isPending: isPendingMutation } =
    useUpdateUploaderConfigMutation();

  const form = useForm<WorkspaceConfig>({
    initialValues: {
      ...config,
      allowedFileTypes: config?.allowedFileTypes || [],
      allowedSources: config?.allowedSources || [],
      requiredMetadataFields: config?.requiredMetadataFields || [],
      allowedOrigins: config?.allowedOrigins || [],
    },
    validate: {
      allowedSources: value =>
        value.length === 0 ? 'Please select at least one source' : null,
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
          minNumberOfFiles: form.values.minNumberOfFiles || 0,
          maxNumberOfFiles: form.values.maxNumberOfFiles || 0,
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
      <ContainerOverlay visible={isPending || isPendingMutation} />
      <Paper withBorder p="sm" mb={50}>
        <SimpleGrid cols={{ base: 1, xl: 2 }}>
          <Stack p="md">
            {/* Allowed input sources */}
            <MultiSelect
              label="Allowed input sources"
              description="Allowed input sources for your uploader"
              placeholder="Select allowed input sources"
              data={allowedSources || []}
              {...form.getInputProps('allowedSources')}
              disabled={isPending}
              searchable
            />
            {/* Min file size */}
            {/* <NumberInput
              label="Min file size"
              description="Enter minimum file size in bytes"
              {...form.getInputProps("minFileSize")}
              min={0}
            /> */}

            {/* Max file size */}
            <NumberInput
              label="Max file size"
              description="Enter maximum file size in bytes"
              placeholder="Enter maximum file size in bytes"
              {...form.getInputProps('maxFileSize')}
              min={0}
            />

            {/* Min number of files */}
            {/* <NumberInput
              label="Min number of files"
              description="Specify the minimum number of files required"
              {...form.getInputProps("minNumberOfFiles")}
              min={0}
            /> */}

            {/* Max number of files */}
            <NumberInput
              label="Max number of files"
              description="Specify the maximum number of files allowed"
              placeholder="Specify the maximum number of files allowed"
              {...form.getInputProps('maxNumberOfFiles')}
              min={1}
            />
          </Stack>
          <Stack p="md">
            {/* Allowed file types */}
            <TagsInput
              label="Allowed mime types"
              description="Allowed mime types for your uploader"
              placeholder="Comma separated mime types"
              {...form.getInputProps('allowedFileTypes')}
              min={0}
            />

            {/*Auth Endpoint */}
            <TagsInput
              label="Allowed origins"
              description="Allowed origins for your uploader"
              placeholder="Comma separated origins"
              {...form.getInputProps('allowedOrigins')}
              min={0}
            />
            <TagsInput
              label="Required metadata fields"
              placeholder="Comma separated fields"
              description="Required metadata fields for your uploader"
              {...form.getInputProps('requiredMetadataFields')}
              min={0}
            />
          </Stack>
        </SimpleGrid>
      </Paper>
      {/* <Title order={5} opacity={0.7} mb="sm" >Uploader Settings</Title> */}
      <Paper withBorder p="sm">
        <SimpleGrid cols={{ sm: 1, lg: 2 }}>
          <Stack p="md">
            <Group justify="space-between" className={classes.item}>
              <div>
                <Text size="sm">Allow pause and resume</Text>
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
                <Text size="sm">Enable image editing</Text>
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
          <Stack p="md">
            <Group justify="space-between" className={classes.item}>
              <div>
                <Text size="sm">Use compression</Text>
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
                <Text size="sm">Use fault tolerant mode</Text>
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
      </Paper>

      <Group justify="flex-end" gap="md" mt="xl">
        <DiscardButton type="reset" />
        <SaveButton type="submit" />
      </Group>
    </form>
  );
};

export default UploaderConfigForm;
