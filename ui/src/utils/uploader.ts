export class Uploader {
  private file: File;
  private tenantId: string;
  private workspaceId: string;
  private apiKey: string;
  private baseUrl: string;
  private metadata: Record<string, any>;
  private xhr: XMLHttpRequest | null = null;
  private uploadId: string | null = null;
  private isUploading = false;
  private abortedErr = 'upload_aborted';
  private status = 'unknown';

  private eventHandlers = {
    progress: (computable: boolean, loaded: number, total: number) => {},
    start: () => {},
    complete: (serverFileId: string) => {},
    error: (message: string) => {},
    abort: () => {},
  };

  constructor(
    file: File,
    tenantId: string,
    workspaceId: string,
    apiKey: string,
    metadata: Record<string, any> = {},
    baseUrl: string = 'http://localhost:8080',
  ) {
    this.file = file;
    this.tenantId = tenantId;
    this.workspaceId = workspaceId;
    this.apiKey = apiKey;
    this.baseUrl = baseUrl;
    this.metadata = metadata;
  }

  addMetadata(key: string, value: any) {
    this.metadata[key] = value;
  }

  removeMetadata(key: string) {
    delete this.metadata[key];
  }

  async start() {
    if (this.isUploading) return;
    this.isUploading = true;
    this.eventHandlers.start();

    try {
      this.status = 'generating_presigned_url';
      const presignedData = await this.getPresignedUrl();

      this.status = 'uploading';
      this.uploadId = presignedData.uploadId;

      await this.uploadToS3(presignedData.uploadUrl, presignedData.method);

      this.status = 'finishing';
      await this.completeUpload();
      this.status = 'complete';
      this.eventHandlers.complete(this.uploadId!);
    } catch (error) {
      if ((error as Error).message === this.abortedErr) return;
      this.eventHandlers.error((error as Error).message);
      try {
        if (this.status === 'uploading') {
          await this.completeUpload('Failed');
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  cancel() {
    if (this.xhr) {
      this.xhr.abort();
    }
    this.eventHandlers.abort();
    this.isUploading = false;
    this.completeUpload('Cancelled');
  }

  private async getPresignedUrl() {
    const url = `${this.baseUrl}/tenants/${this.tenantId}/workspaces/${this.workspaceId}/uploads`;

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Api-Key': this.apiKey,
      },
      body: JSON.stringify({
        fileName: this.file.name,
        contentType: this.file.type,
        contentLength: this.file.size,
        uploadUrlValiditySecs: 900,
        metadata: this.metadata,
      }),
    });

    if (!response.ok) {
      const resp = await response.json();
      throw new Error(resp.message);
    }

    return response.json();
  }

  private async uploadToS3(uploadUrl: string, method: string) {
    return new Promise<void>((resolve, reject) => {
      this.xhr = new XMLHttpRequest();
      this.xhr.open(method, uploadUrl);

      this.xhr.upload.onprogress = event => {
        this.eventHandlers.progress(
          event.lengthComputable,
          event.loaded,
          event.total,
        );
      };

      this.xhr.onload = () => {
        if (this.xhr!.status >= 200 && this.xhr!.status < 300) {
          resolve();
        } else {
          reject(new Error(`Upload failed with status ${this.xhr!.status}`));
        }
      };

      this.xhr.onerror = () =>
        reject(new Error('Network error during upload.'));
      this.xhr.onabort = () => reject(new Error(this.abortedErr));

      this.xhr.setRequestHeader('Content-Type', this.file.type);
      this.xhr.setRequestHeader('If-None-Match', '*');
      this.xhr.send(this.file);
    });
  }

  private async completeUpload(status: string = 'Finished') {
    if (!this.uploadId) throw new Error('No uploadId provided.');

    const url = `${this.baseUrl}/tenants/${this.tenantId}/workspaces/${this.workspaceId}/uploads/${this.uploadId}/finish`;

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'X-Api-Key': this.apiKey,
      },
      body: JSON.stringify({
        status: status,
      }),
    });

    if (!response.ok) {
      const resp = await response.json();
      throw new Error(resp.message);
    }
  }

  onProgress(
    callback: (computable: boolean, loaded: number, total: number) => void,
  ) {
    this.eventHandlers.progress = callback;
  }

  onComplete(callback: (serverFileId: string) => void) {
    this.eventHandlers.complete = callback;
  }

  onError(callback: (message: string) => void) {
    this.eventHandlers.error = callback;
  }

  onAbort(callback: () => void) {
    this.eventHandlers.abort = callback;
  }
}
