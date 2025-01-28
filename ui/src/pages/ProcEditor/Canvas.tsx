import { useMemo, useRef, useState } from 'react';
import { ReactFlow, Controls, Background, BackgroundVariant, Panel } from '@xyflow/react';
import { Transition, useMantineColorScheme, useMantineTheme } from '@mantine/core';
import { BaseNode } from './Node/BaseNode';
import '@xyflow/react/dist/style.css';
import SampleNodeForm from './Form';
import { useCanvas, useDragAndDrop } from '../../hooks/DndCanvas';

export const ProcessorCanvas = () => {
  const reactFlowWrapper = useRef(null);
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();
  const bg = colorScheme === "dark" ? "#0A0A0A" : theme.colors.gray[0];
  const nodeTypes = useMemo(() => ({ baseNode: BaseNode }), []);

  const [activeNode, setActiveNode] = useState("");
  const { nodes, edges, onConnect, onEdgesChange, onNodesChange } = useCanvas();
  const { onDrop, onDragOver } = useDragAndDrop();

  return (
    <div style={{ width: '100%', height: '92vh' }} ref={reactFlowWrapper} className='reactflow-wrapper'>
      <ReactFlow
        fitView
        fitViewOptions={{ padding: 3 }}
        style={{ background: bg }}
        colorMode={colorScheme === "auto" ? "dark" : colorScheme}
        nodeTypes={nodeTypes}
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onDrop={onDrop}
        onDragOver={onDragOver}
        onNodeDoubleClick={(e) => {
          console.log(e)
          setActiveNode("e.id");
        }}
      >
        <Controls />
        <Panel position="top-right">
          <Transition mounted={activeNode !== ""} transition="pop" duration={100} timingFunction="ease">
            {(styles) => <div style={styles} className="transition" >
              <SampleNodeForm workspaceId={""} setActive={setActiveNode} />
            </div>}
          </Transition>
        </Panel>
        <Background variant={BackgroundVariant.Dots} gap={12} size={1} />
      </ReactFlow>
    </div>
  );
}