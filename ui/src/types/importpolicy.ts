export type ImportPolicy = {
    id?: string;
    name: string;
    allowedMimeTypes: string[];
    allowedSources: string[];
    maxFileSizeKb: number;
    maxFileCount: number;
    updatedAt?: number;
    createdAt?: number;
    createdBy?: string;
    updatedBy?: string;
};