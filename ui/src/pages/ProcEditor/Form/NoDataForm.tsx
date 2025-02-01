import { useForm } from "@mantine/form";
import { Button, Group, Stack, Text, ScrollArea } from "@mantine/core";
import { IconDeviceFloppy, IconX } from "@tabler/icons-react";
import { NodeFormProps } from ".";
import { CommonForm } from "./Common";

const NoDataForm: React.FC<NodeFormProps> = ({ nodeData, saveNodeData, setOpenedNodeId }) => {

    const form = useForm({});


    return (
        <form onSubmit={form.onSubmit(saveNodeData)}>
            <Stack justify="space-between" h="82vh">
                <ScrollArea scrollbarSize={6}>
                    <Group justify="space-between" align="center" mb="lg">
                        <Text size="lg" fw={500}>Task</Text>
                        <IconX size={18} onClick={() => setOpenedNodeId("")} cursor="pointer" />
                    </Group>
                    <Stack gap="lg" pr="lg">
                        <CommonForm form={form} />
                    </Stack>
                </ScrollArea>
            </Stack>
            <Group justify="center">
                <Button type="submit" leftSection={<IconDeviceFloppy size={16} />} >Save</Button>
            </Group>
        </form>
    )
}

export default NoDataForm


