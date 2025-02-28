import axios from 'axios';
import { TENANT_ID_HEADER, TENANT_ID_KEY } from '../constants/tenancy';
import { getApiDomain } from './config';

const axiosInstance = axios.create({
  baseURL: getApiDomain(),
  withCredentials: true,
  withXSRFToken: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor
axiosInstance.interceptors.request.use(
  config => {
    const tenantId = localStorage.getItem(TENANT_ID_KEY);
    if (tenantId) {
      if (!config.headers) {
        config.headers = {};
      }
      config.headers[TENANT_ID_HEADER] = tenantId;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  },
);

// Response Interceptor
axiosInstance.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    if (error.response && error.response.status === 401) {
      console.error('Unauthorized! Redirecting to login...');
      // localStorage.removeItem('uploadpilottoken');
      // window.location.href = '/auth';
    }
    return Promise.reject(error);
  },
);

export default axiosInstance;
