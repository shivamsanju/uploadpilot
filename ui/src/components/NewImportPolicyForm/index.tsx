import {
    TextInput,
    MultiSelect,
    Stack,
    Button,
    Group,
    NumberInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";

export type ImportPolicy = {
    name: string;
    allowedMimeTypes: string[];
    allowedSources: string[];
    maxFileSize: number;
    maxFileCount: number;
};

type ImportPolicyFormProps = {
    onSubmit: (policy: ImportPolicy) => void;
    showSubmitButton?: boolean;
    showCancelButton?: boolean;
    onCancel?: () => void;
};

const ImportPolicyForm = ({
    onSubmit,
    showSubmitButton,
    showCancelButton,
    onCancel,
}: ImportPolicyFormProps) => {
    const form = useForm<ImportPolicy>({
        initialValues: {
            name: "",
            allowedMimeTypes: [],
            allowedSources: [],
            maxFileSize: 0,
            maxFileCount: 0,
        },
        validate: {
            name: (value) => (!value.trim() ? "Name is required" : null),
            maxFileSize: (value) =>
                value <= 0 ? "Maximum file size must be greater than 0" : null,
            maxFileCount: (value) =>
                value <= 0 ? "Maximum file count must be greater than 0" : null,
        },
    });

    return (
        <form
            onSubmit={form.onSubmit((values) => {
                console.log(values);
                onSubmit(values);
            })}
        >
            <Stack gap="md" h="50vh">
                <TextInput
                    label="Policy Name"
                    placeholder="Enter policy name"
                    withAsterisk
                    {...form.getInputProps("name")}
                />

                <MultiSelect
                    label="Allowed MIME Types"
                    placeholder="Select allowed MIME types"
                    data={[
                        "image/jpeg",
                        "image/png",
                        "application/pdf",
                        "text/plain",
                    ]}
                    {...form.getInputProps("allowedMimeTypes")}
                />

                <MultiSelect
                    label="Allowed Sources"
                    placeholder="Select allowed sources"
                    data={[
                        "Google Drive",
                        "Dropbox",
                        "Local Upload",
                        "OneDrive",
                    ]}
                    {...form.getInputProps("allowedSources")}
                />

                <NumberInput
                    label="Maximum File Size (in MB)"
                    placeholder="Enter max file size"
                    withAsterisk
                    {...form.getInputProps("maxFileSize")}
                />

                <NumberInput
                    label="Maximum File Count"
                    placeholder="Enter max file count"
                    withAsterisk
                    {...form.getInputProps("maxFileCount")}
                />
            </Stack>

            <Group justify="flex-end" mt="md" gap="sm">
                {showCancelButton && <Button onClick={onCancel}>Cancel</Button>}
                {showSubmitButton && <Button type="submit">Submit</Button>}
            </Group>
        </form>
    );
};

export default ImportPolicyForm;
