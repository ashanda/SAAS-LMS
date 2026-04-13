import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost",
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  const tenantDomain = localStorage.getItem("tenant_domain") || "guru.localhost";

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  config.headers["X-Tenant-Domain"] = tenantDomain;
  return config;
});

export default api;