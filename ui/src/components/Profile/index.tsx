import { useEffect, useState } from 'react';
import { Container, Card, Text, Group, Button, TextInput, Grid, FileButton, Avatar, Stack } from '@mantine/core';
import { useForm } from '@mantine/form';

function ProfilePage() {
    const [isEditing, setIsEditing] = useState(false);
    const form = useForm({
        initialValues: {
            firstName: "John",
            lastName: "Doe",
            image: null,
            email: "",
            organization: "CodeMonk"
        },
        validate: {
            firstName: (value) => (value ? null : 'First name is required'),
            lastName: (value) => (value ? null : 'Last name is required'),
        },
    });

    useEffect(() => {
        const usermetadata = sessionStorage.getItem('usermetadata');
        if (usermetadata) {
            form.setValues(JSON.parse(usermetadata));
        }
    }, [])

    const handleEditToggle = () => {
        setIsEditing((prev) => !prev);
    };

    const handleSave = (values: any) => {
        console.log('Saved values:', values);
        setIsEditing(false);
    };

    return form.values.email ? (
        <Container size="md" mt="xl">
            <Card shadow="sm" padding="lg" radius="md" withBorder>
                <Group>
                    <Group p="center">
                        <Avatar size={100} src={form.values.image} alt="Profile Photo" />
                    </Group>
                    <Stack>
                        <Group p="center">
                            <FileButton onChange={() => { }} accept="image/png,image/jpeg">
                                {(props) => <Button {...props}>Upload New Photo</Button>}
                            </FileButton>
                            <Button variant="outline" color="red">
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
                                value={form.values.firstName}
                                disabled={!isEditing}
                            />
                        </Grid.Col>
                        <Grid.Col span={6}>
                            <TextInput
                                label="Last Name"
                                placeholder="Enter your last name"
                                {...form.getInputProps('lastName')}
                                value={form.values.lastName}
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
                                <Button onClick={handleEditToggle} variant="outline" color="gray">
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