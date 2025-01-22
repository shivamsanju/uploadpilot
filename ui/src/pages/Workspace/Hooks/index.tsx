import { Box, Text, Title } from "@mantine/core"
import { useParams } from "react-router-dom"
import { AppLoader } from "../../../components/Loader/AppLoader"
import { HooksMarketPlace } from "./Marketplace";

const WorkspaceWebhooksPage = () => {
    const { workspaceId } = useParams();

    if (!workspaceId) {
        return <AppLoader h="70vh" />
    }

    return (
        <Box pb="xl">
            <Title order={3} opacity={0.7}>Hooks</Title>
            <Text c="dimmed" mb="md" size="xs" mt={2}>
                Hooks let you add more fine tuned transformations and validations to your data.
            </Text>
            <HooksMarketPlace />
        </Box>
    )
}

export default WorkspaceWebhooksPage