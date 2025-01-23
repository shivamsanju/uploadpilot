export type UploaderConfig = {
    minFileSize?: number;
    maxFileSize?: number;
    minNumberOfFiles?: number;
    maxNumberOfFiles?: number;
    allowedFileTypes: string[];
    allowedSources: string[];
    requiredMetadataFields?: string[];
    allowPauseAndResume?: boolean;
    enableImageEditing?: boolean;
    useCompression?: boolean;
    useFaultTolerantMode?: boolean;
    authEndpoint?: string;
};