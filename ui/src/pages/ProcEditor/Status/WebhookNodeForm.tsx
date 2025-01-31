import { useForm } from "@mantine/form";
import { Button, Group, Modal, Stack, Text, TextInput, UnstyledButton } from "@mantine/core";
import { IconClockBolt, IconHttpPost, IconLock, IconX } from "@tabler/icons-react";
import { useState } from "react";
import { CodeHighlight } from "@mantine/code-highlight";
import { NodeFormProps } from ".";

const getSampleCurl = (url: string) => {
    return `
curl -X POST ${url || "https://example.com/webhook"} \\
-H "Content-Type: application/json" \\
-H "Secret: <your secret>" \\
-d '{
  "file_url": "https://example.com/uploads/file123.jpg",
  "file_name": "file123.jpg",
  "file_size": 1048576,
  "upload_time": "2025-01-22T10:30:00Z"
}'
`
}



const WebhookNodeForm: React.FC<NodeFormProps> = ({ nodeData, saveNodeData, setOpenedNodeId }) => {
    const [openSampleReq, setOpenSampleReq] = useState(false);

    const form = useForm({
        initialValues: {
            url: nodeData?.url || "",
            secret: nodeData?.secret || "",
        },
        validate: {
            url: (value) => (/^http(s)?:\/\//.test(value) ? null : 'Invalid URL'),
            secret: (value) => (value.length === 0 ? "Please enter a signing secret" : null),
        },
    });


    return (
        <form onSubmit={form.onSubmit(saveNodeData)}>
            <Group justify="space-between" align="center" mb="lg">
                <Text size="lg" fw={500}>Webhook</Text>
                <IconX size={18} onClick={() => setOpenedNodeId("")} cursor="pointer" />
            </Group>
            <Stack gap="lg" >
                <TextInput
                    leftSection={<IconHttpPost size={16} color="#E0A526" />}
                    withAsterisk
                    label="Target URL"
                    type="url"
                    description={(
                        <Group align="center" justify="space-between" p={0} m={0}>
                            A post request will be sent to this URL
                            <UnstyledButton variant="subtle" onClick={() => setOpenSampleReq(true)} p={0} m={0}>
                                <Text >View sample request</Text>
                            </UnstyledButton>
                        </Group>
                    )}
                    placeholder="Enter the webhook url"
                    {...form.getInputProps('url')}
                />

                <TextInput
                    leftSection={<IconLock size={16} />}
                    withAsterisk
                    label="Signing Secret"
                    description="Signing secret for the webhook"
                    placeholder="Enter the signing secret"
                    {...form.getInputProps('secret')}
                />
            </Stack>
            <Group justify="flex-end" mt="xl">
                <Button type="submit" leftSection={<IconClockBolt size={16} />} >Save</Button>
            </Group>
            <Modal
                transitionProps={{ transition: 'pop' }}
                opened={openSampleReq}
                title={"Sample request"}
                onClose={() => setOpenSampleReq(false)}
                size="xl"
            >
                <CodeHighlight p="md" code={getSampleCurl(form.values.url)} language="cmd" />
            </Modal>
        </form>
    )
}

export default WebhookNodeForm


