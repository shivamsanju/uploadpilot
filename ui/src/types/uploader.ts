export type WorkspaceConfig = {
  minFileSize?: number;
  maxFileSize?: number;
  allowedContentTypes?: string[];
  allowedOrigins?: string[];
  maxUploadURLLifetimeSecs?: number;
  requiredMetadataFields?: string[];
  allowPauseAndResume?: boolean;
  enableImageEditing?: boolean;
  useCompression?: boolean;
  useFaultTolerantMode?: boolean;
};
