import { useForm } from "@mantine/form";
import {
  Box,
  Button,
  Container,
  Group,
  MultiSelect,
  Paper,
  ScrollArea,
  SelectProps,
  SimpleGrid,
  Stepper,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { IconCheck, IconClockBolt, IconFile } from "@tabler/icons-react";
import { MIME_TYPES } from "../../../utils/mime";
import { MIME_TYPE_ICONS } from "../../../utils/fileicons";
import {
  useCreateProcessorMutation,
  useGetAllWorkflowTemplates,
} from "../../../apis/processors";
import { useNavigate, useParams } from "react-router-dom";
import { useState } from "react";
import { ContainerOverlay } from "../../../components/Overlay";
import { Template } from "./Template";
import { ErrorCard } from "../../../components/ErrorCard/ErrorCard";

const iconProps = {
  stroke: 1.5,
  opacity: 0.6,
  size: 14,
};

const renderSelectOption: SelectProps["renderOption"] = ({
  option,
  checked,
}) => {
  let Icon = MIME_TYPE_ICONS[option.value];
  if (!Icon) {
    Icon = IconFile;
  }
  return (
    <Group flex="1" gap="xs">
      <Icon />
      {option.label}
      {/* <Text c="dimmed" ml="sm">({option.value})</Text> */}
      {checked && (
        <IconCheck style={{ marginInlineStart: "auto" }} {...iconProps} />
      )}
    </Group>
  );
};

const NewprocessorPage = () => {
  const [active, setActive] = useState(0);
  const { workspaceId } = useParams();
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

  const { mutateAsync: createMutateAsync, isPending: isCreating } =
    useCreateProcessorMutation();
  const { isPending, templates, error } = useGetAllWorkflowTemplates(
    workspaceId || ""
  );

  const handleAdd = async (values: any) => {
    try {
      await createMutateAsync({
        workspaceId: workspaceId || "",
        processor: {
          workspaceId,
          name: values.name,
          triggers: values.triggers,
          enabled: values.enabled,
          templateKey: values.templateKey,
        },
      });
      form.reset();
      navigate(`/workspace/${workspaceId}/processors`);
    } catch (error) {
      console.error(error);
    }
  };

  const nextStep = () => {
    const val = form.validate();
    if (val.hasErrors) {
      return;
    }
    setActive((current) => (current < 3 ? current + 1 : current));
  };

  const prevStep = () => {
    setActive((current) => (current > 0 ? current - 1 : current));
  };

  const selectTemplate = (templateKey: string) => {
    if (templateKey === form.values.templateKey) {
      form.setFieldValue("templateKey", "");
    } else {
      form.setFieldValue("templateKey", templateKey);
    }
  };

  return (
    <Box mb={50}>
      <Container>
        <Group justify="center" mb="xl">
          <Title order={3} opacity={0.7}>
            New Processor
          </Title>
        </Group>
        <Paper p="xl" withBorder>
          <form onSubmit={form.onSubmit(handleAdd)}>
            <Stepper size="xs" active={active}>
              <Stepper.Step
                label="Name"
                description="Choose a name for your processor"
              >
                <TextInput
                  mt="xl"
                  withAsterisk
                  label="Name"
                  description="Name of the processor"
                  type="name"
                  placeholder="Enter a name"
                  {...form.getInputProps("name")}
                />
              </Stepper.Step>
              <Stepper.Step label="Trigger" description="Select triggers">
                <MultiSelect
                  mt="xl"
                  searchable
                  leftSection={<IconClockBolt size={16} />}
                  label="Trigger"
                  description="File type to trigger the processor"
                  placeholder="Select file type"
                  data={MIME_TYPES}
                  {...form.getInputProps("triggers")}
                  renderOption={renderSelectOption}
                />
              </Stepper.Step>
              <Stepper.Step label="Workflow" description="Select a workflow">
                <Text size="sm" mb="sm">
                  Choose workflow from an prebuilt templates
                </Text>
                <ScrollArea scrollbarSize={6} h="55vh">
                  {error && (
                    <ErrorCard message={error?.message} title={"Error"} />
                  )}
                  <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }} spacing="xl">
                    {templates &&
                      templates.length > 0 &&
                      templates.map((t: any) => (
                        <Box onClick={() => selectTemplate(t.key)}>
                          <Template
                            template={t}
                            selected={t.key === form.values.templateKey}
                          />
                        </Box>
                      ))}
                  </SimpleGrid>
                </ScrollArea>
              </Stepper.Step>
            </Stepper>
            <ContainerOverlay visible={isCreating || isPending} />
            <Group justify="flex-end" mt={50}>
              <Button
                variant="default"
                onClick={prevStep}
                disabled={active === 0}
              >
                Back
              </Button>
              {active === 2 && <Button type="submit">Create</Button>}
              {active < 2 && <Button onClick={nextStep}>Next</Button>}
            </Group>
          </form>
        </Paper>
      </Container>
    </Box>
  );
};

export default NewprocessorPage;
