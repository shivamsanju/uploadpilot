import {
    TextInput,
    Textarea,
    Select,
    Button,
    Group,
    Stack,
    TagsInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";

const CreateCodebaseForm = ({ onSubmit }: { onSubmit: (codebase: any) => void }) => {
    const form = useForm({
        initialValues: {
            name: "",
            description: "",
            tags: [] as string[],
            lang: "",
            basePath: ".",
            url: "",
        },

        validate: {
            name: (value) => (!value || value.trim().length === 0 ? "Name is required" : null),
            url: (value) => (!value || value.trim().length === 0 ? "Git URL is required" : null),
            description: (value) =>
                !value || value.trim().length === 0 ? "Description is required" : null,
            lang: (value) => (value ? null : "Please select a language"),
        },
    });

    const languages = [
        { value: "go", label: "Go" },
    ];

    return (
        <form
            onSubmit={form.onSubmit((values) => {
                const newCodebase = {
                    name: values.name,
                    description: values.description,
                    tags: values.tags,
                    lang: values.lang,
                    updated: new Date().toISOString(),
                    members: [],
                    basePath: values.basePath,
                    pinned: false,
                    gitUrl: values.url,
                };
                onSubmit(newCodebase);
            })}
        >
            <Stack gap="md">
                <TextInput
                    label="Name"
                    placeholder="Enter codebase name"
                    {...form.getInputProps("name")}
                />

                <TextInput
                    label="Git URL"
                    placeholder="Enter Git URL"
                    {...form.getInputProps("url")}
                />

                <TextInput
                    label="Base Path"
                    placeholder="Specify the base path of the package"
                    {...form.getInputProps("basePath")}
                />

                <Textarea
                    label="Description"
                    placeholder="Enter a brief description"
                    autosize
                    minRows={3}
                    {...form.getInputProps("description")}
                />

                <TagsInput
                    label="Press Enter to submit a tag"
                    placeholder="Enter tags"
                    {...form.getInputProps("tags")}
                />

                <Select
                    label="Language"
                    placeholder="Select a language"
                    data={languages}
                    {...form.getInputProps("lang")}
                />

                <Group p="right" mt="md">
                    <Button type="submit">Add Codebase</Button>
                </Group>
            </Stack>
        </form>
    );
};

export default CreateCodebaseForm;
