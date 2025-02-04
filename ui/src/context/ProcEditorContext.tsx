import { createContext, useContext } from "react";

export type ProcEditorContextType = {
  workspaceId: string | undefined;
  processorId: string | undefined;
  processor: any;
  isPending: boolean;
  isUpdating: boolean;
  error: Error | null;
  nodes: any[];
  edges: any[];
  nodesData: any;
  openedNodeId: string;
  connectionStateNodeId: any;
  openedBlocksModal: boolean;
  onNodesChange: (changes: any) => void;
  onEdgesChange: (changes: any) => void;
  onNodesDelete: (nodes: any) => void;
  setNodes: (nodes: any) => void;
  setEdges: (edges: any) => void;
  setNodesData: (nodesData: any) => void;
  handleSave: () => Promise<void>;
  handleDiscard: () => Promise<void>;
  setOpenedNodeId: (id: string) => void;
  setconnectionStateNodeId: (connectionStateNodeId: any) => void;
  onSelectNewNode: (item: any, type: string) => void;
  onConnectEnd: (fromId: any) => void;
  openBlocksModal: () => void;
  closeBlocksModal: () => void;
  getLayoutedElements: (nodes: any[], edges: any[], options: any) => any;
};

export const ProcEditorContext = createContext<ProcEditorContextType>({
  nodes: [],
  edges: [],
  nodesData: {},
  workspaceId: undefined,
  processorId: undefined,
  processor: {},
  isPending: false,
  isUpdating: false,
  error: null,
  openedNodeId: "",
  connectionStateNodeId: null,
  openedBlocksModal: false,
  onNodesChange: () => {},
  onEdgesChange: () => {},
  onNodesDelete: () => {},
  setNodes: () => {},
  setEdges: () => {},
  setNodesData: () => {},
  handleSave: async () => {},
  handleDiscard: async () => {},
  setOpenedNodeId: () => {},
  setconnectionStateNodeId: () => {},
  onConnectEnd: () => {},
  onSelectNewNode: () => {},
  openBlocksModal: () => {},
  closeBlocksModal: () => {},
  getLayoutedElements: () => {},
});

export const useCanvas = () => useContext(ProcEditorContext);
