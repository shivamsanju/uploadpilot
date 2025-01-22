export type Webhook = {
    id?: string;
    url: string;
    signingSecret: string;
    event: string;
    method: string;
    workspaceId: string;
    enabled?: boolean;
    createdAt?: string;
    createdBy?: string;
    updatedAt?: string;
    updatedBy?: string;
};