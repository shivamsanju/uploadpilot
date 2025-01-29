import { Box } from "@mantine/core";
import { useCanvas } from "../../../../hooks/DndCanvas";
import WebhookNodeForm from "./WebhookNodeForm";

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
    const { nodes, setNodes, openedNodeId, setOpenedNodeId } = useCanvas();

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
        <Box>
            <Form nodeData={nodeData} saveNodeData={saveNodeData} setOpenedNodeId={setOpenedNodeId} />
        </Box>
    )
}

