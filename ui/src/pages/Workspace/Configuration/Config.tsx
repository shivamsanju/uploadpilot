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
    Button,
    Paper,
    TextInput,
    Tooltip,
    LoadingOverlay,
} from "@mantine/core";
import { UploaderConfig } from "../../../types/uploader";
import { useForm } from "@mantine/form";
import { MIME_TYPES } from "../../../utils/mime";
import classes from "./Form.module.css";
import { useParams } from "react-router-dom";
import { useGetAllAllowedSources } from "../../../apis/workspace";
import { ErrorLoadingWrapper } from "../../../components/ErrorLoadingWrapper";
import { IconDeviceFloppy, IconInfoCircle, IconRestore } from "@tabler/icons-react";
import { useUpdateUploaderConfigMutation } from "../../../apis/uploader";
import { showNotification } from "@mantine/notifications";

const w = "300px";
const authEndpointTooltip = `
If you have a custom authentication endpoint, enter it here.\n
 We will send a request with all headers you set in uploader to this endpoint for authentication.\n
You can use this to authenticate the upload by setting your token in the Authorization header.\n
You can leave this field empty if you don't have a custom authentication endpoint.
`;

type NewUploaderConfigProps = {
    config: UploaderConfig;
};

const UploaderConfigForm: React.FC<NewUploaderConfigProps> = ({
    config,
}) => {
    const { workspaceId } = useParams();
    const { isPending, error, allowedSources } = useGetAllAllowedSources(workspaceId || "");
    const { mutateAsync, isPending: isPendingMutation } = useUpdateUploaderConfigMutation()


    const form = useForm<UploaderConfig>({
        initialValues: {
            ...config,
            allowedFileTypes: config?.allowedFileTypes || [],
            allowedSources: config?.allowedSources || [],
            requiredMetadataFields: config?.requiredMetadataFields || [],
            authEndpoint: config?.authEndpoint || "",
        },
        validate: {
            allowedSources: (value) => value.length === 0 ? "Please select at least one source" : null,
            authEndpoint: (value) => (value && !/^https?:\/\//.test(value)) ? "Please enter a valid URL" : null
        }
    });

    const handleEditAndSaveButton = async () => {
        if (!workspaceId) {
            showNotification({
                color: "red",
                title: "Error",
                message: "Workspace ID is not available"
            })
            return
        };
        if (form.isDirty()) {
            mutateAsync({
                workspaceId: workspaceId,
                config: {
                    ...form.values,
                    minFileSize: form.values.minFileSize || 0,
                    maxFileSize: form.values.maxFileSize || 0,
                    minNumberOfFiles: form.values.minNumberOfFiles || 0,
                    maxNumberOfFiles: form.values.maxNumberOfFiles || 0,
                }
            }).catch((error) => {
                console.log(error)
            })
        }
        form.resetDirty();
    }

    const handleResetButton = () => {
        form.reset();
        form.resetDirty();
    }


    return (
        <ErrorLoadingWrapper error={error} isPending={isPending}>
            <form onSubmit={form.onSubmit(handleEditAndSaveButton)} onReset={handleResetButton}>
                <LoadingOverlay visible={isPendingMutation} />
                <Paper p="sm">
                    <SimpleGrid cols={{ base: 1, xl: 2 }}>
                        <Stack p="md">

                            {/* Allowed input sources */}
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Allowed input sources *</Text>
                                    <Text size="xs" c="dimmed">
                                        Allowed input sources for your uploader
                                    </Text>
                                </div>
                                <MultiSelect
                                    w={w}
                                    size="xs"
                                    data={allowedSources || []}
                                    {...form.getInputProps("allowedSources")}
                                    disabled={isPending}
                                // disabled={isPending || type === "view"}
                                />
                            </Group>

                            {/* Allowed file types */}
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Allowed file types</Text>
                                    <Text size="xs" c="dimmed">
                                        Select allowed file types for the file uploader
                                    </Text>
                                </div>
                                <MultiSelect
                                    w={w}
                                    size="xs"
                                    data={MIME_TYPES}
                                    {...form.getInputProps("allowedFileTypes")}
                                // disabled={isPending || type === "view"}
                                />
                            </Group>

                            {/*Auth Endpoint */}
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Auth Endpoint</Text>
                                    <Text size="xs" c="dimmed">
                                        Enter a auth endpoint{"  "}
                                        <Tooltip
                                            w="300px"
                                            multiline
                                            transitionProps={{ duration: 200 }}
                                            label={authEndpointTooltip}
                                        >
                                            <IconInfoCircle size={14} />
                                        </Tooltip>
                                    </Text>
                                </div>
                                <TextInput
                                    w={w}
                                    type="url"
                                    size="xs"
                                    {...form.getInputProps("authEndpoint")}
                                    // disabled={isPending || type === "view"}
                                    min={0}
                                />
                            </Group>
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Required metadata fields</Text>
                                    <Text size="xs" c="dimmed">
                                        Separated by commas
                                    </Text>
                                </div>
                                <TagsInput
                                    w={w}
                                    size="xs"
                                    {...form.getInputProps("requiredMetadataFields")}
                                    // disabled={isPending || type === "view"}
                                    min={0}
                                />
                            </Group>
                        </Stack>
                        <Stack p="md">
                            {/* Min file size */}
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Min file size (bytes)</Text>
                                    <Text size="xs" c="dimmed">
                                        Enter minimum file size in bytes
                                    </Text>
                                </div>
                                <NumberInput
                                    w={w}
                                    size="xs"
                                    {...form.getInputProps("minFileSize")}
                                    // disabled={isPending || type === "view"}
                                    min={0}
                                />
                            </Group>

                            {/* Max file size */}
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Max file size (bytes)</Text>
                                    <Text size="xs" c="dimmed">
                                        Enter maximum file size in bytes
                                    </Text>
                                </div>
                                <NumberInput
                                    w={w}
                                    size="xs"
                                    {...form.getInputProps("maxFileSize")}
                                    // disabled={isPending || type === "view"}
                                    min={0}
                                />
                            </Group>

                            {/* Min number of files */}
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Min number of files</Text>
                                    <Text size="xs" c="dimmed">
                                        Specify the minimum number of files required
                                    </Text>
                                </div>
                                <NumberInput
                                    w={w}
                                    size="xs"
                                    {...form.getInputProps("minNumberOfFiles")}
                                    // disabled={isPending || type === "view"}
                                    min={0}
                                />
                            </Group>

                            {/* Max number of files */}
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Max number of files</Text>
                                    <Text size="xs" c="dimmed">
                                        Specify the maximum number of files allowed
                                    </Text>
                                </div>
                                <NumberInput
                                    w={w}
                                    size="xs"
                                    {...form.getInputProps("maxNumberOfFiles")}
                                    // disabled={isPending || type === "view"}
                                    min={1}
                                />
                            </Group>
                        </Stack>
                    </SimpleGrid>
                </Paper>
                <Paper p="sm" mt={50}>
                    <SimpleGrid cols={{ sm: 1, lg: 2 }}>
                        <Stack p="md">
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Allow pause and resume</Text>
                                    <Text size="xs" c="dimmed">
                                        Toggle to allow pause and resume in the uploader
                                    </Text>
                                </div>
                                <Switch
                                    className={classes.switch}
                                    size="lg"
                                    onLabel="ON" offLabel="OFF"
                                    checked={form.values.allowPauseAndResume}
                                    onChange={(e) => form.setFieldValue("allowPauseAndResume", e.target.checked)}
                                // disabled={isPending || type === "view"}
                                />
                            </Group>

                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Enable image editing</Text>
                                    <Text size="xs" c="dimmed">
                                        Toggle to enable image editing in the uploader ui
                                    </Text>
                                </div>
                                <Switch
                                    className={classes.switch}
                                    size="lg"
                                    onLabel="ON" offLabel="OFF"
                                    checked={form.values.enableImageEditing}
                                    onChange={(e) => form.setFieldValue("enableImageEditing", e.target.checked)}
                                // disabled={isPending || type === "view"}
                                />
                            </Group>
                        </Stack>
                        <Stack p="md">
                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Use compression</Text>
                                    <Text size="xs" c="dimmed">
                                        Toggle to enable compression while uploading files
                                    </Text>
                                </div>
                                <Switch
                                    className={classes.switch}
                                    size="lg"
                                    onLabel="ON" offLabel="OFF"
                                    checked={form.values.useCompression}
                                    onChange={(e) => form.setFieldValue("useCompression", e.target.checked)}
                                // disabled={isPending || type === "view"}
                                />
                            </Group>

                            <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Use fault tolerant mode</Text>
                                    <Text size="xs" c="dimmed">
                                        Fault tolerant mode allows to recover from browser crashes
                                    </Text>
                                </div>
                                <Switch
                                    className={classes.switch}
                                    size="lg"
                                    onLabel="ON" offLabel="OFF"
                                    checked={form.values.useFaultTolerantMode}
                                    onChange={(e) => form.setFieldValue("useFaultTolerantMode", e.target.checked)}
                                // disabled={isPending || type === "view"}
                                />
                            </Group>
                        </Stack>
                    </SimpleGrid >
                </Paper>
                <Transition
                    mounted={form.isDirty()}
                    transition="fade-up"
                    duration={400}
                    timingFunction="ease"
                >
                    {(styles) => <div style={styles}>
                        <Group justify="center" gap="md" mt="xl">
                            <Button
                                variant="light"
                                type="reset"
                                leftSection={<IconRestore size={18} />}
                            >
                                Reset
                            </Button>
                            <Button
                                type="submit"
                                leftSection={<IconDeviceFloppy size={18} />}
                            >
                                Save
                            </Button>
                        </Group>
                    </div>}
                </Transition>
            </form >
        </ErrorLoadingWrapper >
    );
};

export default UploaderConfigForm;
