import { Button, Group } from "@mantine/core";
import { useReactFlow } from "@xyflow/react";
import { useCanvas } from "../../context/ProcEditorContext";

export const Toolbar = () => {
  const {
    nodes,
    setNodes,
    edges,
    setEdges,
    isUpdating,
    isPending,
    handleSave,
    getLayoutedElements,
  } = useCanvas();
  const { fitView } = useReactFlow();

  const onFormat = () => {
    const { nodes: n, edges: e } = getLayoutedElements(nodes, edges, {});
    setNodes(n);
    setEdges(e);
    fitView({ padding: 1 });
  };

  return (
    <Group justify="center">
      <Button variant="default" c="dimmed" onClick={() => onFormat()}>
        Format
      </Button>
      {/* <Button variant="default" c="dimmed" loading={isUpdating || isPending} onClick={handleDiscard}>Discard</Button> */}
      <Button loading={isUpdating || isPending} onClick={handleSave}>
        Save
      </Button>
    </Group>
  );
};
