import { Paper } from "@mantine/core";
import WebhookNodeForm from "./WebhookNodeForm";
import { useCanvas } from "../../../context/EditorCtx";
import NoDataForm from "./NoDataForm";

export type NodeFormProps = {
  nodeData: any;
  saveNodeData: (data: any) => void;
  setOpenedNodeId: (id: string) => void;
};

const getNodeForm = (key: string): any => {
  switch (key) {
    case "webhook":
      return WebhookNodeForm;
    default:
      return NoDataForm;
  }
};

export const NodeForm = () => {
  const { nodes, setNodes, openedNodeId, setOpenedNodeId } = useCanvas();

  const saveNodeData = (data: any) => {
    setNodes((nds: any) =>
      nds.map((node: any) => {
        if (node.id === openedNodeId) {
          return {
            ...node,
            data: {
              ...node.data,
              ...data,
              isComplete: true,
            },
          };
        }
        return node;
      }),
    );
    setOpenedNodeId("");
  };

  const key = nodes?.find((node) => node.id === openedNodeId)?.key;
  const nodeData = nodes?.find((node) => node.id === openedNodeId)?.data;

  const Form = getNodeForm(key);

  if (!Form || !openedNodeId) return <></>;

  return (
    <Paper h="88vh" px="xl" py="sm" withBorder>
      <Form
        nodeData={nodeData}
        saveNodeData={saveNodeData}
        setOpenedNodeId={setOpenedNodeId}
      />
    </Paper>
  );
};
