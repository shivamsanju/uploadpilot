import { create } from 'zustand';
import { Workflow } from '../types/workflow';
import { Connector } from '../types/connector';
import { ImportPolicy } from '../types/importpolicy';



type Store = {
    // Workflows
    workflows: Workflow[];
    setWorkflows: (workflows: Workflow[]) => void;

    // Connectors
    connectors: Connector[];
    setConnectors: (connectors: Connector[]) => void;

    // Import policies
    importPolicies: ImportPolicy[];
    setImportPolicies: (importPolicies: ImportPolicy[]) => void;

    // Import Policy Details
    importPolicyDetails: ImportPolicy | null;
    setImportPolicyDetails: (importPolicyDetails: ImportPolicy | null) => void;

};


const useAppStore = create<Store>((set) => ({
    // Workflows
    workflows: [],
    setWorkflows: (workflows) => set({ workflows }),

    // Connectors
    connectors: [],
    setConnectors: (connectors) => set({ connectors }),

    // Import policies
    importPolicies: [],
    setImportPolicies: (importPolicies) => set({ importPolicies }),

    // Import Policy Details
    importPolicyDetails: null,
    setImportPolicyDetails: (importPolicyDetails) => set({ importPolicyDetails }),
}));

export default useAppStore;
