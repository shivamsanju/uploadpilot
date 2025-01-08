import {
    TextInput,
    Textarea,
    Select,
    Button,
    Group,
    Stack,
    Card,
    TagsInput,
    MultiSelect,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useState } from "react";

const CreateWorkflowForm = ({ onSubmit }: { onSubmit: (codebase: any) => void }) => {
    const [page, setPage] = useState(1);

    const form = useForm({
        initialValues: {
            name: "",
            description: "",
            tags: [] as string[],
            sources: [] as string[],
            destination: "",
            s3BucketName: "",
            s3Region: "",
            s3AccessKey: "",
            s3SecretKey: "",
        },

        validate: {
            name: (value) => (!value || value.trim().length === 0 ? "Name is required" : null),
            description: (value) => (!value || value.trim().length === 0 ? "Description is required" : null),
            sources: (value) => ((page === 2 && value.length === 0) ? "Please select at least one source" : null),
            destination: (value) => ((page === 2 && !value) ? "Please select a destination" : null),
            s3BucketName: (value, values) => ((page === 3 && values.destination === "s3" && !value) ? "Please enter a bucket name" : null),
            s3Region: (value, values) => ((page === 3 && values.destination === "s3" && !value) ? "Please enter a region" : null),
            s3AccessKey: (value, values) => ((page === 3 && values.destination === "s3" && !value) ? "Please enter an access key" : null),
            s3SecretKey: (value, values) => ((page === 3 && values.destination === "s3" && !value) ? "Please enter a secret key" : null),
        },
    });

    const handleNextPage = () => {
        if (form.validate().hasErrors) {
            return;
        }
        setPage(page + 1);
    };

    const handlePrevPage = () => {
        setPage(page - 1);
    }


    return (
        <form
            onSubmit={form.onSubmit((values) => {
                const newWorkflow = {
                    name: values.name,
                    description: values.description,
                    tags: values.tags,
                    updated: new Date().toISOString(),
                    sources: values.sources,
                    destination: {
                        type: values.destination,
                        config: {}
                    }
                };
                if (values.destination === "s3") {
                    newWorkflow.destination.config = {
                        [values.destination]: {
                            bucketName: values.s3BucketName,
                            region: values.s3Region,
                            accessKey: values.s3AccessKey,
                            secretKey: values.s3SecretKey,
                        }
                    };
                }
                onSubmit(newWorkflow);
            })}
        >
            {page === 1 &&
                <Card withBorder radius="md" title="Workflow Metadata">
                    <Stack gap="md" h="40vh">
                        <TextInput
                            label="Name"
                            placeholder="Enter workflow name"
                            {...form.getInputProps("name")}
                        />

                        <Textarea
                            label="Description"
                            placeholder="Enter a brief description"
                            autosize
                            minRows={6}
                            {...form.getInputProps("description")}
                        />

                        <TagsInput
                            label="Press Enter to submit a tag"
                            placeholder="Enter tags"
                            {...form.getInputProps("tags")}
                        />
                    </Stack>

                    <Group justify="flex-end" mt="md">
                        <Button onClick={handleNextPage}>Next</Button>
                    </Group>
                </Card>
            }

            {page === 2 &&
                <Card withBorder radius="md" title="Source and Destination">
                    <Stack gap="md" h="40vh">
                        <MultiSelect
                            size="xs"
                            label="Allowed input sources"
                            placeholder="Select all allowed input sources"
                            data={[
                                { value: "file-upload", label: "File Upload" },
                                { value: "google-drive", label: "Google Drive" },
                                { value: "dropbox", label: "Dropbox" },
                                { value: "onedrive", label: "OneDrive" },
                                { value: "box", label: "Box" },
                                { value: "sharepoint", label: "Sharepoint" },
                                { value: "s3", label: "S3" },

                            ]}
                            onChange={(values) => form.setFieldValue("sources", values)}
                        />
                        <Select
                            label="Destination"
                            placeholder="Select Destination"
                            data={[
                                { value: "s3", label: "AWS S3" },
                                { value: "gcs", label: "Google Cloud Storage" },
                                { value: "azure", label: "Azure Blob Storage" },
                                { value: "local", label: "Local" },
                            ]}
                            {...form.getInputProps("destination")}
                        />
                    </Stack>

                    <Group justify="flex-end" mt="md" gap="md">
                        <Button onClick={handlePrevPage}>Prev</Button>
                        <Button onClick={handleNextPage}>Next</Button>
                    </Group>
                </Card>
            }

            {page === 3 &&
                <Card withBorder radius="md" title="Destination Configuration">
                    {form.values.destination === "s3" && (
                        <Stack gap="md" h="40vh">
                            <TextInput
                                label="S3 Access Key"
                                placeholder="Enter S3 Access Key"
                                {...form.getInputProps("s3AccessKey")}
                            />
                            <TextInput
                                label="S3 Secret Key"
                                placeholder="Enter S3 Access Key"
                                {...form.getInputProps("s3SecretKey")}
                            />
                            <Select
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
                            <TextInput
                                label="S3 Bucket Name"
                                placeholder="Enter S3 Bucket Name"
                                {...form.getInputProps("s3BucketName")}
                            />
                        </Stack>
                    )}
                    < Group justify="flex-end" mt="md" gap="md">
                        <Button onClick={handlePrevPage}>Prev</Button>
                        <Button type="submit">Submit</Button>
                    </Group>
                </Card>
            }
        </form >
    );
};

export default CreateWorkflowForm;
