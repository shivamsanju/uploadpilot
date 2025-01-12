export type UploaderConfig = {
    minFileSize?: number;
    maxFileSize?: number;
    minNumberOfFiles?: number;
    maxNumberOfFiles?: number;
    maxTotalFileSize?: number;
    allowedFileTypes: string[];
    allowedSources: string[];
    requiredMetadataFields?: string[];
    theme?: 'dark' | 'light' | 'auto';
    showStatusBar?: boolean;
    showProgressBar?: boolean;
    allowPauseAndResume?: boolean;
    enableImageEditing?: boolean;
    useCompression?: boolean;
    useFaultTolerantMode?: boolean;
};

export type Uploader = {
    id?: string;
    name: string;
    description: string;
    tags: string[];
    config?: UploaderConfig;
    dataStoreId?: string;
    updatedAt?: number;
    createdAt?: number;
    createdBy?: string;
    updatedBy?: string;
}

export type CreateUploaderForm = {
    name: string;
    description: string;
    tags: string[];
    connectorId: string;
    dataStoreId: string;
    dataStoreName: string;
    bucket: string;

} & UploaderConfig;