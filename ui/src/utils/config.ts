import { TENANT_ID_KEY } from '../constants/tenancy';

export const getApiDomain = () => {
  const apiUrl = process.env.REACT_APP_BACKEND_URL || `http://localhost:8080`;
  return apiUrl;
};

export const getWebsiteDomain = () => {
  const websiteUrl =
    process.env.REACT_APP_WEBSITE_URL || `http://localhost:3000`;
  return websiteUrl;
};

export const getAppName = () => {
  return process.env.REACT_APP_APP_NAME || 'UploadPilot';
};

export const getUploadApiDomain = () => {
  const uploadApiUrl =
    process.env.REACT_APP_BACKEND_URL || `http://localhost:8080`;
  return uploadApiUrl;
};

export const getTenantId = (): string | null => {
  const tenantId = localStorage.getItem(TENANT_ID_KEY);
  return tenantId;
};
