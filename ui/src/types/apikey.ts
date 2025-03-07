export type APIKey = {
  id: string;
  name: string;
  workspaceId: string;
  revoked: boolean;
  createdBy: string;
  createdAt: Date;
  expiresAt: Date;
};

export type APIWorkspacePerm = {
  id: string;
  read: boolean;
  manage: boolean;
  upload: boolean;
};

export type CreateApiKeyData = {
  name: string;
  expiresAt: Date;
  tenantRead: boolean;
  workspacePerms: APIWorkspacePerm[];
};
