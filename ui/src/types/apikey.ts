export type APIKey = {
  id: string;
  name: string;
  workspaceId: string;
  revoked: boolean;
  createdBy: string;
  createdAt: Date;
  expiresAt: Date;
};

export type APIKeyPerm = {
  workspaceName: string;
  workspaceId: string;
  canRead: boolean;
  canManage: boolean;
  canUpload: boolean;
};
export type CreateApiKeyData = {
  name: string;
  expiresAt: Date;
  canManageAcc: boolean;
  canReadAcc: boolean;
  permissions: APIKeyPerm[];
};
