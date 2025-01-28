import { useForm } from "@mantine/form";
import { Button, Group, LoadingOverlay, Modal, Paper, Select, SelectProps, Stack, Text, TextInput, UnstyledButton } from "@mantine/core";
import { IconCheck, IconCircleCheck, IconClockBolt, IconExclamationCircle, IconHttpPost, IconLock, IconTrash, IconX } from "@tabler/icons-react";
import { useState } from "react";
import { CodeHighlight } from "@mantine/code-highlight";
import { useCreateWebhookMutation, useUpdateWebhookMutation } from "../../../apis/webhooks";

const EVENTS = ["File uploaded", "File upload failed", "File deleted"];


const iconProps = {
    stroke: 1.5,
    opacity: 0.6,
    size: 18,
};

const icons: Record<string, React.ReactNode> = {
    "File upload failed": <IconExclamationCircle color="#E0A526" {...iconProps} />,
    "File uploaded": <IconCircleCheck color="green" {...iconProps} />,
    "File deleted": <IconTrash color="red" {...iconProps} />,
};

const renderSelectOption: SelectProps['renderOption'] = ({ option, checked }) => (
    <Group flex="1" gap="xs">
        {icons[option.value]}
        {option.label}
        {checked && <IconCheck style={{ marginInlineStart: 'auto' }} {...iconProps} />}
    </Group>
);



type Props = {
    workspaceId: string,
    setActive: React.Dispatch<React.SetStateAction<string>>
};

const SampleNodeForm: React.FC<Props> = ({ workspaceId, setActive }) => {
    const [openSampleReq, setOpenSampleReq] = useState(false);

    const form = useForm({
        initialValues: {
            url: "",
            event: "",
            signingSecret: "",
            enabled: true
        },
        validate: {
            url: (value) => (/^http(s)?:\/\//.test(value) ? null : 'Invalid URL'),
            event: (value) => (value.length === 0 ? "Please select an event" : null),
            signingSecret: (value) => (value.length === 0 ? "Please enter a signing secret" : null),
        },
    });

    const { mutateAsync: createMutateAsync, isPending: isCreating } = useCreateWebhookMutation();
    const { mutateAsync: updateMutateAsync, isPending: isUpdating } = useUpdateWebhookMutation();

    const handleAdd = async (values: any) => {
        const body = {
            workspaceId, webhook: {
                method: "POST",
                url: values.url,
                enabled: values.enabled,
                event: values.event,
                signingSecret: values.signingSecret,
                workspaceId
            }
        }
    }

    return (
        <Paper w="30vw" miw="400" radius={"md"} p="md" h="88vh" withBorder>
            <form onSubmit={form.onSubmit(handleAdd)}>
                <LoadingOverlay visible={isCreating || isUpdating} overlayProps={{ radius: "sm", blur: 1 }} />
                <Group justify="space-between" align="center" mb="lg">
                    <Text size="lg" fw={500}>Add Webhook</Text>
                    <IconX size={18} onClick={() => setActive("")} cursor="pointer" />
                </Group>
                <Stack gap="xl">
                    <Select
                        size="sm"
                        withAsterisk
                        leftSection={<IconClockBolt size={16} />}
                        label="Events"
                        description="Event to trigger the webhook"
                        placeholder="Select event"
                        data={EVENTS}
                        {...form.getInputProps('event')}
                        renderOption={renderSelectOption}

                    />
                    <TextInput
                        size="sm"
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
                        size="sm"
                        leftSection={<IconLock size={16} />}
                        withAsterisk
                        label="Signing Secret"
                        description="Signing secret for the webhook"
                        placeholder="Enter the signing secret"
                        {...form.getInputProps('signingSecret')}
                    />
                </Stack>
            </form>
        </Paper>
    )
}

export default SampleNodeForm


