import { Button, Group, Paper, Stack, Text, Title } from "@mantine/core";

type ErrorCardProps = {
    title: string
    message: string
    h?: string
}
export const ErrorCard: React.FC<ErrorCardProps> = ({ title, message, h }) => {
    const refreshPage = () => window.location.reload();
    return (
        <Paper shadow="xs" p="sm" radius="xs" withBorder>
            <Stack p="xl" align='center' justify='center' h={h} >

                <Title order={3} mb="lg">{title}</Title>
                <Text size="md" ta="center" mb="md" >
                    {message}
                </Text>
                <Group justify="center">
                    <Button size="md" variant="primary" onClick={refreshPage}>
                        Refresh the page
                    </Button>
                </Group>
            </Stack>
        </Paper>
    )
}