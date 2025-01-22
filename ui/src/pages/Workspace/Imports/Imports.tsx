import { Box, Text, Title } from "@mantine/core"
import ImportsList from "./List"

const ImportsPage = () => {
    return (
        <Box mb={50}>
            <Title order={3} opacity={0.7}>Imports</Title>
            <Text c="dimmed" size="xs" mt={2} mb="lg">
                Imports may take a few seconds to appear here after upload. Please click the refresh button if you don't see any imports.
            </Text>
            <ImportsList />
        </Box>
    )
}

export default ImportsPage