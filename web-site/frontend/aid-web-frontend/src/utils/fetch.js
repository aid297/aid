// export const ROOT_URL = "http://172.20.232.212:9900";
export const ROOT_URL = "http://127.0.0.1:9900";
export const API_BASE_URL = `${ROOT_URL}/api/v1`; // 设置根 URL
const DEFAULT_HEADERS = {
    "Content-Type": "application/json",
    Authorization: localStorage.getItem("token") ? `Bearer ${localStorage.getItem("token")}` : "", // 动态添加 Token
};

export const fetcher = {
    get: async (endpoint, options = {}) => f(endpoint, "GET", options),
    post: async (endpoint, options = {}) => f(endpoint, "POST", options),
    put: async (endpoint, options = {}) => f(endpoint, "PUT", options),
    delete: async (endpoint, options = {}) => f(endpoint, "DELETE", options),
    patch: async (endpoint, options = {}) => f(endpoint, "PATCH", options),
};

const f = async (endpoint, method = "GET", options = {}) => {
    console.log(`Fetch -> ${method} ${API_BASE_URL}${endpoint}`, options);
    const { headers, ...restOptions } = options;

    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            method,
            ...restOptions,
            headers: {
                ...DEFAULT_HEADERS, // 默认 headers
                ...headers, // 合并传入的 headers
            },
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return await response.json();
    } catch (error) {
        console.error("Fetch error:", error);
        throw error;
    }
};

import Axios from "axios";

const axiosIns = Axios.create({
    baseUrl: API_BASE_URL,
    headers: DEFAULT_HEADERS,
});

export const axios = {
    get: async (endpoint, options = {}) => a(endpoint, "GET", options),
    post: async (endpoint, options = {}) => a(endpoint, "POST", options),
    put: async (endpoint, options = {}) => a(endpoint, "PUT", options),
    delete: async (endpoint, options = {}) => a(endpoint, "DELETE", options),
    patch: async (endpoint, options = {}) => a(endpoint, "PATCH", options),
};

const a = async (endpoint, method = "GET", options = {}) => {
    const { headers, data, params, ...restOptions } = options;

    try {
        const res = await axiosIns.request({
            url: `${API_BASE_URL}${endpoint}`,
            method,
            headers: {
                ...DEFAULT_HEADERS, // 默认 headers
                ...headers, // 合并传入的 headers
            },
            data, // POST/PUT 请求体
            params, // GET 请求参数
            ...restOptions, // 其他选项
        });

        console.log(`Axios success -> ${endpoint}`, res);
        return res.data; // 返回响应数据
    } catch (err) {
        console.error(`Axios error -> ${endpoint}`, err);
        throw err;
    }
};
