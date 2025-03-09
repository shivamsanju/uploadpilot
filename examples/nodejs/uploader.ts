export class NodeFile extends Blob {
  name: string;
  lastModified: number;

  constructor(buffer: Buffer, name: string, type = "application/octet-stream") {
    super([buffer], { type });
    this.name = name;
    this.lastModified = Date.now();
  }
}

export class Uploader {
  private apiKey: string;
  private tenantId: string;
  private workspaceId: string;
  private baseUrl: string;

  constructor(
    tenantId: string,
    workspaceId: string,
    apiKey: string,
    baseUrl: string = "http://localhost:8080"
  ) {
    this.tenantId = tenantId;
    this.workspaceId = workspaceId;
    this.apiKey = apiKey;
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

  async uploadMultiple(
    files: File[] | NodeFile[],
    metadata: Record<string, any> = {}
  ): Promise<boolean[]> {
    if (!files.length) {
      throw new Error("No files provided for upload.");
    }

    return Promise.all(files.map((file) => this.upload(file, metadata)));
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
      const resp = await response.text();
      throw new Error("Failed to get upload URL: " + resp);
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
      const resp = await uploadResponse.text();
      throw new Error("Failed to upload: " + resp);
    }
  }

  private async completeUpload(uploadId: string): Promise<boolean> {
    if (!uploadId) {
      throw new Error("No uploadId provided.");
    }

    const completeEndpoint = `${this.baseUrl}/tenants/${this.tenantId}/workspaces/${this.workspaceId}/uploads/${uploadId}/finish`;

    const response = await fetch(completeEndpoint, {
      method: "POST",
      headers: {
        "X-Api-Key": this.apiKey,
      },
    });

    if (!response.ok) {
      const err = await response.text();
      throw new Error("Failed to complete upload: " + err);
    }

    return true;
  }
}
