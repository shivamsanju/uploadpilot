import { Box, Text, Title } from "@mantine/core"
import { useParams } from "react-router-dom"
import { AppLoader } from "../../../components/Loader/AppLoader"
import { UtilitiesGrid } from "./Grid";

const UtilitiesPage = () => {
    const { workspaceId } = useParams();

    if (!workspaceId) {
        return <AppLoader h="70vh" />
    }

    return (
        <Box mb={50} mr="sm">
            <Title order={3} opacity={0.7}>Utilities</Title>
            <Text c="dimmed" mb="md" mt={2}>
                Add more features and transformations to your files
            </Text>
            <UtilitiesGrid />
        </Box>
    )
}

export default UtilitiesPage