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
        const token = localStorage.getItem("token");
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        console.log("Request Intercepted:", config);
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Response Interceptor
axiosInstance.interceptors.response.use(
    (response) => {
        // Handle successful response
        console.log("Response Intercepted:", response);
        return response;
    },
    (error) => {
        // Handle response error (e.g., refresh token or log out)
        if (error.response && error.response.status === 401) {
            console.error("Unauthorized! Redirecting to login...");
            // Redirect to login or handle token refresh
        }
        return Promise.reject(error);
    }
);

export default axiosInstance;
