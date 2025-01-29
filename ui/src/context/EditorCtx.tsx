import { createContext, useEffect, useState } from 'react';
import { addEdge, applyNodeChanges, applyEdgeChanges } from '@xyflow/react';
import { useParams } from 'react-router-dom';
import { useGetProcessor, useUpdateProcessorTaskMutation } from '../apis/processors';

type ProcEditorContextType = {
    workspaceId: string | undefined;
    processorId: string | undefined;
    isPending: boolean;
    isUpdating: boolean;
    error: Error | null;
    nodes: any[];
    edges: any[];
    openedNodeId: string;
    onNodesChange: (changes: any) => void;
    onEdgesChange: (changes: any) => void;
    onConnect: (connection: any) => void;
    setNodes: (nodes: any) => void;
    setEdges: (edges: any) => void;
    handleSave: () => Promise<void>;
    handleDiscard: () => Promise<void>;
    setOpenedNodeId: (id: string) => void
}

export const ProcEditorContext = createContext<ProcEditorContextType>({
    nodes: [],
    edges: [],
    workspaceId: undefined,
    processorId: undefined,
    isPending: false,
    isUpdating: false,
    error: null,
    openedNodeId: '',
    onNodesChange: (changes: any) => { },
    onEdgesChange: (changes: any) => { },
    onConnect: (connection: any) => { },
    setNodes: (nodes: any) => { },
    setEdges: (edges: any) => { },
    handleSave: async () => { },
    handleDiscard: async () => { },
    setOpenedNodeId: (id: string) => { }
});

export const ProcEditorProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { workspaceId, processorId } = useParams();

    const { isPending, error, processor, invalidate } = useGetProcessor(workspaceId as string, processorId as string);
    const { mutateAsync, isPending: isUpdating } = useUpdateProcessorTaskMutation();

    const [nodes, setNodes] = useState<any[]>([]);
    const [edges, setEdges] = useState<any[]>([]);
    const [openedNodeId, setOpenedNodeId] = useState('');


    const onNodesChange = (changes: any) => {
        setNodes((prevNodes) => applyNodeChanges(changes, prevNodes));
    };

    const onEdgesChange = (changes: any) => {
        setEdges((prevEdges) => applyEdgeChanges(changes, prevEdges));
    };

    const onConnect = (connection: any) => {
        setEdges((prevEdges) => addEdge(connection, prevEdges));
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

    useEffect(() => {
        if (processor && processor.tasks) {
            setNodes(processor.tasks.nodes);
            setEdges(processor.tasks.edges);
        }
    }, [processor]);

    return (
        <ProcEditorContext.Provider
            value={{
                nodes,
                edges,
                workspaceId,
                processorId,
                isPending,
                isUpdating,
                error,
                openedNodeId,
                onNodesChange,
                onEdgesChange,
                onConnect,
                setNodes,
                setEdges,
                handleSave,
                handleDiscard,
                setOpenedNodeId
            }}
        >
            {children}
        </ProcEditorContext.Provider>
    );
};
