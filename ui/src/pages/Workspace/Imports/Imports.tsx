import { Box, Title } from "@mantine/core"
import ImportsList from "./List"

const ImportsPage = () => {
    return (
        <Box mb={50}>
            <Title order={3} mb="lg" opacity={0.7}>Imports</Title>
            <ImportsList />
        </Box>
    )
}

export default ImportsPage