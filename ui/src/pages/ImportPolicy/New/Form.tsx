import React from "react";
import {
    TextInput,
    MultiSelect,
    NumberInput,
    Button,
    Group,
    Stack,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { ImportPolicy } from "../../../types/importpolicy";

const allowedMimeTypesOptions = [
    "image/jpeg",
    "image/png",
    "application/pdf",
    "text/plain",
];

const allowedSourcesOptions = ["Upload", "Google Drive", "Dropbox", "OneDrive"];

type NewImportPolicyFormProps = {
    type?: "create" | "edit" | "view";
    initialValues?: ImportPolicy;
    onSubmit: (values: ImportPolicy) => void;
    onCancel?: () => void;
    showCancelButton?: boolean;
    showSubmitButton?: boolean;
}
const ImportPolicyForm: React.FC<NewImportPolicyFormProps> = ({
    type = "create",
    initialValues,
    onSubmit,
    onCancel,
    showCancelButton,
    showSubmitButton,
}) => {
    const form = useForm<ImportPolicy>({
        initialValues: initialValues || {
            name: "",
            allowedMimeTypes: [],
            allowedSources: [],
            maxFileSizeKb: 0,
            maxFileCount: 0,
        },
        validate: {
            name: (value) => (value.trim() === "" ? "Name is required" : null),
            maxFileSizeKb: (value) =>
                value <= 0 ? "Max file size must be greater than 0" : null,
            maxFileCount: (value) =>
                value <= 0 ? "Max file count must be greater than 0" : null,
        },
    });

    return (
        <>
            <form>
                <Stack gap="md">
                    <TextInput

                        label="Policy Name"
                        placeholder="Enter policy name"
                        {...form.getInputProps("name")}
                        required
                        disabled={type === "view"}
                    />

                    <MultiSelect
                        size="xs"
                        label="Allowed MIME Types"
                        placeholder="Select allowed MIME types"
                        data={allowedMimeTypesOptions}
                        {...form.getInputProps("allowedMimeTypes")}
                        withAsterisk
                        disabled={type === "view"}
                    />

                    <MultiSelect
                        size="xs"
                        label="Allowed Sources"
                        placeholder="Select allowed sources"
                        data={allowedSourcesOptions}
                        {...form.getInputProps("allowedSources")}
                        withAsterisk
                        disabled={type === "view"}
                    />

                    <NumberInput
                        size="xs"
                        label="Max File Size (in Kb)"
                        placeholder="Enter max file size"
                        {...form.getInputProps("maxFileSizeKb")}
                        withAsterisk
                        disabled={type === "view"}
                    />

                    <NumberInput
                        size="xs"
                        label="Max File Count"
                        placeholder="Enter max file count"
                        {...form.getInputProps("maxFileCount")}
                        withAsterisk
                        disabled={type === "view"}
                    />

                    {type !== "view" && <Group justify="flex-end" mt="sm" gap="md">
                        {showCancelButton && <Button onClick={onCancel}>Cancel</Button>}
                        {showSubmitButton && <Button onClick={() => {
                            if (form.validate().hasErrors) {
                                return;
                            };
                            onSubmit(form.values);
                        }}>Create</Button>}
                    </Group>}
                </Stack>
            </form>
        </>
    );
};

export default ImportPolicyForm;
