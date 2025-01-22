import React from "react";
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
} from "@mantine/core";
import { UploaderConfig } from "../../../types/uploader";
import { useForm } from "@mantine/form";
import { MIME_TYPES } from "../../../utils/mime";
import classes from "./Form.module.css";
import { useParams } from "react-router-dom";
import { useGetAllAllowedSources } from "../../../apis/workspace";
import { ErrorLoadingWrapper } from "../../../components/ErrorLoadingWrapper";
import { IconDeviceFloppy, IconRestore } from "@tabler/icons-react";
import { useUpdateUploaderConfigMutation } from "../../../apis/uploader";
import { showNotification } from "@mantine/notifications";

const w = "300px";

type NewUploaderConfigProps = {
    config: UploaderConfig;
};

const UploaderConfigForm: React.FC<NewUploaderConfigProps> = ({
    config,
}) => {
    const { workspaceId } = useParams();
    const { isPending, error, allowedSources } = useGetAllAllowedSources(workspaceId || "");
    const { mutateAsync } = useUpdateUploaderConfigMutation()


    const form = useForm<UploaderConfig>({
        initialValues: {
            ...config,
            allowedFileTypes: config?.allowedFileTypes || [],
            allowedSources: config?.allowedSources || [],
            requiredMetadataFields: config?.requiredMetadataFields || []
        },
        validate: {
            allowedSources: (value) => value.length === 0 ? "Please select at least one source" : null,
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
                config: form.values
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

                            {/* Max total file size */}
                            {/* <Group justify="space-between" className={classes.item}>
                                <div>
                                    <Text size="sm">Max total file size (bytes)</Text>
                                    <Text size="xs" c="dimmed">
                                        Enter the maximum total file size allowed
                                    </Text>
                                </div>
                                <NumberInput
                                    w={w}
                                    size="xs"
                                    {...form.getInputProps("maxTotalFileSize")}
                                    // disabled={isPending || type === "view"}
                                    min={0}
                                />
                            </Group> */}
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
                                        Toggle to enable file compression
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
