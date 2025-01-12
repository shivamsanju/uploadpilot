import { Button, Card, Group, Select, Stack, Text } from "@mantine/core"
import { notifications } from "@mantine/notifications"

type ImportsProps = {
    uploaderDetails: any
}
const Imports: React.FC<ImportsProps> = ({ uploaderDetails }) => {

    const handleEdit = () => {
        notifications.show({
            title: 'Coming Soon',
            message: 'This feature is coming soon',
            color: 'yellow',
        })
    }
    return (
        <Stack gap="md" >
            <Card shadow="xs" p="sm" radius="xs" withBorder>
                <Group justify="space-between" mb="sm">
                    <Text size="lg" fw="bold" mb="sm">DataStore</Text>
                    <Button size="xs" onClick={handleEdit}>Edit</Button>
                </Group>
                <Select
                    value={uploaderDetails.dataStoreDetails.name}
                    disabled
                    data={[uploaderDetails.dataStoreDetails.name]}
                />
            </Card>

        </Stack>
    )
}

export default Imports