import { Box, Title } from "@mantine/core"
import UserList from "./List"

const WorkspaceUsersPage = () => {
    return (
        <Box>
            <Title order={3} mb="lg" opacity={0.7}>Users</Title>
            <UserList />
        </Box>
    )
}

export default WorkspaceUsersPage