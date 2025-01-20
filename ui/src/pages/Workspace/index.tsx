import { Box, Button, Divider, Group, Stack, Text, TextInput } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useCreateWorkspaceMutation, useGetWorkspaces } from '../../apis/workspace';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { AppLoader } from '../../components/Loader/AppLoader';
import { useNavigate } from 'react-router-dom';
import classes from "./Workspace.module.css"

const WorkspaceLandingPage = () => {
    const { isPending, error, workspaces } = useGetWorkspaces();
    const { mutateAsync } = useCreateWorkspaceMutation();
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

    const handleCreateWorkspace = (values: any) => {
        mutateAsync(values.name).then((id) => {
            console.log(id)
            navigate(`/workspaces/${id}`);
        })
    }


    return (
        <Group justify='center' pb="xl">
            {
                error ?
                    <ErrorCard title="Error" message={error.message} h="55vh" /> :
                    isPending ? <AppLoader h="55vh" />
                        : <Box>
                            <form onSubmit={form.onSubmit((values) => handleCreateWorkspace(values))}>
                                <Stack h="20vh" w="40vw" miw="300" mt="xl">
                                    <Text size="xl" fw={700} opacity={0.7}>{workspaces && workspaces.length > 0 ? 'Create a new workspace' : 'Create a new workspace to get started'}</Text>
                                    <TextInput
                                        placeholder="Enter a workspace name"
                                        {...form.getInputProps('name')}
                                    />
                                    {/* <TagsInput
                                        placeholder="Enter tags(comma-separated)"
                                        {...form.getInputProps('tags')}
                                    /> */}
                                    <Button type="submit">Create</Button>
                                </Stack>
                            </form>
                            {workspaces && workspaces.length > 0 && (
                                <>
                                    <Divider />
                                    <Stack mt="md" >
                                        <Text size="xl" fw={700} opacity={0.7}>Choose an existing workspace</Text>
                                        {workspaces && workspaces.length > 0 && workspaces.map((workspace: any) => (
                                            <Group justify='space-between' key={workspace.id} className={classes.wsItem} pt="lg">
                                                <Text fw="bold" opacity={0.7} >{workspace.name}</Text>
                                                <Button key={workspace.id} variant='outline' onClick={() => navigate(`/workspaces/${workspace.id}`)}>Open</Button>
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