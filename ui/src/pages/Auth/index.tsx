import { Button, Image, Box, Stack, Text } from '@mantine/core';
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


    if (loading) {
        return <AppLoader h="100vh" />;
    }

    return (
        <div className={classes.wrapper}>
            <TokenCatcher />
            <Box className={classes.form} p="70" bg="white" pt="lg">
                <Stack gap="xs" mb="60">
                    <Logo2 enableOnClick={false} />
                </Stack>
                <Stack mt="xl">
                    <Text size="xs" ta="center" c="dimmed">
                        By continuing, you agree to our Terms of Service and acknowledge you have read our Privacy Policy
                    </Text>
                    <Button
                        variant='outline'
                        leftSection={<Image src={GoogleIcon} width={20} height={20} />}
                        onClick={() => handleLogin('google')}
                        size="sm"
                    >
                        Google
                    </Button>
                    <Text ta="center" c="dimmed">or</Text>
                    <Button
                        variant='outline'
                        leftSection={<Image src={GithubIcon} width={25} height={25} />}
                        onClick={() => handleLogin('github')}
                        size="sm"
                    >
                        Github
                    </Button>
                </Stack>
            </Box>
        </div>
    );
}

export default AuthPage