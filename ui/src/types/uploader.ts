export type WorkspaceConfig = {
  minFileSize?: number;
  maxFileSize?: number;
  allowedContentTypes?: string[];
  allowedOrigins?: string[];
  maxUploadURLLifetimeSecs?: number;
  requiredMetadataFields?: string[];
};
