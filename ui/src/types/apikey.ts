export type APIKey = {
  id: string;
  name: string;
  workspaceId: string;
  revoked: boolean;
  createdBy: string;
  createdAt: Date;
  expiresAt: Date;
};

export type CreateApiKeyData = {
  name: string;
  expiresAt: Date;
};
