import { useForm } from "@mantine/form";
import {
  ActionIcon,
  Anchor,
  Box,
  Breadcrumbs,
  Button,
  Group,
  Paper,
  TagsInput,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { IconChevronLeft } from "@tabler/icons-react";
import {
  useUpdateProcessorMutation,
  useGetProcessor,
} from "../../../apis/processors";
import { useNavigate, useParams } from "react-router-dom";
import { ContainerOverlay } from "../../../components/Overlay";
import { ErrorCard } from "../../../components/ErrorCard/ErrorCard";
import { useEffect } from "react";

const ProcessorSettingsPage = () => {
  const { workspaceId, processorId } = useParams();
  const navigate = useNavigate();

  const form = useForm({
    initialValues: {
      name: "",
      triggers: [],
      enabled: true,
      templateKey: "",
    },
    validate: {
      name: (value) => (value ? null : "Name is required"),
    },
  });

  const { mutateAsync, isPending: isCreating } = useUpdateProcessorMutation();
  const { isPending, error, processor } = useGetProcessor(
    workspaceId || "",
    processorId || ""
  );

  const handleUpdate = async (values: any) => {
    try {
      await mutateAsync({
        workspaceId: workspaceId || "",
        processorId: processorId || "",
        processor: {
          workspaceId,
          name: values.name,
          triggers: values.triggers,
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
      <Breadcrumbs separator=">">
        <Anchor href={`/`}>Workspaces</Anchor>
        <Anchor href={`/workspace/${workspaceId}/processors`}>
          Processors
        </Anchor>
        <Text>{processorId}</Text>
      </Breadcrumbs>
      <Group mt="xs" mb="xl">
        <ActionIcon
          variant="default"
          radius="xl"
          size="sm"
          onClick={() => navigate(`/workspace/${workspaceId}/processors`)}
        >
          <IconChevronLeft size={16} />
        </ActionIcon>
        <Title order={3}>Settings</Title>
      </Group>
      <Paper p="xl" withBorder>
        <ContainerOverlay visible={isCreating || isPending} />
        <form onSubmit={form.onSubmit(handleUpdate)}>
          <TextInput
            mt="xl"
            withAsterisk
            label="Name"
            description="Name of the processor"
            type="name"
            placeholder="Enter a name"
            {...form.getInputProps("name")}
          />
          <TagsInput
            mt="xl"
            label="Trigger"
            description="File type to trigger the processor"
            placeholder="Enter comma separated file type"
            {...form.getInputProps("triggers")}
            min={0}
          />
          <Group justify="flex-end" mt={50}>
            <Button type="submit">Update</Button>
          </Group>
        </form>
      </Paper>
    </Box>
  );
};

export default ProcessorSettingsPage;
