import { Button, Image, Paper, Stack, Text } from '@mantine/core';
import classes from './Auth.module.css';
import { getApiDomain } from '../../utils/config';
import TokenCatcher from './TokenCatcher';
import { useEffect, useState } from 'react';
import { AppLoader } from '../../components/Loader/AppLoader';
import { useNavigate } from 'react-router-dom';
import { Logo2 } from '../../components/Logo/Logo2';
import GoogleIcon from "../../assets/icons/google.svg";
import GithubIcon from "../../assets/icons/github.svg";

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
                    <Logo2 enableOnClick={false} />
                    <Text size="md" fw={500}>
                        Welcome to UploadPilot, Login with
                    </Text>
                </Stack>
                <Stack mt="xl">
                    <Button
                        variant='default'
                        leftSection={<Image src={GoogleIcon} width={20} height={20} />}
                        onClick={() => handleLogin('google')}
                        size="sm"
                    >
                        Google
                    </Button>
                    <Text size="xs" ta="center" c="dimmed">OR</Text>
                    <Button
                        variant='default'
                        leftSection={<Image src={GithubIcon} width={20} height={20} />}
                        onClick={() => handleLogin('github')}
                        size="sm"
                    >
                        Github
                    </Button>
                </Stack>
            </Paper>
        </div>
    );
}

export default AuthPage