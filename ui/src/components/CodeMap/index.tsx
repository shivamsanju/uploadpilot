import ReactFlow, {
    ReactFlowProvider,
    Background,
    Controls,
    Handle,
    Position,
    useNodesState,
    useEdgesState,
} from 'reactflow';
import 'reactflow/dist/style.css';
import { Card, Text } from '@mantine/core';
import { useEffect } from 'react';
import axios from 'axios';
import { notifications } from '@mantine/notifications';
import { getApiDomain } from '../../config';
import getElkLayout from './ElkLayout';
import CodeNode from '../CustomNodes/CodeNode';


// Define the node types
const nodeTypes = {
    mantineNode: CodeNode,
};


const CodeMap = ({ codebaseId }: any) => {
    const [nodes, setNodes, onNodesChange] = useNodesState([]);
    const [edges, setEdges, onEdgesChange] = useEdgesState([]);

    const getCodebaseMap = async () => {
        try {
            const resp = await axios.get(getApiDomain() + "/codebase/" + codebaseId + "/codemap");
            return resp.data;
        } catch {
            notifications.show({
                title: 'Error',
                message: 'Failed to fetch codemap',
                color: 'red',
            })
        }
    }


    useEffect(() => {
        getCodebaseMap().then((data) => {
            getElkLayout(data.nodes, data.edges, "DOWN").then(({ nodes, edges }) => {
                setNodes(nodes);
                setEdges(edges)
            });
            ;
        });
    }, [])

    return (
        <ReactFlowProvider>
            <div style={{ height: '100vh', width: '100%' }}>
                <ReactFlow
                    nodes={nodes}
                    edges={edges}
                    onNodesChange={onNodesChange}
                    onEdgesChange={onEdgesChange}
                    nodeTypes={nodeTypes}
                    fitView
                    minZoom={0.1}
                >
                    <Background color="#aaa" gap={16} />
                    <Controls />
                </ReactFlow>
            </div>
        </ReactFlowProvider>
    );
};

export default CodeMap;
