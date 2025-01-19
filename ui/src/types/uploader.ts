export type UploaderConfig = {
    minFileSize?: number;
    maxFileSize?: number;
    minNumberOfFiles?: number;
    maxNumberOfFiles?: number;
    maxTotalFileSize?: number;
    allowedFileTypes: string[];
    allowedSources: string[];
    requiredMetadataFields?: string[];
    allowPauseAndResume?: boolean;
    enableImageEditing?: boolean;
    useCompression?: boolean;
    useFaultTolerantMode?: boolean;
};

export type Datastore = {
    name: string;
    connectorId: string;
    connectorName: string;
    connectorType: string;
    bucket: string;
}

export type Uploader = {
    id?: string;
    name: string;
    description: string;
    tags: string[];
    config?: UploaderConfig;
    dataStore?: Datastore;
    updatedAt?: number;
    createdAt?: number;
    createdBy?: string;
    updatedBy?: string;
}

export type CreateUploaderForm = {
    name: string;
    description: string;
    tags: string[];
    dataStoreName: string;
} & UploaderConfig & Datastore;