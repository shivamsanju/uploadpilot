export type Workflow = {
    id?: string;
    name: string;
    description: string;
    tags: string[];
    importPolicyId?: string;
    dataStoreId?: string;
    updatedAt?: number;
    createdAt?: number;
    createdBy?: string;
    updatedBy?: string;
}


export type CreateWorkflowForm = {
    name: string;
    description: string;
    tags: string[];
    importPolicyId: string;
    connectorId: string;
    dataStoreId: string;
    dataStoreName: string;
    bucket: string;
}