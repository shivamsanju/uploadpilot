import { ActionIcon, Box, Group, Text, Title } from "@mantine/core"
import UserList from "./List"
import { IconUserPlus } from "@tabler/icons-react"
import { useParams } from "react-router-dom"
import { AppLoader } from "../../../components/Loader/AppLoader"
import { useState } from "react"

const WorkspaceUsersPage = () => {
    const [workspaceUsersOpened, setWorkspaceUsersOpened] = useState(false);
    const { workspaceId } = useParams();


    if (!workspaceId) {
        return <AppLoader h="70vh" />
    }

    return (
        <Box pb="xl">
            <Group justify="space-between">
                <Title order={3} opacity={0.7}>Users</Title>
                <ActionIcon size="lg" variant="outline" onClick={() => setWorkspaceUsersOpened(true)} mr="xl">
                    <IconUserPlus size={20} />
                </ActionIcon>
            </Group>
            <Text c="dimmed" size="xs" mb="lg">Manage users and roles in this workspace</Text>

            <UserList opened={workspaceUsersOpened} setOpened={setWorkspaceUsersOpened} />
        </Box>
    )
}

export default WorkspaceUsersPage