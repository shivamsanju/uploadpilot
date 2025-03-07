export class Uploader {
  private apiKey: string;
  private tenantId: string;
  private workspaceId: string;
  private baseUrl: string;

  constructor(
    apiKey: string,
    tenantId: string,
    workspaceId: string,
    baseUrl: string = "http://localhost:8080"
  ) {
    this.apiKey = apiKey;
    this.tenantId = tenantId;
    this.workspaceId = workspaceId;
    this.baseUrl = baseUrl;
  }

  async upload(
    file: File,
    metadata: Record<string, any> = {}
  ): Promise<boolean> {
    if (!file) {
      throw new Error("No file provided for upload.");
    }

    const uploadId = await this.getPresignedUrl(file, metadata);
    await this.uploadToS3(file, uploadId.uploadUrl, uploadId.method);
    return this.completeUpload(uploadId.uploadId);
  }

  private async getPresignedUrl(
    file: File,
    metadata: Record<string, any>
  ): Promise<{ uploadUrl: string; method: string; uploadId: string }> {
    const uploadUrlEndpoint = `${this.baseUrl}/tenants/${this.tenantId}/workspaces/${this.workspaceId}/uploads`;

    const response = await fetch(uploadUrlEndpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Api-Key": this.apiKey,
      },
      body: JSON.stringify({
        fileName: file.name,
        contentType: file.type,
        contentLength: file.size,
        uploadUrlValiditySecs: 900,
        metadata,
      }),
    });

    if (!response.ok) {
      throw new Error("Failed to get upload URL.");
    }

    return response.json();
  }

  private async uploadToS3(
    file: File,
    uploadUrl: string,
    method: string
  ): Promise<void> {
    const uploadResponse = await fetch(uploadUrl, {
      method,
      body: file,
      headers: {
        "Content-Type": file.type,
        "if-none-match": "*",
      },
    });

    if (!uploadResponse.ok) {
      throw new Error("Upload to S3 failed.");
    }
  }

  private async completeUpload(uploadId: string): Promise<boolean> {
    if (!uploadId) {
      throw new Error("No uploadId provided.");
    }

    const completeEndpoint = `${this.baseUrl}/tenants/${this.tenantId}/workspaces/${this.workspaceId}/uploads/${uploadId}/complete`;

    const response = await fetch(completeEndpoint, {
      method: "POST",
      headers: {
        "X-Api-Key": this.apiKey,
      },
    });

    if (!response.ok) {
      throw new Error("Failed to complete upload.");
    }

    return true;
  }
}
