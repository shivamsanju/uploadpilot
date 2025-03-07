var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
export class Uploader {
    constructor(apiKey, tenantId, workspaceId, baseUrl = "http://localhost:8080") {
        this.apiKey = apiKey;
        this.tenantId = tenantId;
        this.workspaceId = workspaceId;
        this.baseUrl = baseUrl;
    }
    upload(file_1) {
        return __awaiter(this, arguments, void 0, function* (file, metadata = {}) {
            if (!file) {
                throw new Error("No file provided for upload.");
            }
            const uploadId = yield this.getPresignedUrl(file, metadata);
            yield this.uploadToS3(file, uploadId.uploadUrl, uploadId.method);
            return this.completeUpload(uploadId.uploadId);
        });
    }
    getPresignedUrl(file, metadata) {
        return __awaiter(this, void 0, void 0, function* () {
            const uploadUrlEndpoint = `${this.baseUrl}/tenants/${this.tenantId}/workspaces/${this.workspaceId}/uploads`;
            const response = yield fetch(uploadUrlEndpoint, {
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
        });
    }
    uploadToS3(file, uploadUrl, method) {
        return __awaiter(this, void 0, void 0, function* () {
            const uploadResponse = yield fetch(uploadUrl, {
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
        });
    }
    completeUpload(uploadId) {
        return __awaiter(this, void 0, void 0, function* () {
            if (!uploadId) {
                throw new Error("No uploadId provided.");
            }
            const completeEndpoint = `${this.baseUrl}/tenants/${this.tenantId}/workspaces/${this.workspaceId}/uploads/${uploadId}/complete`;
            const response = yield fetch(completeEndpoint, {
                method: "POST",
                headers: {
                    "X-Api-Key": this.apiKey,
                },
            });
            if (!response.ok) {
                throw new Error("Failed to complete upload.");
            }
            return true;
        });
    }
}
