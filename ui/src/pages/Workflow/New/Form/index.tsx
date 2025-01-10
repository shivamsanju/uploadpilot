import React, { useState } from "react";
import { Button, Group, Paper, Stack, Stepper, Title } from "@mantine/core";
import { useForm } from "@mantine/form";
import WorkflowMetadataPage from "./WorkflowMetadataForm";
import ImportPolicyPage from "./ImportPolicyPage";
import DataStorePage from "./DataStoreForm";
import { CreateWorkflowForm } from "../../../../types/workflow";

const CreateWorkflowFormPage = ({ onSubmit }: { onSubmit: (workflow: any) => void }) => {
    const [page, setPage] = useState(1);

    const form = useForm<CreateWorkflowForm>({
        initialValues: {
            name: "",
            description: "",
            tags: [] as string[],
            importPolicyId: "",
            connectorId: "",
            dataStoreId: "",
            dataStoreName: "",
            bucket: "",
        },
        validate: {
            name: (value) => (!value.trim() ? "Name is required" : null),
            description: (value) => (!value.trim() ? "Description is required" : null),
            importPolicyId: (value) =>
                page === 2 && !value ? "Please select an import policy" : null,
            connectorId: (value) => page === 3 && !value ? "Please select a connector" : null,
            bucket: (value) => page === 3 && !value ? "Please enter a bucket name" : null,
            dataStoreName: (value) => page === 3 && !value ? "Please enter a datastore name" : null
        },
    });

    const handleNextPage = () => {
        if (!form.validate().hasErrors) setPage((prev) => prev + 1);
    };

    const handlePrevPage = () => setPage((prev) => prev - 1);


    const handleSubmit = () => {
        if (form.validate().hasErrors) {
            console.log(form.errors);
            return;
        };
        onSubmit(form.values);
    };

    return (
        <form onSubmit={form.onSubmit(handleSubmit)}>
            <Group justify='center' mb="sm">
                <Title order={3} mb="sm" opacity={0.8}>Create new workflow</Title>
            </Group>
            <Stepper active={page - 1} mb="sm" mt={30} size='sm'>
                <Stepper.Step label="Step 1" description="Fill Metadata">
                    Add Workflow metadata
                </Stepper.Step>
                <Stepper.Step label="Step 2" description="Add import policy">
                    Add import policy
                </Stepper.Step>
                <Stepper.Step label="Final step" description="Create datastore">
                    Create Datastore
                </Stepper.Step>
                <Stepper.Completed>
                    Completed, click back button to get to previous step
                </Stepper.Completed>
            </Stepper>
            <Paper shadow='xs' p="lg" radius="md" withBorder>
                <Stack h="57vh" style={{ overflowY: 'auto' }}>
                    {page === 1 && <WorkflowMetadataPage form={form} />}
                    {page === 2 && <ImportPolicyPage form={form} />}
                    {page === 3 && <DataStorePage form={form} />}
                </Stack>
                <Group justify="space-between" mt="xl">
                    {page > 1 ? <Button onClick={handlePrevPage}>Prev</Button> : <div />}
                    {page < 3 ? (
                        <Button onClick={handleNextPage}>Next</Button>
                    ) : (
                        <Button onClick={handleSubmit}>Create</Button>
                    )}
                </Group>
            </Paper>
        </form>
    );
};

export default CreateWorkflowFormPage;
