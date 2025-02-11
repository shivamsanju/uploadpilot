import {
  Group,
  MultiSelect,
  NumberInput,
  Stack,
  Text,
  Switch,
  TagsInput,
  SimpleGrid,
  Transition,
  Paper,
  TextInput,
  Tooltip,
  SelectProps,
  LoadingOverlay,
} from "@mantine/core";
import { UploaderConfig } from "../../types/uploader";
import { useForm } from "@mantine/form";
import { MIME_TYPES } from "../../utils/mime";
import classes from "./Form.module.css";
import { useParams } from "react-router-dom";
import { useGetAllAllowedSources } from "../../apis/workspace";
import { IconCheck, IconFile, IconInfoCircle } from "@tabler/icons-react";
import { useUpdateUploaderConfigMutation } from "../../apis/uploader";
import { showNotification } from "@mantine/notifications";
import { MIME_TYPE_ICONS } from "../../utils/fileicons";
import { DiscardButton } from "../../components/Buttons/DiscardButton";
import { SaveButton } from "../../components/Buttons/SaveButton";

const authEndpointTooltip = `
If you have a custom authentication endpoint, enter it here.\n
 We will send a request with all headers you set in uploader to this endpoint for authentication.\n
You can use this to authenticate the upload by setting your token in the Authorization header.\n
You can leave this field empty if you don't have a custom authentication endpoint.
`;

type NewUploaderConfigProps = {
  config: UploaderConfig;
};

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

const UploaderConfigForm: React.FC<NewUploaderConfigProps> = ({ config }) => {
  const { workspaceId } = useParams();
  const { isPending, allowedSources } = useGetAllAllowedSources(
    workspaceId || ""
  );
  const { mutateAsync, isPending: isPendingMutation } =
    useUpdateUploaderConfigMutation();

  const form = useForm<UploaderConfig>({
    initialValues: {
      ...config,
      allowedFileTypes: config?.allowedFileTypes || [],
      allowedSources: config?.allowedSources || [],
      requiredMetadataFields: config?.requiredMetadataFields || [],
      authEndpoint: config?.authEndpoint || "",
    },
    validate: {
      allowedSources: (value) =>
        value.length === 0 ? "Please select at least one source" : null,
      authEndpoint: (value) =>
        value && !/^https?:\/\//.test(value)
          ? "Please enter a valid URL"
          : null,
    },
  });

  const handleEditAndSaveButton = async () => {
    if (!workspaceId) {
      showNotification({
        color: "red",
        title: "Error",
        message: "Workspace ID is not available",
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
      }).catch((error) => {
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
      <LoadingOverlay
        visible={isPending || isPendingMutation}
        overlayProps={{ backgroundOpacity: 0 }}
        zIndex={1000}
      />
      <Paper withBorder p="sm" mb={50}>
        <SimpleGrid cols={{ base: 1, xl: 2 }}>
          <Stack p="md">
            {/* Allowed input sources */}
            <MultiSelect
              label="Allowed input sources"
              description="Allowed input sources for your uploader"
              data={allowedSources || []}
              {...form.getInputProps("allowedSources")}
              disabled={isPending}
              searchable
            />

            {/* Allowed file types */}
            <MultiSelect
              label="Allowed file types"
              description="Allowed file types for your uploader"
              data={MIME_TYPES}
              {...form.getInputProps("allowedFileTypes")}
              searchable
              renderOption={renderSelectOption}
            />

            {/*Auth Endpoint */}
            <TextInput
              label={
                <Text c="dimmed">
                  Enter a auth endpoint
                  <Tooltip
                    w="300px"
                    multiline
                    transitionProps={{ duration: 200 }}
                    label={authEndpointTooltip}
                  >
                    <IconInfoCircle
                      size={14}
                      style={{ cursor: "pointer", marginLeft: "5px" }}
                    />
                  </Tooltip>
                </Text>
              }
              description="Enter an auth endpoint"
              type="url"
              {...form.getInputProps("authEndpoint")}
              min={0}
            />
            <TagsInput
              label="Required metadata fields"
              description="Required metadata fields for your uploader"
              {...form.getInputProps("requiredMetadataFields")}
              min={0}
            />
          </Stack>
          <Stack p="md">
            {/* Min file size */}
            <NumberInput
              label="Min file size"
              description="Enter minimum file size in bytes"
              {...form.getInputProps("minFileSize")}
              min={0}
            />

            {/* Max file size */}
            <NumberInput
              label="Max file size"
              description="Enter maximum file size in bytes"
              {...form.getInputProps("maxFileSize")}
              min={0}
            />

            {/* Min number of files */}
            <NumberInput
              label="Min number of files"
              description="Specify the minimum number of files required"
              {...form.getInputProps("minNumberOfFiles")}
              min={0}
            />

            {/* Max number of files */}
            <NumberInput
              label="Max number of files"
              description="Specify the maximum number of files allowed"
              {...form.getInputProps("maxNumberOfFiles")}
              min={1}
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
                onChange={(e) =>
                  form.setFieldValue("allowPauseAndResume", e.target.checked)
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
                onChange={(e) =>
                  form.setFieldValue("enableImageEditing", e.target.checked)
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
                onChange={(e) =>
                  form.setFieldValue("useCompression", e.target.checked)
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
                onChange={(e) =>
                  form.setFieldValue("useFaultTolerantMode", e.target.checked)
                }
              />
            </Group>
          </Stack>
        </SimpleGrid>
      </Paper>
      <Transition
        mounted={form.isDirty()}
        transition="fade-up"
        duration={400}
        timingFunction="ease"
      >
        {(styles) => (
          <div style={styles}>
            <Group justify="center" gap="md" mt="xl">
              <DiscardButton type="reset" />
              <SaveButton type="submit" />
            </Group>
          </div>
        )}
      </Transition>
    </form>
  );
};

export default UploaderConfigForm;
