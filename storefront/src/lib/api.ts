import type { ApiResponse } from '../types/api';
import { $user } from './store';

/**
 * Smart API Fetcher for Astro SSR + Client-side
 * 
 * Features:
 * - Auto-detects SSR vs Client environment
 * - Forwards cookies in SSR mode (for session-based auth)
 * - Supports Bearer token auth (for client-side)
 * - Uses appropriate API URL based on environment
 */

// Detect if running on server (SSR) or client (browser)
const isServer = typeof window === 'undefined';

// Select appropriate API URL based on environment
const BASE_URL = isServer
    ? (import.meta.env.INTERNAL_API_URL || 'http://localhost:8080/api/v1')
    : (import.meta.env.PUBLIC_API_URL || 'http://localhost:8080/api/v1');

type RequestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';

interface RequestOptions {
    method?: RequestMethod;
    body?: any;
    headers?: Record<string, string>;
    auth?: boolean; // Default true - adds Bearer token if available
    astroRequest?: Request; // For SSR: pass Astro.request to forward cookies
    retry?: number; // Number of retry attempts (default: 0)
    retryDelay?: number; // Delay between retries in ms (default: 1000)
    onLoading?: (loading: boolean) => void; // Callback for loading state
}

/**
 * Main API fetch function
 * 
 * Usage in SSR (Astro pages):
 * ```typescript
 * const response = await api<Product[]>('/products', {
 *   astroRequest: Astro.request // Important for cookie forwarding
 * });
 * ```
 * 
 * Usage in Client (browser):
 * ```typescript
 * const response = await api<Product[]>('/products');
 * ```
 */
export async function api<T>(endpoint: string, options: RequestOptions = {}): Promise<ApiResponse<T>> {
    const {
        method = 'GET',
        body,
        headers = {},
        auth = true,
        astroRequest,
        retry = 0,
        retryDelay = 1000,
        onLoading
    } = options;

    // Set loading state
    if (onLoading) {
        onLoading(true);
    }

    const config: RequestInit = {
        method,
        headers: {
            'Content-Type': 'application/json',
            ...headers,
        },
    };

    // Add request body if provided
    if (body) {
        config.body = JSON.stringify(body);
    }

    // SSR Mode: Forward cookies from Astro.request
    if (isServer && astroRequest) {
        const cookies = astroRequest.headers.get('cookie');
        if (cookies) {
            // @ts-ignore - headers is Record<string, string>
            config.headers['Cookie'] = cookies;
        }
    }

    // Client Mode: Include credentials for cookie-based auth
    if (!isServer) {
        config.credentials = 'include'; // Important: sends cookies automatically
    }

    // Add Bearer token if available (for both SSR and Client)
    if (auth) {
        let token = $user.get()?.access_token;

        // Fallback for Client: Try reading from cookie if store is empty (e.g. on refresh/hydrate)
        if (!token && !isServer) {
            const match = document.cookie.match(new RegExp('(^| )access_token=([^;]+)'));
            if (match) token = match[2];
        }

        if (token) {
            // @ts-ignore
            config.headers['Authorization'] = `Bearer ${token}`;
        }
    }

    // Retry logic
    let lastError: Error | null = null;
    for (let attempt = 0; attempt <= retry; attempt++) {
        try {
            const url = `${BASE_URL}${endpoint}`;

            // Log request in development
            if (import.meta.env.DEV) {
                console.log(`[API ${isServer ? 'SSR' : 'Client'}] ${method} ${url} (attempt ${attempt + 1}/${retry + 1})`);
            }

            const response = await fetch(url, config);
            const data = await response.json();

            if (!response.ok) {
                // Don't retry on client errors (4xx)
                if (response.status >= 400 && response.status < 500) {
                    throw new Error(data.message || 'API request failed');
                }
                // Retry on server errors (5xx)
                throw new Error(data.message || 'API request failed');
            }

            // Clear loading state on success
            if (onLoading) {
                onLoading(false);
            }

            return data as ApiResponse<T>;
        } catch (error: any) {
            lastError = error;
            console.error(`API Request Error [${method} ${endpoint}] (attempt ${attempt + 1}/${retry + 1}):`, error);

            // Don't retry on the last attempt
            if (attempt < retry) {
                // Wait before retrying
                await new Promise(resolve => setTimeout(resolve, retryDelay));
            }
        }
    }

    // Clear loading state on error
    if (onLoading) {
        onLoading(false);
    }

    throw lastError || new Error('API request failed');
}

/**
 * Helper function for GET requests
 */
export async function apiGet<T>(endpoint: string, options?: Omit<RequestOptions, 'method'>): Promise<ApiResponse<T>> {
    return api<T>(endpoint, { ...options, method: 'GET' });
}

/**
 * Helper function for POST requests
 */
export async function apiPost<T>(endpoint: string, body: any, options?: Omit<RequestOptions, 'method' | 'body'>): Promise<ApiResponse<T>> {
    return api<T>(endpoint, { ...options, method: 'POST', body });
}

/**
 * Helper function for PUT requests
 */
export async function apiPut<T>(endpoint: string, body: any, options?: Omit<RequestOptions, 'method' | 'body'>): Promise<ApiResponse<T>> {
    return api<T>(endpoint, { ...options, method: 'PUT', body });
}

/**
 * Helper function for DELETE requests
 */
export async function apiDelete<T>(endpoint: string, options?: Omit<RequestOptions, 'method'>): Promise<ApiResponse<T>> {
    return api<T>(endpoint, { ...options, method: 'DELETE' });
}
