import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { applyNodeChanges, applyEdgeChanges, getIncomers, getOutgoers, getConnectedEdges, ConnectionLineType } from '@xyflow/react';
import { useParams } from 'react-router-dom';
import { useGetProcessor, useUpdateProcessorTaskMutation } from '../apis/processors';
import { v4 as uuid } from 'uuid';
import { useDisclosure } from '@mantine/hooks';
import { stratify, tree } from 'd3-hierarchy';
import { isDataNode } from '../components/EditorNode/IsDataNode';
type ProcEditorContextType = {
    workspaceId: string | undefined;
    processorId: string | undefined;
    processor: any;
    isPending: boolean;
    isUpdating: boolean;
    error: Error | null;
    nodes: any[];
    edges: any[];
    openedNodeId: string;
    onNodesChange: (changes: any) => void;
    onEdgesChange: (changes: any) => void;
    onNodesDelete: (nodes: any) => void;
    setNodes: (nodes: any) => void;
    setEdges: (edges: any) => void;
    handleSave: () => Promise<void>;
    handleDiscard: () => Promise<void>;
    setOpenedNodeId: (id: string) => void
    connectionStateNodeId: any;
    setconnectionStateNodeId: (connectionStateNodeId: any) => void;
    onSelectNewNode: (item: any, type: string) => void
    onConnectEnd: (fromId: any) => void
    openedBlocksModal: boolean,
    openBlocksModal: () => void
    closeBlocksModal: () => void
    getLayoutedElements: (nodes: any[], edges: any[], options: any) => any
}

export const ProcEditorContext = createContext<ProcEditorContextType>({
    nodes: [],
    edges: [],
    workspaceId: undefined,
    processorId: undefined,
    processor: {},
    isPending: false,
    isUpdating: false,
    error: null,
    openedNodeId: '',
    onNodesChange: (changes: any) => { },
    onEdgesChange: (changes: any) => { },
    onNodesDelete: (nodes: any) => { },
    setNodes: (nodes: any) => { },
    setEdges: (edges: any) => { },
    handleSave: async () => { },
    handleDiscard: async () => { },
    setOpenedNodeId: (id: string) => { },
    connectionStateNodeId: null,
    setconnectionStateNodeId: () => { },
    onConnectEnd: (fromId) => { },
    onSelectNewNode: (item, type) => { },
    openedBlocksModal: false,
    openBlocksModal: () => { },
    closeBlocksModal: () => { },
    getLayoutedElements: (nodes, edges, options) => { },

});

export const ProcEditorProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { workspaceId, processorId } = useParams();
    const [openedBlocksModal, { open: openBlocksModal, close: closeBlocksModal }] = useDisclosure();
    const { isPending, error, processor, invalidate } = useGetProcessor(workspaceId as string, processorId as string);
    const { mutateAsync, isPending: isUpdating } = useUpdateProcessorTaskMutation();

    const [nodes, setNodes] = useState<any[]>([]);
    const [edges, setEdges] = useState<any[]>([]);
    const [openedNodeId, setOpenedNodeId] = useState('');
    const [connectionStateNodeId, setconnectionStateNodeId] = useState<any>(null);
    const g = tree();

    const getLayoutedElements = useCallback((nodes: any, edges: any, options: any) => {
        if (nodes.length === 0) return { nodes, edges };
        if (!document) return { nodes, edges };
        const { width, height }: any = document.querySelector(`[data-id="${nodes[0].id}"]`)?.getBoundingClientRect();
        const hierarchy = stratify()
            .id((node: any) => node.id)
            .parentId((node: any) => edges.find((edge: any) => edge.target === node.id)?.source);
        const root = hierarchy(nodes);
        const layout = g.nodeSize([width * 2, height * 5])(root);

        return {
            nodes: layout
                .descendants()
                .map((node: any) => ({ ...nodes.find((n: any) => n.id === node.id), position: { x: node.x, y: node.y } })),
            edges,
        };
    }, [g]);

    const onNodesChange = (changes: any) => {
        setNodes((prevNodes) => applyNodeChanges(changes, prevNodes));
    };

    const onEdgesChange = (changes: any) => {
        setEdges((prevEdges) => applyEdgeChanges(changes, prevEdges));
    };

    const handleSave = async () => {
        try {
            await mutateAsync({
                processorId: processorId!,
                workspaceId: workspaceId!,
                tasks: {
                    nodes,
                    edges
                }
            })
        } catch (error) {
            console.log(error)
        }
    }

    const handleDiscard = async () => {
        await invalidate();
    }

    const onNodesDelete = useCallback(
        (deleted: any) => {
            setEdges(
                deleted.reduce((acc: any, node: any) => {
                    const incomers = getIncomers(node, nodes, edges);
                    const outgoers = getOutgoers(node, nodes, edges);
                    const connectedEdges = getConnectedEdges([node], edges);

                    const remainingEdges = acc.filter(
                        (edge: any) => !connectedEdges.includes(edge),
                    );

                    const createdEdges = incomers.flatMap(({ id: source }) =>
                        outgoers.map(({ id: target }) => ({
                            id: `${source}->${target}`,
                            deletable: false,
                            source,
                            target,
                        })),
                    );

                    return [...remainingEdges, ...createdEdges];
                }, edges),
            );
        },
        [nodes, edges],
    );

    const onSelectNewNode = useCallback(
        (item: any, type: string) => {
            const node = nodes.find((n: any) => n.id === connectionStateNodeId);
            const numEdges = edges.filter((e: any) => e.source === connectionStateNodeId).length;
            if (!node) return;

            const id = uuid();

            const newNode: any = {
                id: id,
                type,
                position: {
                    x: node?.position?.x + (350 * numEdges),
                    y: node?.position?.y + 200,
                },
                key: item?.key,
                retry: 0,
                continueOnError: false,
                timeoutMilSec: 1000000,
                data: {
                    ...item,
                    isComplete: isDataNode(item?.key) ? false : true
                },
                deletable: true,
            };
            setNodes((nds: any[]) => nds.concat(newNode));
            setEdges((eds: any[]) =>
                eds.concat({ id: uuid(), source: connectionStateNodeId, target: id, deletable: false }),
            );
            setconnectionStateNodeId(null);
        },
        [connectionStateNodeId, nodes, edges],
    );

    const onConnectEnd = useCallback((fromNodeId: any) => {
        setconnectionStateNodeId(fromNodeId);
    }, []);

    useEffect(() => {
        if (processor && processor.tasks) {
            setNodes(processor.tasks.nodes);
            setEdges(processor.tasks.edges);
        }
    }, [processor]);




    return (
        <ProcEditorContext.Provider
            value={{
                nodes: nodes,
                edges: edges.map((e: any) => ({ ...e, type: ConnectionLineType.SmoothStep, animated: true })),
                workspaceId,
                processorId,
                processor,
                isPending,
                isUpdating,
                error,
                openedNodeId,
                connectionStateNodeId,
                openedBlocksModal,
                closeBlocksModal,
                openBlocksModal,
                onNodesChange,
                onEdgesChange,
                onNodesDelete,
                setNodes,
                setEdges,
                handleSave,
                handleDiscard,
                setOpenedNodeId,
                setconnectionStateNodeId,
                onConnectEnd,
                onSelectNewNode,
                getLayoutedElements,
            }}
        >
            {children}
        </ProcEditorContext.Provider>
    );
};


export const useCanvas = () => useContext(ProcEditorContext);