import { TextInput, Textarea, Stack, TagsInput } from "@mantine/core";

const WorkflowMetadataForm = ({ form }: { form: any }) => (
    <Stack gap="md">
        <TextInput
            label="Workflow Name"
            placeholder="Enter workflow name"
            {...form.getInputProps("name")}
        />
        <Textarea
            label="Description"
            placeholder="Enter a brief description"
            autosize
            minRows={4}
            {...form.getInputProps("description")}
        />
        <TagsInput
            label="Tags"
            placeholder="Enter tags (comma-separated)"
            {...form.getInputProps("tags")}
        />
    </Stack>
);

export default WorkflowMetadataForm;
