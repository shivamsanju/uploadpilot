import { Paper, Stack, Text } from '@mantine/core';
import { Logo } from '../../components/Logo/Logo';
import { GoogleButton } from './GoogleButton';
import { GithubButton } from './GithubButton';
import classes from './Auth.module.css';
import { getApiDomain } from '../../utils/config';
import TokenCatcher from './TokenCatcher';
import { useEffect, useState } from 'react';
import { AppLoader } from '../../components/Loader/AppLoader';
import { useNavigate } from 'react-router-dom';

const AuthPage = () => {
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    const handleLogin = (provider: string) => {
        window.location.href = getApiDomain() + `/auth/${provider}/authorize`
    }

    useEffect(() => {
        const token = localStorage.getItem('uploadpilottoken');
        if (token) {
            navigate('/', { replace: true });
        }
        setLoading(false);
    }, [navigate]);

    console.log(loading);

    if (loading) {
        return <AppLoader h="100vh" />;
    }

    return (
        <div className={classes.wrapper}>
            <TokenCatcher />
            <Paper className={classes.form} radius={0} p={30}>
                <Stack align="center" gap="xs">
                    <Logo enableOnClick={false} />
                    <Text size="sm" fw={500}>
                        Welcome to UploadPilot, Login with
                    </Text>
                </Stack>
                <Stack gap="md" mt="xl">
                    <GoogleButton radius="xl" onClick={() => handleLogin('google')}>Google</GoogleButton>
                    <GithubButton radius="xl" onClick={() => handleLogin('github')}>Github</GithubButton>
                </Stack>
            </Paper>
        </div>
    );
}

export default AuthPage