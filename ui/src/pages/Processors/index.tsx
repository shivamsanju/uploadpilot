import { ActionIcon, Box, Group, Title } from "@mantine/core"
import ProcessorsList from "./List"
import { IconPlus } from "@tabler/icons-react"
import { useState } from "react"
import { useParams } from "react-router-dom"
import { AppLoader } from "../../components/Loader/AppLoader"

const WorkspaceProcessorsPage = () => {
    const { workspaceId } = useParams();
    const [webhookAddModalOpened, setWorkspaceProcessorsOpened] = useState(false);

    if (!workspaceId) {
        return <AppLoader h="70vh" />
    }

    return (
        <Box mb={50}>
            <Group justify="space-between" mb="xl">
                <Title order={3} opacity={0.7}>Processors</Title>
                <ActionIcon size="lg" variant="outline" onClick={() => setWorkspaceProcessorsOpened(true)} mr="lg">
                    <IconPlus size={20} />
                </ActionIcon>
            </Group>
            {/* <Text c="dimmed" size="xs" mb="lg">Processors let you process uploaded files</Text> */}
            <ProcessorsList opened={webhookAddModalOpened} setOpened={setWorkspaceProcessorsOpened} />
        </Box>
    )
}

export default WorkspaceProcessorsPage