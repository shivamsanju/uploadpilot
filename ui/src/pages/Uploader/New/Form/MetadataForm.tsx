import { TextInput, Textarea, Stack, TagsInput } from "@mantine/core";
import { CreateUploaderForm } from "../../../../types/uploader";
import { UseFormReturnType } from "@mantine/form";

type UploaderMetadataFormProps = {
    form: UseFormReturnType<CreateUploaderForm, (values: CreateUploaderForm) => CreateUploaderForm>;
}

const UploaderMetadataForm: React.FC<UploaderMetadataFormProps> = ({ form }) => (
    <Stack gap="md">
        <TextInput
            withAsterisk
            label="Uploader Name"
            placeholder="Give a name to your uploader"
            {...form.getInputProps("name")}
        />
        <Textarea
            withAsterisk
            label="Description"
            placeholder="Give a brief description of your uploader"
            autosize
            minRows={6}
            {...form.getInputProps("description")}
        />
        <TagsInput
            label="Tags"
            placeholder="Enter tags(comma-separated)"
            {...form.getInputProps("tags")}
        />
    </Stack>
);

export default UploaderMetadataForm;
