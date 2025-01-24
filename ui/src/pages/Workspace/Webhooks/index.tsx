import { ActionIcon, Box, Group, Text, Title } from "@mantine/core"
import WebhooksList from "../Webhooks/List"
import { IconCodeVariablePlus } from "@tabler/icons-react"
import { useState } from "react"
import { useParams } from "react-router-dom"
import { AppLoader } from "../../../components/Loader/AppLoader"

const WorkspaceWebhooksPage = () => {
    const { workspaceId } = useParams();
    const [webhookAddModalOpened, setWorkspaceWebhooksOpened] = useState(false);

    if (!workspaceId) {
        return <AppLoader h="70vh" />
    }

    return (
        <Box mb={50}>
            <Group justify="space-between" m={0} p={0}>
                <Title order={3} opacity={0.7}>Webhooks</Title>
                <ActionIcon size="lg" variant="outline" onClick={() => setWorkspaceWebhooksOpened(true)} mr="xl">
                    <IconCodeVariablePlus size={20} />
                </ActionIcon>
            </Group>
            <Text c="dimmed" size="xs" mb="lg">Webhooks let you recieve events from your uploader</Text>
            <WebhooksList opened={webhookAddModalOpened} setOpened={setWorkspaceWebhooksOpened} />
        </Box>
    )
}

export default WorkspaceWebhooksPage