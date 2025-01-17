export type ConnectorType = "s3" | "gcs" | "azure";
export type Connector = {
    id?: string;
    name: string;
    type: ConnectorType;
    tags?: string[];
    s3Region?: string;
    s3AccessKey?: string;
    s3SecretKey?: string;
    gcsApiKey?: string;
    azureAccountName?: string;
    azureAccountKey?: string;
    updatedAt?: number;
    createdAt?: number;
    createdBy?: string;
    updatedBy?: string;
}

