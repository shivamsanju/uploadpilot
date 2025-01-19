import { useState } from "react";
import { Button, Container, Group, Paper, Stack, Stepper, Title } from "@mantine/core";
import { useForm } from "@mantine/form";
import MetadataPage from "./Form/MetadataForm";
import ConfigurationForm from "./Form/ConfigurationForm";
import DataStorePage from "./Form/DataStoreForm";
import { CreateUploaderForm } from "../../../types/uploader";
import { useNavigate } from "react-router-dom";
import { useCreateUploaderMutation } from "../../../apis/uploader";
import { notifications } from "@mantine/notifications";

const CreateNewUploaderPage = () => {
    const [page, setPage] = useState(1);

    const navigate = useNavigate();
    const { mutateAsync } = useCreateUploaderMutation();

    const handleSubmit = async () => {
        if (form.validate().hasErrors) {
            console.log(form.errors);
            return;
        };

        mutateAsync({
            name: form.values.name,
            description: form.values.description,
            tags: form.values.tags,
            config: {
                minFileSize: form.values.minFileSize,
                maxFileSize: form.values.maxFileSize,
                minNumberOfFiles: form.values.minNumberOfFiles,
                maxNumberOfFiles: form.values.maxNumberOfFiles,
                maxTotalFileSize: form.values.maxTotalFileSize,
                allowedFileTypes: form.values.allowedFileTypes,
                allowedSources: form.values.allowedSources,
                requiredMetadataFields: form.values.requiredMetadataFields,
                enableImageEditing: form.values.enableImageEditing,
                useCompression: form.values.useCompression,
                useFaultTolerantMode: form.values.useFaultTolerantMode,
            },
            dataStore: {
                name: form.values.dataStoreName,
                connectorId: form.values.connectorId,
                connectorName: form.values.connectorName,
                connectorType: form.values.connectorType,
                bucket: form.values.bucket,
            }
        }).then(() => {
            navigate('/uploaders');
        })
    }

    const form = useForm<CreateUploaderForm>({
        initialValues: {
            name: "",
            description: "",
            tags: [] as string[],
            dataStoreName: "",
            connectorId: "",
            connectorName: "",
            connectorType: "",
            bucket: "",
            allowedFileTypes: [] as string[],
            allowedSources: [] as string[],
            requiredMetadataFields: [] as string[],
            enableImageEditing: false,
            useCompression: false,
            useFaultTolerantMode: false
        },
        validate: {
            name: (value) => {
                if (!value.trim()) {
                    return "Name is required";
                }
                if (value.length > 100 || value.length < 2) {
                    return "Name must be between 2 and 100 characters";
                }

                return null;
            },
            tags: (value) => value.length > 5 ? "Maximum of 5 tags allowed" : null,
            description: (value) => {
                if (!value.trim()) {
                    return "Description is required";
                }
                if (value.length > 1000 || value.length < 10) {
                    return "Description must be between 10 and 1000 characters";
                }
                return null;
            },
            allowedFileTypes: (value) => page === 2 && value.length === 0 ? "Please select at least one file type" : null,
            allowedSources: (value) => page === 2 && value.length === 0 ? "Please select at least one source" : null,
            bucket: (value) => {
                if (page === 3 && !value) {
                    return "Please enter a bucket name";
                }
                const bucketRegex = /^[a-z][a-z0-9]*$/;
                if (page === 3 && !bucketRegex.test(value)) {
                    return "Bucket name must start with a lowercase letter and contain only lowercase letters and numbers";
                }
                return null;
            },
            dataStoreName: (value) => page === 3 && !value ? "Please enter a datastore name" : null,
            connectorId: (value) => page === 3 && !value ? "Please select a connector" : null,
        },
    });

    const handleNextPage = () => {
        if (!form.validate().hasErrors) {
            setPage((prev) => prev + 1)
        } else {
            notifications.show({
                title: "Error",
                message: "Form has errors. Please fix errors and try again",
                color: "red",
            })
        }
    };

    const handlePrevPage = () => {
        if (page === 1) {
            navigate('/uploaders', { replace: true });
        }
        setPage((prev) => prev - 1);
    }


    return (
        <Container p="sm" >
            <form onSubmit={form.onSubmit(handleSubmit)}>
                <Group justify='center' mb="sm">
                    <Title order={3} mb="sm" opacity={0.8}>Create new uploader</Title>
                </Group>
                <Stepper active={page - 1} mb="sm" mt={30} size='sm'>
                    <Stepper.Step label="Step 1" description="Fill Metadata">
                        Uploader details
                    </Stepper.Step>
                    <Stepper.Step label="Step 2" description="Configuration">
                        Configuration
                    </Stepper.Step>
                    <Stepper.Step label="Final step" description="Create datastore">
                        Create Datastore
                    </Stepper.Step>
                    <Stepper.Completed>
                        Completed, click back button to get to previous step
                    </Stepper.Completed>
                </Stepper>
                <Paper p="lg" radius="md" withBorder>
                    <Stack h="57vh" style={{ overflowY: 'auto' }} pr="lg">
                        {page === 1 && <MetadataPage form={form} />}
                        {page === 2 && <ConfigurationForm form={form} />}
                        {page === 3 && <DataStorePage form={form} />}
                    </Stack>
                    <Group justify="space-between" mt="xl">
                        <Button variant="light" onClick={handlePrevPage}>{page > 1 ? "Prev" : "Cancel"}</Button>
                        {page < 3 ? (
                            <Button onClick={handleNextPage}>Next</Button>
                        ) : (
                            <Button onClick={handleSubmit}>Create</Button>
                        )}
                    </Group>
                </Paper>
            </form>
        </Container>
    );
};

export default CreateNewUploaderPage;
