import { useForm } from "@mantine/form";
import { Button, Group, LoadingOverlay, Modal, Select, SelectProps, Stack, Text, TextInput, UnstyledButton } from "@mantine/core";
import { IconCheck, IconCircleCheck, IconClockBolt, IconExclamationCircle, IconHttpPost, IconLock, IconTrash } from "@tabler/icons-react";
import { useState } from "react";
import { CodeHighlight } from "@mantine/code-highlight";
import { useCreateWebhookMutation, useUpdateWebhookMutation } from "../../../../apis/webhooks";

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

type Props = {
    setOpened: React.Dispatch<React.SetStateAction<boolean>>,
    workspaceId: string,
    mode?: 'edit' | 'add' | 'view',
    setMode: React.Dispatch<React.SetStateAction<'edit' | 'add' | 'view'>>,
    initialValues?: any,
    setInitialValues: React.Dispatch<React.SetStateAction<any>>
};

const AddWebhookForm: React.FC<Props> = ({ setInitialValues, setMode, setOpened, workspaceId, mode = 'add', initialValues }) => {
    const [openSampleReq, setOpenSampleReq] = useState(false);

    const form = useForm({
        initialValues: ((mode === 'edit' || mode === 'view') && initialValues) ? initialValues : {
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
        try {
            if (mode === 'edit') {
                await updateMutateAsync({
                    webhookId: initialValues.id,
                    ...body
                });
            } else {
                await createMutateAsync(body);
            }
            setOpened(false);
            setInitialValues(null);
            setMode('add');
        } catch (error) {
            console.error(error);
        }
    }

    return (
        <form onSubmit={form.onSubmit(handleAdd)}>
            <LoadingOverlay visible={isCreating || isUpdating} overlayProps={{ radius: "sm", blur: 1 }} />
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
                    disabled={mode === 'view'}

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
                                <Text size="xs">View sample request</Text>
                            </UnstyledButton>
                        </Group>
                    )}
                    placeholder="Enter the webhook url"
                    {...form.getInputProps('url')}
                    disabled={mode === 'view'}
                />

                <TextInput
                    size="sm"
                    leftSection={<IconLock size={16} />}
                    withAsterisk
                    label="Signing Secret"
                    description="Signing secret for the webhook"
                    placeholder="Enter the signing secret"
                    {...form.getInputProps('signingSecret')}
                    disabled={mode === 'view'}
                />
            </Stack>
            {mode !== "view" && (
                <Group justify="flex-end" mt={50}>
                    <Button type="submit" size="sm" >
                        {mode === 'edit' ? 'Save' : 'Add'}
                    </Button>
                </Group>
            )}
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

export default AddWebhookForm


