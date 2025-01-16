import React from "react";
import {
    Group,
    MultiSelect,
    NumberInput,
    Stack,
    Text,
    Switch,
    TagsInput,
    SegmentedControl,
} from "@mantine/core";
import { CreateUploaderForm } from "../../../../types/uploader";
import { UseFormReturnType } from "@mantine/form";
import { useGetAllAllowedSources } from "../../../../apis/uploader";
import ErrorCard from "../../../../components/ErrorCard/ErrorCard";
import { MIME_TYPES } from "../../../../utils/mime";
import classes from "./Form.module.css";

const w = "15vw";

type NewUploaderConfigProps = {
    form: UseFormReturnType<CreateUploaderForm, (values: CreateUploaderForm) => CreateUploaderForm>;
    type?: "create" | "edit" | "view";
};

const UploaderConfigForm: React.FC<NewUploaderConfigProps> = ({
    form,
    type = "create",
}) => {
    const { isPending, error, allowedSources } = useGetAllAllowedSources();

    if (error) {
        return <ErrorCard title={error.name} message={error.message} h="50vh" />;
    }

    return (
        <form>
            <Stack p="md">
                {/* Allowed file types */}
                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
                    <div>
                        <Text size="sm">Allowed file types *</Text>
                        <Text size="xs" c="dimmed">
                            Select allowed file types for the file uploader
                        </Text>
                    </div>
                    <MultiSelect
                        w={w}
                        size="xs"
                        data={MIME_TYPES}
                        {...form.getInputProps("allowedFileTypes")}
                        disabled={isPending || type === "view"}
                    />
                </Group>

                {/* Allowed input sources */}
                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
                    <div>
                        <Text size="sm">Allowed input sources *</Text>
                        <Text size="xs" c="dimmed">
                            Select allowed input sources for the file uploader
                        </Text>
                    </div>
                    <MultiSelect
                        w={w}
                        size="xs"
                        data={allowedSources}
                        {...form.getInputProps("allowedSources")}
                        disabled={isPending || type === "view"}
                    />
                </Group>

                {/* Min file size */}
                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
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
                        disabled={isPending || type === "view"}
                        min={0}
                    />
                </Group>

                {/* Max file size */}
                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
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
                        disabled={isPending || type === "view"}
                        min={0}
                    />
                </Group>

                {/* Min number of files */}
                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
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
                        disabled={isPending || type === "view"}
                        min={0}
                    />
                </Group>

                {/* Max number of files */}
                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
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
                        disabled={isPending || type === "view"}
                        min={1}
                    />
                </Group>

                {/* Max total file size */}
                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
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
                        disabled={isPending || type === "view"}
                        min={0}
                    />
                </Group>

                {/* Required metadata fields */}
                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
                    <div>
                        <Text size="sm">Required metadata fields</Text>
                        <Text size="xs" c="dimmed">
                            Enter the required metadata fields (separated by commas)
                        </Text>
                    </div>
                    <TagsInput
                        w={w}
                        size="xs"
                        {...form.getInputProps("requiredMetadataFields")}
                        disabled={isPending || type === "view"}
                        min={0}
                    />
                </Group>

                {/* Advanced options */}
                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
                    <div>
                        <Text size="sm">Show status bar</Text>
                        <Text size="xs" c="dimmed">
                            Toggle to show the status bar in the uploader
                        </Text>
                    </div>
                    <Switch
                        className={classes.switch}
                        size="lg"
                        onLabel="ON" offLabel="OFF"
                        checked={form.values.showStatusBar}
                        onChange={(e) => form.setFieldValue("showStatusBar", e.target.checked)}
                        disabled={isPending || type === "view"}
                    />
                </Group>

                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
                    <div>
                        <Text size="sm">Show progress bar</Text>
                        <Text size="xs" c="dimmed">
                            Toggle to show the progress bar in the uploader
                        </Text>
                    </div>
                    <Switch
                        className={classes.switch}
                        size="lg"
                        onLabel="ON" offLabel="OFF"
                        checked={form.values.showProgress}
                        onChange={(e) => form.setFieldValue("showProgress", e.target.checked)}
                        disabled={isPending || type === "view"}
                    />
                </Group>

                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
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
                        disabled={isPending || type === "view"}
                    />
                </Group>

                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
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
                        disabled={isPending || type === "view"}
                    />
                </Group>

                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
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
                        disabled={isPending || type === "view"}
                    />
                </Group>

                <Group justify="space-between" className={classes.item} wrap="nowrap" gap="xl">
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
                        disabled={isPending || type === "view"}
                    />
                </Group>
            </Stack>
        </form>
    );
};

export default UploaderConfigForm;
