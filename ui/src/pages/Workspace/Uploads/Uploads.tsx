import { Box, Text, Title } from "@mantine/core"
import UploadList from "./List"

const UploadsPage = () => {
    return (
        <Box mb={50}>
            <Title order={3} opacity={0.7}>Uploads</Title>
            <Text c="dimmed" size="xs" mt={2} mb="lg">
                Uploaded files may take a few seconds to appear here after upload. Please click the refresh button if you don't see any new files.
            </Text>
            <UploadList />
        </Box>
    )
}

export default UploadsPage