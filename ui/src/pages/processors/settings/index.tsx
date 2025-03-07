import { Box, Button, Group, TagsInput, TextInput, Title } from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconSettings } from '@tabler/icons-react';
import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  useGetProcessor,
  useUpdateProcessorMutation,
} from '../../../apis/processors';
import { ErrorCard } from '../../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../../components/Overlay';
import { useSetBreadcrumbs } from '../../../hooks/breadcrumb';

const ProcessorSettingsPage = () => {
  const { workspaceId, processorId } = useParams();
  const setBreadcrumbs = useSetBreadcrumbs();

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Workspaces', path: '/' },
      { label: 'Processors', path: `/workspace/${workspaceId}/processors` },
      { label: 'Settings' },
    ]);
  }, [setBreadcrumbs, workspaceId]);

  const form = useForm({
    initialValues: {
      name: '',
      triggers: [],
      enabled: true,
      templateKey: '',
    },
    validate: {
      name: value => (value ? null : 'Name is required'),
    },
  });

  const { mutateAsync, isPending: isCreating } = useUpdateProcessorMutation();
  const { isPending, error, processor } = useGetProcessor(
    workspaceId || '',
    processorId || '',
  );

  const handleUpdate = async (values: any) => {
    try {
      await mutateAsync({
        workspaceId: workspaceId || '',
        processorId: processorId || '',
        processor: {
          name: values.name,
          triggers: values.triggers || [],
        },
      });
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    if (processor) {
      form.setValues({
        name: processor.name,
        triggers: processor.triggers,
        enabled: processor.enabled,
        templateKey: processor.templateKey,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [processor]);

  if (error) {
    return <ErrorCard message={error.message} title={error.name} />;
  }

  return (
    <Box mb={50}>
      <Group mb="xl">
        <IconSettings size={24} />
        <Title order={3}>Settings</Title>
      </Group>
      <ContainerOverlay visible={isCreating || isPending} />
      <form onSubmit={form.onSubmit(handleUpdate)}>
        <TextInput
          mt="xl"
          withAsterisk
          label="Name"
          description="Name of the processor"
          type="name"
          placeholder="Enter a name"
          {...form.getInputProps('name')}
        />
        <TagsInput
          mt="xl"
          label="Trigger"
          description="File type to trigger the processor"
          placeholder="Enter comma separated file type"
          {...form.getInputProps('triggers')}
          min={0}
        />
        <Group mt={50}>
          <Button type="submit" variant="white">
            Update
          </Button>
        </Group>
      </form>
    </Box>
  );
};

export default ProcessorSettingsPage;
