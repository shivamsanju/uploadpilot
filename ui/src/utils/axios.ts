import axios from 'axios';
import { getApiDomain, getTenantId } from './config';

const axiosBaseInstance = axios.create({
  baseURL: getApiDomain(),
  withCredentials: true,
  withXSRFToken: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor
axiosBaseInstance.interceptors.request.use(
  config => {
    return config;
  },
  error => {
    return Promise.reject(error);
  },
);

// Response Interceptor
axiosBaseInstance.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    if (error.response && error.response.status === 401) {
      console.error('Unauthorized! Redirecting to login...');
    }
    return Promise.reject(error);
  },
);

// Tenant Instance
const axiosTenantInstance = axios.create({
  baseURL: getApiDomain(),
  withCredentials: true,
  withXSRFToken: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor
axiosTenantInstance.interceptors.request.use(
  config => {
    const tenantId = getTenantId();
    if (!tenantId) {
      return Promise.reject(new Error('tenantId is required'));
    }
    config.baseURL = `${config.baseURL}/tenants/${tenantId}`;
    return config;
  },
  error => {
    return Promise.reject(error);
  },
);

// Response Interceptor
axiosTenantInstance.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    if (error.response && error.response.status === 401) {
      console.error('Unauthorized! Redirecting to login...');
    }
    return Promise.reject(error);
  },
);

export { axiosBaseInstance, axiosTenantInstance };
