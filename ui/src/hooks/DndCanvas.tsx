import { useCallback, useContext } from 'react';
import { useReactFlow } from '@xyflow/react';
import { v4 as uuid } from 'uuid';
import { DnDContext } from '../context/DnD';
import { ProcEditorContext } from '../context/EditorCtx';

export const useDragAndDrop = () => {
    const { type, setType, dataTransfer, setDataTransfer } = useContext(DnDContext);
    const { setNodes } = useCanvas();
    const { screenToFlowPosition } = useReactFlow();

    const onDragOver = useCallback((event: React.DragEvent) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'move';
    }, []);

    const onDrop = useCallback((event: React.DragEvent) => {
        event.preventDefault();
        const itemJson = event.dataTransfer.getData('item');
        const item = JSON.parse(itemJson || '{}');

        if (!type) return;

        const position = screenToFlowPosition({ x: event.clientX, y: event.clientY });
        const newNode: any = {
            id: uuid(),
            type,
            position,
            key: dataTransfer?.key,
            retry: 0,
            continueOnError: false,
            timeoutMilSec: 1000000,
            data: { label: item?.label, description: item?.description, ...dataTransfer },
        };

        setDataTransfer({});

        setNodes((nds: any) => nds.concat(newNode));
    }, [screenToFlowPosition, type, setNodes, setDataTransfer, dataTransfer]);

    return { onDragOver, onDrop, type, dataTransfer, setType, setDataTransfer };
};


export const useCanvas = () => useContext(ProcEditorContext)