import { ActionIcon, Box, Button, Divider, Group, LoadingOverlay, Stack, Text, TextInput } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useCreateWorkspaceMutation, useGetWorkspaces } from '../../apis/workspace';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { AppLoader } from '../../components/Loader/AppLoader';
import { useNavigate } from 'react-router-dom';
import classes from "./Workspace.module.css"
import { IconChevronRight } from '@tabler/icons-react';

const WorkspaceLandingPage = () => {
    const { isPending, error, workspaces } = useGetWorkspaces();
    const { mutateAsync, isPending: isCreating } = useCreateWorkspaceMutation();
    const navigate = useNavigate();

    const form = useForm({
        initialValues: {
            name: '',
            tags: []
        },
        validate: {
            name: (value) => {
                if (!value.trim()) {
                    return 'Workspace name is required';
                }
                if (value.trim().length > 20 || value.trim().length < 2) {
                    return 'Workspace name must be between 2 and 20 characters';
                }
                return null;
            },
        }
    })

    const handleCreateWorkspace = async (values: any) => {
        try {
            const id = await mutateAsync(values.name)
            navigate(`/workspaces/${id}`);
        } catch (error) {
            console.error(error);
        }
    }


    return (
        <Group justify='center' mb="50">
            {
                error ?
                    <ErrorCard title="Error" message={error.message} h="55vh" /> :
                    isPending ? <AppLoader h="55vh" />
                        : <Box>
                            <LoadingOverlay visible={isCreating} overlayProps={{ radius: "sm", blur: 1 }} />
                            <form onSubmit={form.onSubmit((values) => handleCreateWorkspace(values))}>
                                <Stack className={classes.wsForm} >
                                    <Text size="xl" fw={700} opacity={0.7} ta='left'>
                                        {workspaces && workspaces.length > 0 ? 'Create a new workspace' : 'Create a new workspace to get started'}
                                    </Text>
                                    <TextInput
                                        w="100%"
                                        placeholder="Enter a workspace name"
                                        {...form.getInputProps('name')}
                                    />
                                    <Button type="submit" w="100%">Create</Button>
                                </Stack>
                            </form>
                            {workspaces && workspaces.length > 0 && (
                                <>
                                    <Divider mt={50} mb={50} />
                                    <Stack >
                                        <Text size="xl" fw={700} opacity={0.7}>Choose an existing workspace</Text>
                                        {workspaces && workspaces.length > 0 && workspaces.map((workspace: any) => (
                                            <Group justify='space-between' key={workspace.id} className={classes.wsItem} pt="lg">
                                                <Text size="sm" fw="bold" opacity={0.7} >{workspace.name}</Text>
                                                <ActionIcon
                                                    mt={8}
                                                    radius="50%"
                                                    variant="default"
                                                    size="lg"
                                                    onClick={() => navigate(`/workspaces/${workspace.id}`)}
                                                >
                                                    <IconChevronRight color="gray" />
                                                </ActionIcon>
                                            </Group>
                                        ))}
                                    </Stack>
                                </>
                            )}
                        </Box>
            }
        </Group>
    );
}

export default WorkspaceLandingPage