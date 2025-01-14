export function getApiDomain() {
    const apiUrl = process.env.REACT_APP_BACKEND_URL || `http://localhost:8080`;
    return apiUrl;
}

export function getWebsiteDomain() {
    const websiteUrl = process.env.REACT_APP_WEBSITE_URL || `http://localhost:3000`;
    return websiteUrl;
}

export function getAppName() {
    return process.env.REACT_APP_APP_NAME || "Upload Pilot";
}

