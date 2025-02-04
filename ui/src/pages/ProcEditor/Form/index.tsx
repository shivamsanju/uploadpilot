import { Alert, Group, Paper, Text } from "@mantine/core";
import WebhookNodeForm from "./WebhookNodeForm";
import { useCanvas } from "../../../context/ProcEditorContext";
import NoDataForm from "./NoDataForm";
import { IconLock, IconX } from "@tabler/icons-react";

export type NodeFormProps = {
  nodeData: any;
  saveNodeData: (data: any) => void;
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
  const {
    nodes,
    setNodes,
    openedNodeId,
    setOpenedNodeId,
    setNodesData,
    nodesData,
  } = useCanvas();

  const saveNodeData = (data: any) => {
    setNodes((nds: any) =>
      nds.map((node: any) => {
        if (node.id === openedNodeId) {
          return {
            ...node,
            data: {
              ...node.data,
              isComplete: true,
            },
          };
        }
        return node;
      })
    );
    setNodesData((nds: any) => ({ ...nds, [openedNodeId]: data }));
    setOpenedNodeId("");
  };

  const key = nodes?.find((node) => node.id === openedNodeId)?.key;
  const nodeData = nodesData[openedNodeId] || {};

  const Form = getNodeForm(key);

  if (!Form || !openedNodeId) return <></>;

  return (
    <Paper
      h="88vh"
      px="xl"
      py="sm"
      withBorder
      style={{ borderColor: "var(--mantine-color-appcolor-5)" }}
      radius="lg"
    >
      <Group justify="space-between" align="center" mb="lg">
        <Text size="lg" fw={500}>
          {key && key.charAt(0).toUpperCase() + key.slice(1)}
        </Text>
        <IconX size={18} onClick={() => setOpenedNodeId("")} cursor="pointer" />
      </Group>
      <Alert color="appcolor" icon={<IconLock size={20} />} mb="md">
        <Text size="sm">Your data is encrypted and stored safely.</Text>
      </Alert>
      <Form
        nodeData={nodeData}
        saveNodeData={saveNodeData}
        setOpenedNodeId={setOpenedNodeId}
      />
    </Paper>
  );
};
