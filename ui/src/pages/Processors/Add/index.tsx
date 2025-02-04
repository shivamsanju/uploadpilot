import { useForm } from "@mantine/form";
import {
  Button,
  Group,
  LoadingOverlay,
  MultiSelect,
  SelectProps,
  Stack,
  TextInput,
} from "@mantine/core";
import { IconCheck, IconClockBolt, IconFile } from "@tabler/icons-react";
import { MIME_TYPES } from "../../../utils/mime";
import { MIME_TYPE_ICONS } from "../../../utils/fileicons";
import {
  useCreateProcessorMutation,
  useUpdateProcessorMutation,
} from "../../../apis/processors";

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

type Props = {
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
  workspaceId: string;
  mode?: "edit" | "add" | "view";
  setMode: React.Dispatch<React.SetStateAction<"edit" | "add" | "view">>;
  initialValues?: any;
  setInitialValues: React.Dispatch<React.SetStateAction<any>>;
};

const AddWebhookForm: React.FC<Props> = ({
  setInitialValues,
  setMode,
  setOpened,
  workspaceId,
  mode = "add",
  initialValues,
}) => {
  const form = useForm({
    initialValues:
      (mode === "edit" || mode === "view") && initialValues
        ? initialValues
        : {
            name: "",
            triggers: [],
            tasks: {
              nodes: [],
              edges: [],
            },
            enabled: true,
          },
    validate: {
      name: (value) => (value ? null : "Name is required"),
    },
  });

  const { mutateAsync: createMutateAsync, isPending: isCreating } =
    useCreateProcessorMutation();
  const { mutateAsync: updateMutateAsync, isPending: isUpdating } =
    useUpdateProcessorMutation();

  const handleAdd = async (values: any) => {
    try {
      if (mode === "edit") {
        await updateMutateAsync({
          processorId: initialValues.id,
          workspaceId,
          processor: {
            name: values.name,
            triggers: values.triggers,
          },
        });
      } else {
        await createMutateAsync({
          workspaceId,
          processor: {
            workspaceId,
            name: values.name,
            triggers: values.triggers,
            tasks: values.tasks,
            enabled: values.enabled,
            data: {},
          },
        });
      }
      setOpened(false);
      setInitialValues(null);
      setMode("add");
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <form onSubmit={form.onSubmit(handleAdd)}>
      <LoadingOverlay
        visible={isCreating || isUpdating}
        overlayProps={{ backgroundOpacity: 0 }}
        zIndex={1000}
      />
      <Stack gap="xl">
        <TextInput
          withAsterisk
          label="Name"
          description="Name of the processor"
          type="name"
          placeholder="Enter a name"
          {...form.getInputProps("name")}
          disabled={mode === "view"}
        />
        <MultiSelect
          searchable
          leftSection={<IconClockBolt size={16} />}
          label="Trigger"
          description="File type to trigger the processor"
          placeholder="Select file type"
          data={MIME_TYPES}
          {...form.getInputProps("triggers")}
          renderOption={renderSelectOption}
          disabled={mode === "view"}
        />
      </Stack>
      {mode !== "view" && (
        <Group justify="flex-end" mt={50}>
          <Button type="submit">{mode === "edit" ? "Save" : "Create"}</Button>
        </Group>
      )}
    </form>
  );
};

export default AddWebhookForm;
