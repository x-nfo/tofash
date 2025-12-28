import type { ApiResponse } from '../types/api';
import { $user } from './store';

// Base URL from environment or default
const BASE_URL = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080/api/v1';

type RequestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';

interface RequestOptions {
    method?: RequestMethod;
    body?: any;
    headers?: Record<string, string>;
    auth?: boolean; // Default true
}

export async function api<T>(endpoint: string, options: RequestOptions = {}): Promise<ApiResponse<T>> {
    const { method = 'GET', body, headers = {}, auth = true } = options;

    const config: RequestInit = {
        method,
        headers: {
            'Content-Type': 'application/json',
            ...headers,
        },
    };

    if (body) {
        config.body = JSON.stringify(body);
    }

    if (auth) {
        const user = $user.get();
        if (user && user.access_token) {
            // @ts-ignore
            config.headers['Authorization'] = `Bearer ${user.access_token}`;
        }
    }

    try {
        const response = await fetch(`${BASE_URL}${endpoint}`, config);
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.message || 'API request failed');
        }

        return data as ApiResponse<T>;
    } catch (error) {
        console.error(`API Request Error [${method} ${endpoint}]:`, error);
        throw error;
    }
}
