import { useEffect, useState } from 'react';
import { Container, Card, Text, Group, Button, TextInput, Grid, FileButton, Avatar, Stack } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useGetSession } from '../../apis/user';
import { AppLoader } from '../../components/Loader/AppLoader';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';

const ProfilePage = () => {
    const [isEditing, setIsEditing] = useState(false);
    const { isPending, error, session } = useGetSession();
    const form = useForm({
        initialValues: {
            name: "John",
            avatarUrl: null,
            email: "",
            organization: ""
        },
        validate: {
            name: (value) => (value ? null : 'First name is required'),
            email: (value) => (value ? null : 'Email is required'),
        },
    });


    useEffect(() => {
        form.setValues(session)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [session])

    const handleEditToggle = () => {
        setIsEditing((prev) => !prev);
    };

    const handleSave = (values: any) => {
        console.log('Saved values:', values);
        setIsEditing(false);
    };

    if (error) {
        return <ErrorCard title={error.name} message={error.message} h="70vh" />
    }

    return isPending ? <AppLoader h="70vh" /> : (session.email || session.name) ? (
        <Container size="md" mt="xl">
            <Card shadow="sm" padding="lg" radius="md" withBorder>
                <Group>
                    <Group p="center">
                        <Avatar size={100} src={form.values.avatarUrl} alt="Profile Photo" />
                    </Group>
                    <Stack>
                        <Group p="center">
                            <FileButton onChange={() => { }} accept="image/png,image/jpeg">
                                {(props) => <Button {...props}>Upload New Photo</Button>}
                            </FileButton>
                            <Button variant="light">
                                Reset
                            </Button>
                        </Group>
                        <Text size="sm" c="dimmed">
                            Allowed JPG, GIF or PNG. Max size of 800K
                        </Text>
                    </Stack>
                </Group>

                <form onSubmit={form.onSubmit(handleSave)}>
                    <Grid mt="lg">
                        <Grid.Col span={6}>
                            <TextInput
                                label="First Name"
                                placeholder="Enter your first name"
                                {...form.getInputProps('firstName')}
                                value={form.values.name}
                                disabled={!isEditing}
                            />
                        </Grid.Col>
                        <Grid.Col span={6}>
                            <TextInput
                                label="E-mail"
                                placeholder="Enter your email"
                                {...form.getInputProps('email')}
                                value={form.values.email}
                                disabled={true}
                            />
                        </Grid.Col>
                        <Grid.Col span={6}>
                            <TextInput
                                label="Organization"
                                placeholder="Enter your organization"
                                {...form.getInputProps('organization')}
                                value={form.values.organization}
                                disabled={true}
                            />
                        </Grid.Col>
                    </Grid>
                    <Group p="apart" mt="xl">
                        {isEditing ? (
                            <>
                                <Button type="submit">Save Changes</Button>
                                <Button onClick={handleEditToggle} variant="light">
                                    Reset
                                </Button>
                            </>
                        ) : (
                            <Button onClick={handleEditToggle}>Edit Profile</Button>
                        )}
                    </Group>
                </form>
            </Card>
        </Container>
    ) : <></>;
}

export default ProfilePage;