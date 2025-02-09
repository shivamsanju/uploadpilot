import { useMemo, useRef } from "react";
import {
  ReactFlow,
  Controls,
  Background,
  BackgroundVariant,
  ConnectionLineType,
  Panel,
} from "@xyflow/react";
import { useMantineColorScheme, useMantineTheme } from "@mantine/core";
import { BaseNode } from "../../components/EditorNode/BaseNode";
import { BlockSearch } from "../../components/BlockSearch";
import "@xyflow/react/dist/style.css";
import { useCanvas } from "../../context/ProcEditorContext";
import { NodeForm } from "./Form";

export const ProcessorCanvas = () => {
  const reactFlowWrapper = useRef(null);
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();
  const bg = colorScheme === "dark" ? "#0A0A0A" : theme.colors.gray[0];
  const nodeTypes = useMemo(() => ({ baseNode: BaseNode }), []);
  const {
    nodes,
    edges,
    onEdgesChange,
    onNodesChange,
    onNodesDelete,
    openedBlocksModal,
    closeBlocksModal,
  } = useCanvas();

  return (
    <div
      style={{ width: "100%", height: "92vh" }}
      ref={reactFlowWrapper}
      className="reactflow-wrapper"
    >
      <ReactFlow
        fitView
        fitViewOptions={{ padding: 1 }}
        style={{ background: bg }}
        colorMode={colorScheme === "auto" ? "dark" : colorScheme}
        nodeTypes={nodeTypes}
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodesDelete={onNodesDelete}
        connectionLineType={ConnectionLineType.SmoothStep}
      >
        <Controls />
        <Background variant={BackgroundVariant.Dots} gap={12} size={1} />
        <Panel position="top-right" style={{ zIndex: 9999 }}>
          <NodeForm />
        </Panel>
      </ReactFlow>
      <BlockSearch opened={openedBlocksModal} close={closeBlocksModal} />
    </div>
  );
};
