import { useForm } from "@mantine/form";
import { Button, Group, LoadingOverlay, MultiSelect, SelectProps, Stack, TextInput } from "@mantine/core";
import { IconCheck, IconClockBolt, IconFileUnknown } from "@tabler/icons-react";
import { MIME_TYPES } from "../../../utils/mime";
import { MIME_TYPE_ICONS } from "../../../utils/fileicons";
import { useCreateProcessorMutation, useUpdateProcessorMutation } from "../../../apis/processors";


const iconProps = {
    stroke: 1.5,
    opacity: 0.6,
    size: 14,
};

const renderSelectOption: SelectProps['renderOption'] = ({ option, checked }) => {
    let Icon = MIME_TYPE_ICONS[option.value]
    if (!Icon) {
        Icon = IconFileUnknown
    }
    return (
        <Group flex="1" gap="xs">
            <Icon />
            {option.label}
            {checked && <IconCheck style={{ marginInlineStart: 'auto' }} {...iconProps} />}
        </Group>
    )
};


type Props = {
    setOpened: React.Dispatch<React.SetStateAction<boolean>>,
    workspaceId: string,
    mode?: 'edit' | 'add' | 'view',
    setMode: React.Dispatch<React.SetStateAction<'edit' | 'add' | 'view'>>,
    initialValues?: any,
    setInitialValues: React.Dispatch<React.SetStateAction<any>>
};

const AddWebhookForm: React.FC<Props> = ({ setInitialValues, setMode, setOpened, workspaceId, mode = 'add', initialValues }) => {
    const form = useForm({
        initialValues: ((mode === 'edit' || mode === 'view') && initialValues) ? initialValues : {
            name: "",
            triggers: [],
            tasks: {
                nodes: [],
                edges: []
            },
            enabled: true
        },
        validate: {
            name: (value) => (value ? null : 'Name is required'),
            triggers: (value) => (value.length === 0 ? "Please select atleast one file type as a trigger" : null),
        },
    });

    const { mutateAsync: createMutateAsync, isPending: isCreating } = useCreateProcessorMutation();
    const { mutateAsync: updateMutateAsync, isPending: isUpdating } = useUpdateProcessorMutation();

    const handleAdd = async (values: any) => {
        try {
            if (mode === 'edit') {
                await updateMutateAsync({
                    processorId: initialValues.id,
                    workspaceId,
                    processor: {
                        name: values.name,
                        triggers: values.triggers,
                    }
                });
            } else {
                await createMutateAsync({
                    workspaceId,
                    processor: {
                        workspaceId,
                        name: values.name,
                        triggers: values.triggers,
                        tasks: values.tasks,
                        enabled: values.enabled
                    }
                });
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
                <TextInput
                    size="sm"
                    withAsterisk
                    label="Name"
                    description="Name of the processor"
                    type="name"
                    placeholder="Enter a name"
                    {...form.getInputProps('name')}
                    disabled={mode === 'view'}
                />
                <MultiSelect
                    searchable
                    size="sm"
                    withAsterisk
                    leftSection={<IconClockBolt size={16} />}
                    label="Trigger"
                    description="File type to trigger the processor"
                    placeholder="Select file type"
                    data={MIME_TYPES}
                    {...form.getInputProps('triggers')}
                    renderOption={renderSelectOption}
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
        </form>
    )
}

export default AddWebhookForm


