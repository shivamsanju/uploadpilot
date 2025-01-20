import { useLocation, useNavigate } from 'react-router-dom';
import { Alert, Button, Stack } from '@mantine/core';
import { IconMoodSad2 } from '@tabler/icons-react';

const ErrorQueryDisplay = () => {
    const location = useLocation();
    const navigate = useNavigate();

    const queryParams = new URLSearchParams(location.search);
    const error = queryParams.get('error');

    const handleBackToLogin = () => {
        navigate('/auth', { replace: true });
    };

    return (
        <Stack align="center" justify="center" mt="xl">
            <Alert variant="light" color="red" title="Error" icon={<IconMoodSad2 size={18} />}>
                {error || "An error occurred"}
            </Alert>

            <Button
                size="md"
                mt={40}
                variant="outline"
                onClick={handleBackToLogin}
            >
                Back to Login
            </Button>
        </Stack>
    );
};

export default ErrorQueryDisplay;
