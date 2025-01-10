import {
    TextInput,
    Select,
    Stack,
    Button,
    Group,
    TagsInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { Connector } from "../../../../types/connector";


type NewConnectorFormProps = {
    onSubmit: (connector: Connector) => void
    showSubmitButton?: boolean
    showCancelButton?: boolean
    onCancel?: () => void
}

const NewConnectorForm = ({ onSubmit, showSubmitButton, showCancelButton, onCancel }: NewConnectorFormProps) => {
    const form = useForm<Connector>({
        initialValues: {
            name: "",
            type: "s3",
            tags: [],
            s3Region: "",
            s3AccessKey: "",
            s3SecretKey: "",
            gcsApiKey: "",
            azureAccountName: "",
            azureAccountKey: "",

        },
        validate: {
            name: (value) => (!value || value.trim().length === 0 ? "Name is required" : null),
            type: (value) => (value ? null : "Please select a type"),
            s3Region: (value, values) => (values.type === "s3" && !value ? "Please enter a region" : null),
            s3AccessKey: (value, values) => (values.type === "s3" && !value ? "Please enter an access key" : null),
            s3SecretKey: (value, values) => (values.type === "s3" && !value ? "Please enter a secret key" : null),
            gcsApiKey: (value, values) => (values.type === "gcs" && !value ? "Please enter an API key" : null),
            azureAccountName: (value, values) => (values.type === "azure" && !value ? "Please enter an account name" : null),
            azureAccountKey: (value, values) => (values.type === "azure" && !value ? "Please enter an account key" : null),
        },
    });


    return (
        <form onSubmit={form.onSubmit((values) => {
            onSubmit(values)
        })}>
            <Stack gap="md" h="50vh">

                <TextInput label="Name" placeholder="Enter Name" withAsterisk {...form.getInputProps("name")} />
                <TagsInput
                    label="Press Enter to submit a tag"
                    placeholder="Enter tags"
                    {...form.getInputProps("tags")}
                />
                <Select
                    label="Connector Type"
                    placeholder="Select Connector Type"
                    withAsterisk
                    data={[
                        { value: "s3", label: "AWS S3" },
                        { value: "gcs", label: "Google Cloud Storage" },
                        { value: "azure", label: "Azure Blob Storage" },
                    ]}
                    {...form.getInputProps("type")}
                />
                {
                    form.values.type === "s3" && (
                        <>
                            <TextInput
                                withAsterisk
                                label="S3 Access Key"
                                placeholder="Enter S3 Access Key"
                                {...form.getInputProps("s3AccessKey")}
                            />
                            <TextInput
                                withAsterisk

                                label="S3 Secret Key"
                                placeholder="Enter S3 Access Key"
                                {...form.getInputProps("s3SecretKey")}
                            />
                            <Select
                                withAsterisk
                                label="S3 Region"
                                placeholder="Select Region"
                                data={[
                                    { value: "us-east-1", label: "US East (N. Virginia)" },
                                    { value: "us-east-2", label: "US East (Ohio)" },
                                    { value: "us-west-1", label: "US West (N. California)" },
                                    { value: 'ap-south-1', label: 'Asia Pacific (Mumbai)' },
                                    { value: 'ap-northeast-3', label: 'Asia Pacific (Osaka-Local)' },

                                ]}
                                {...form.getInputProps("s3Region")}
                            />
                        </>

                    )
                }
                {
                    form.values.type === "gcs" && (
                        <TextInput
                            withAsterisk
                            label="GCS API Key"
                            placeholder="Enter GCS API Key"
                            {...form.getInputProps("gcsApiKey")}
                        />
                    )
                }
                {
                    form.values.type === "azure" && (
                        <>
                            <TextInput
                                withAsterisk
                                label="Azure Account Name"
                                placeholder="Enter Azure Account Name"
                                {...form.getInputProps("azureAccountName")}
                            />
                            <TextInput
                                withAsterisk
                                label="Azure Account Key"
                                placeholder="Enter Azure Account Key"
                                {...form.getInputProps("azureAccountKey")}
                            />
                        </>
                    )
                }
            </Stack>

            <Group justify="flex-end" mt="md" gap="sm">
                {showCancelButton && <Button onClick={onCancel}>Cancel</Button>}
                {showSubmitButton && <Button type="submit">Submit</Button>}
            </Group>
        </form >
    );
};

export default NewConnectorForm;
