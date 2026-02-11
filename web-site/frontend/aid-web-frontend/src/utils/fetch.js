// export const ROOT_URL = 'http://172.20.232.212:19900';
export const ROOT_URL = 'http://127.0.0.1:19900';
export const API_BASE_URL = `${ROOT_URL}/api/v1`;
const authorization = localStorage.getItem('token') ? `Bearer ${localStorage.getItem('token')}` : '';
const DEFAULT_HEADERS = { 'Content-Type': 'application/json', Authorization: authorization };

export const fetcher = {
	get: async (endpoint, options = {}) => f(endpoint, 'GET', options),
	post: async (endpoint, options = {}) => f(endpoint, 'POST', options),
	put: async (endpoint, options = {}) => f(endpoint, 'PUT', options),
	delete: async (endpoint, options = {}) => f(endpoint, 'DELETE', options),
	patch: async (endpoint, options = {}) => f(endpoint, 'PATCH', options),
};

const f = async (endpoint, method = 'GET', options = {}) => {
	const { headers, ...restOptions } = options;

	try {
		await fetch(`${API_BASE_URL}${endpoint}`, { method, ...restOptions, headers: { ...DEFAULT_HEADERS, ...headers } });
	} catch (error) {
		console.error('Fetch error:', error);
		throw error;
	}
};

import Axios from 'axios';

const axiosIns = Axios.create({ baseUrl: API_BASE_URL, headers: DEFAULT_HEADERS });

export const axios = {
	get: async (endpoint, options = {}) => sendAxios(endpoint, 'GET', options),
	post: async (endpoint, options = {}) => sendAxios(endpoint, 'POST', options),
	put: async (endpoint, options = {}) => sendAxios(endpoint, 'PUT', options),
	delete: async (endpoint, options = {}) => sendAxios(endpoint, 'DELETE', options),
	patch: async (endpoint, options = {}) => sendAxios(endpoint, 'PATCH', options),
};

const sendAxios = async (endpoint, method = 'GET', options = {}) => {
	const { headers, body, params, ...restOptions } = options;

	try {
		return await axiosIns.request({ url: `${API_BASE_URL}${endpoint}`, method, headers: { ...DEFAULT_HEADERS, ...headers }, data: body, params, ...restOptions });
	} catch (err) {
		console.error(`Axios error -> ${endpoint}`, err);
		throw err;
	}
};
