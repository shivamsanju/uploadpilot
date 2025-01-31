import { Box, Button, Group, Stack } from "@mantine/core";
import WebhookNodeForm from "./WebhookNodeForm";
import { useCanvas } from "../../../context/EditorCtx";

export type NodeFormProps = {
    nodeData: any,
    saveNodeData: (data: any) => void
    setOpenedNodeId: (id: string) => void
};


const getNodeForm = (key: string): React.FC<NodeFormProps> => {
    switch (key) {
        case 'webhook':
            return WebhookNodeForm
        default:
            return () => <>Not found</>
    }
}

export const NodeForm = () => {
    const { nodes, setNodes, openedNodeId, setOpenedNodeId, isUpdating, isPending, handleSave, handleDiscard } = useCanvas();

    const saveNodeData = (data: any) => {
        setNodes((nds: any) => nds.map((node: any) => {
            if (node.id === openedNodeId) {
                return {
                    ...node,
                    data: {
                        ...node.data,
                        ...data
                    }
                };
            }
            return node;
        }));
        setOpenedNodeId('');
    }

    const key = nodes?.find((node) => node.id === openedNodeId)?.key;
    const nodeData = nodes?.find((node) => node.id === openedNodeId)?.data;

    const Form = getNodeForm(key);
    return (
        <Stack justify="space-between" h="100%" p="sm" px="md">
            <Box >
                <Form nodeData={nodeData} saveNodeData={saveNodeData} setOpenedNodeId={setOpenedNodeId} />
            </Box>
            <Group justify="center" >
                <Button variant="default" c="dimmed" loading={isUpdating || isPending} onClick={handleDiscard}>Discard</Button>
                <Button loading={isUpdating || isPending} onClick={handleSave}>Save</Button>
            </Group>
        </Stack>
    )
}

