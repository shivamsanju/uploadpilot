import {
    Anchor,
    Button,
    Checkbox,
    Divider,
    Group,
    Paper,
    PaperProps,
    PasswordInput,
    Stack,
    Text,
    TextInput,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useToggle } from '@mantine/hooks';
import axiosInstance from '../../utils/axios';
import { notifications } from '@mantine/notifications';
import { Logo } from '../../components/Logo/Logo';


type FormValues = {
    firstName?: string;
    lastName?: string;
    email: string;
    password: string;
    confirmPassword?: string;
    rootPassword?: string;
    terms?: boolean;
}
const AuthPage = (props: PaperProps) => {
    const [type, toggle] = useToggle(['login', 'register']);
    const form = useForm<FormValues>({
        initialValues: {
            firstName: '',
            lastName: '',
            email: '',
            password: '',
            confirmPassword: '',
            rootPassword: '',
            terms: true,
        },

        validate: {
            firstName: (val) => ((type === 'register' && (!val || val.length < 2)) ? 'First name should include at least 2 letters' : null),
            lastName: (val) => ((type === 'register' && (!val || val.length < 2)) ? 'Last name should include at least 2 letters' : null),
            email: (val) => (/^\S+@\S+$/.test(val) ? null : 'Invalid email'),
            password: (val) => (val.length <= 6 ? 'Password should include at least 6 characters' : null),
            confirmPassword: (val, values) => (type === 'register' && val !== values.password ? 'Passwords did not match' : null),
            rootPassword: (val) => (type === 'register' && !val ? 'Please enter the root password' : null),
        },
    });


    const redirectToHomePage = () => {
        window.location.href = '/';
    }

    const onSubmit = async (values: FormValues) => {
        let response;
        try {
            if (type === 'register') {
                response = await axiosInstance.post('/auth/signup', values)
            } else {
                response = await axiosInstance.post('/auth/login', values)
            }

            if (response.status === 200) {
                localStorage.setItem('token', response.data.token);
                localStorage.setItem('refreshToken', response.data.refreshToken);
                redirectToHomePage();
            } else {
                console.log(response.data);
                notifications.show({
                    title: 'Error',
                    message: response.data.message,
                    color: 'red',
                })
            }
        } catch (error: any) {
            notifications.show({
                title: 'Error',
                message: error.response.data.message,
                color: 'red',
            })
        }

    }



    return (
        <Group justify="center" align="center" style={{ minHeight: '100vh' }}>
            <Paper radius="md" p="xl" withBorder {...props} w="30vw" h={type === 'register' ? '82vh' : '45vh'}>
                <Stack align="center" gap="xs">
                    <Logo enableOnClick={false} />
                    <Text size="sm" fw={500}>
                        Welcome to UploadPilot, {type === 'register' ? 'Create  User' : 'Login'}
                    </Text>
                </Stack>

                <form onSubmit={form.onSubmit(onSubmit)}>
                    <Stack>
                        {
                            type === 'register' && (
                                <>
                                    <PasswordInput
                                        size='xs'
                                        required
                                        label="Root password"
                                        placeholder="Root passwor for creating user"
                                        value={form.values.rootPassword}
                                        onChange={(event) => form.setFieldValue('rootPassword', event.currentTarget.value)}
                                        error={form.errors.rootPassword && 'Passwords did not match'}
                                        radius="md"
                                    />
                                    <Divider label="New User Details" labelPosition="center" />
                                </>
                            )
                        }
                        {type === 'register' && (
                            <TextInput
                                withAsterisk
                                label="First Name"
                                placeholder="First name"
                                value={form.values.firstName}
                                onChange={(event) => form.setFieldValue('firstName', event.currentTarget.value)}
                                error={form.errors.firstName && 'First name should include at least 2 letters'}
                                radius="md"
                            />

                        )}

                        {type === 'register' && (
                            <TextInput
                                withAsterisk
                                label="Last Name"
                                placeholder="Last name"
                                value={form.values.lastName}
                                onChange={(event) => form.setFieldValue('lastName', event.currentTarget.value)}
                                error={form.errors.lastName && 'Last name should include at least 2 letters'}
                                radius="md"
                            />

                        )}

                        <TextInput
                            required
                            label="Email"
                            placeholder="hello@mantine.dev"
                            value={form.values.email}
                            onChange={(event) => form.setFieldValue('email', event.currentTarget.value)}
                            error={form.errors.email && 'Invalid email'}
                            radius="md"
                        />

                        <PasswordInput
                            size='xs'
                            required
                            label="Password"
                            placeholder="Your password"
                            value={form.values.password}
                            onChange={(event) => form.setFieldValue('password', event.currentTarget.value)}
                            error={form.errors.password && 'Password should include at least 6 characters'}
                            radius="md"
                        />

                        {
                            type === 'register' && (
                                <PasswordInput
                                    size='xs'
                                    required
                                    label="Confirm password"
                                    placeholder="Your password"
                                    value={form.values.confirmPassword}
                                    onChange={(event) => form.setFieldValue('confirmPassword', event.currentTarget.value)}
                                    error={form.errors.confirmPassword && 'Passwords did not match'}
                                    radius="md"
                                />
                            )
                        }


                        {type === 'register' && (
                            <Checkbox
                                label="I accept terms and conditions"
                                checked={form.values.terms}
                                onChange={(event) => form.setFieldValue('terms', event.currentTarget.checked)}
                            />
                        )}
                    </Stack>

                    <Group justify="space-between" mt="xl">
                        <Anchor size="xs" component="button" type="button" c="dimmed" onClick={() => toggle()}>
                            {type === 'register'
                                ? 'Already have an account? Login'
                                : "Create User (Needs root permission)"}
                        </Anchor>
                        <Button type="submit" radius="xl">
                            {type === 'register' ? 'Create' : 'Login'}
                        </Button>
                    </Group>
                </form>
            </Paper>
        </Group >
    );
}

export default AuthPage