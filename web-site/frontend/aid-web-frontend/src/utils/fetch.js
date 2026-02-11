export let ROOT_URL = 'http://172.20.232.212:9900';
export let API_BASE_URL = `api/v1`;
export let API_URL = '';
const authorization = localStorage.getItem('token') ? `Bearer ${localStorage.getItem('token')}` : '';
const DEFAULT_HEADERS = { 'Content-Type': 'application/json', Authorization: authorization };

async function initConfig() {
    try {
        const response = await fetch('/web/config.json');
        const config = await response.json();
        ROOT_URL = config.ROOT_URL;
        API_BASE_URL = config.API_BASE_URL;
        API_URL = `${ROOT_URL}/${API_BASE_URL}`;
    } catch (error) {
        console.error('Failed to load config.json, using default URL', error);
    }
}

export const fetcher = {
    get: async (endpoint, options = {}) => f(endpoint, 'GET', options),
    post: async (endpoint, options = {}) => f(endpoint, 'POST', options),
    put: async (endpoint, options = {}) => f(endpoint, 'PUT', options),
    delete: async (endpoint, options = {}) => f(endpoint, 'DELETE', options),
    patch: async (endpoint, options = {}) => f(endpoint, 'PATCH', options),
};

const f = async (endpoint, method = 'GET', options = {}) => {
    await initConfig(); // 确保在每次请求前都加载最新的配置
    const { headers, ...restOptions } = options;

    try {
        await fetch(`${API_URL}${endpoint}`, { method, ...restOptions, headers: { ...DEFAULT_HEADERS, ...headers } });
    } catch (error) {
        console.error('Fetch error:', error);
        throw error;
    }
};

import Axios from 'axios';

const axiosIns = Axios.create({ baseURL: API_URL, headers: DEFAULT_HEADERS });

export const axios = {
    get: async (endpoint, options = {}) => sendAxios(endpoint, 'GET', options),
    post: async (endpoint, options = {}) => sendAxios(endpoint, 'POST', options),
    put: async (endpoint, options = {}) => sendAxios(endpoint, 'PUT', options),
    delete: async (endpoint, options = {}) => sendAxios(endpoint, 'DELETE', options),
    patch: async (endpoint, options = {}) => sendAxios(endpoint, 'PATCH', options),
};

const sendAxios = async (endpoint, method = 'GET', options = {}) => {
    await initConfig(); // 确保在每次请求前都加载最新的配置
    const { headers, body, params, ...restOptions } = options;

    try {
        return await axiosIns.request({ url: `${API_URL}${endpoint}`, method, headers: { ...DEFAULT_HEADERS, ...headers }, data: body, params, ...restOptions });
    } catch (err) {
        console.error(`Axios error -> ${endpoint}`, err);
        throw err;
    }
};
