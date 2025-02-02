import axios from "axios";
import { getApiDomain } from "./config";

const axiosInstance = axios.create({
  baseURL: getApiDomain(),
  headers: {
    "Content-Type": "application/json",
  },
});

// Request Interceptor
axiosInstance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("uploadpilottoken");
    if (token) {
      if (!config.headers) {
        config.headers = {};
      }
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Response Interceptor
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      console.error("Unauthorized! Redirecting to login...");
      localStorage.removeItem("uploadpilottoken");
      window.location.href = "/auth";
    }
    return Promise.reject(error);
  },
);

export default axiosInstance;
