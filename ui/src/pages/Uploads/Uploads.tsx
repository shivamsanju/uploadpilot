import { Box, Text, Title, Group } from "@mantine/core"
import UploadList from "./List"

const UploadsPage = () => {
    return (
        <Box mb={50}>
            <Group align="center" justify="space-between">
                <Title order={3} opacity={0.7} >Uploads</Title>
                <Text c="dimmed" opacity={0.7} mr="lg">
                    Uploaded files may take a few seconds to appear here after upload. Please click the refresh button if you don't see any new files.
                </Text>
            </Group>
            <UploadList />
        </Box>
    )
}

export default UploadsPage