import {
  Button,
  Group,
  NumberInput,
  SimpleGrid,
  Stack,
  TagsInput,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { showNotification } from '@mantine/notifications';
import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useUpdateUploaderConfigMutation } from '../../apis/uploader';
import { ContainerOverlay } from '../../components/Overlay';
import { WorkspaceConfig } from '../../types/uploader';

type NewUploaderConfigProps = {
  config: WorkspaceConfig;
  isPending: boolean;
};

const UploaderConfigForm: React.FC<NewUploaderConfigProps> = ({
  config,
  isPending,
}) => {
  const { workspaceId } = useParams();

  const { mutateAsync, isPending: isUpdating } =
    useUpdateUploaderConfigMutation();

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

  useEffect(() => {
    form.setValues({
      ...config,
      allowedContentTypes: config?.allowedContentTypes || [],
      requiredMetadataFields: config?.requiredMetadataFields || [],
      allowedOrigins: config?.allowedOrigins || [],
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [config]);

  const handleResetButton = () => {
    form.reset();
    form.resetDirty();
  };

  return (
    <form
      onSubmit={form.onSubmit(handleEditAndSaveButton)}
      onReset={handleResetButton}
    >
      <ContainerOverlay visible={isPending || isUpdating} />
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
            label="Max file size (in bytes)"
            placeholder="Enter maximum file size in bytes"
            {...form.getInputProps('maxFileSize')}
            min={0}
          />
          {/* Min file size */}
          <NumberInput
            label="Min file size (in bytes)"
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

      <Group gap="md" mt="70" justify="flex-end">
        <Button variant="default" type="reset">
          Reset
        </Button>
        <Button type="submit">Save</Button>
      </Group>
    </form>
  );
};

export default UploaderConfigForm;
