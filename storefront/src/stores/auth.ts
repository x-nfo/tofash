import { atom, computed } from 'nanostores';
import type { User, ApiResponse } from '../types/api';
import { api, apiGet, apiPost, apiPut } from '../lib/api';

// Custom persistent atom for complex objects
function createPersistentAtom<T>(key: string, initialValue: T) {
    const store = atom<T>(initialValue);

    // Load from localStorage on client
    if (typeof window !== 'undefined') {
        const saved = localStorage.getItem(key);
        if (saved) {
            try {
                store.set(JSON.parse(saved));
            } catch (e) {
                console.error(`Failed to parse ${key} from localStorage:`, e);
            }
        }

        // Subscribe to changes and save to localStorage
        store.subscribe((value) => {
            localStorage.setItem(key, JSON.stringify(value));
        });
    }

    return store;
}

// Auth State
export const $authUser = createPersistentAtom<User | null>('authUser', null);
export const $authToken = createPersistentAtom<string | null>('authToken', null);
export const $authLoading = atom<boolean>(false);
export const $authError = atom<string | null>(null);

// Computed values
export const $isAuthenticated = computed([$authToken], (token) => !!token);
export const $userRole = computed([$authUser], (user) => user?.role_name || user?.role || null);
export const $isSuperAdmin = computed([$userRole], (role) => role === 'Super Admin');

// Helper to set auth state
function setAuthState(user: User | null, token: string | null) {
    $authUser.set(user);
    $authToken.set(token);

    // Also set cookie for SSR compatibility
    if (typeof document !== 'undefined') {
        if (token) {
            document.cookie = `access_token=${token}; path=/; max-age=604800`; // 7 days
        } else {
            document.cookie = 'access_token=; path=/; max-age=0';
        }
    }
}

// Helper to get token from cookie (for SSR compatibility)
export function getTokenFromCookie(): string | null {
    if (typeof document === 'undefined') return null;
    const match = document.cookie.match(new RegExp('(^| )access_token=([^;]+)'));
    return match ? match[2] : null;
}

// Actions
export async function signin(email: string, password: string): Promise<ApiResponse<any>> {
    try {
        $authLoading.set(true);
        $authError.set(null);

        const response = await apiPost<any>('/auth/signin', { email, password }, { auth: false });

        if (response.status && response.data) {
            const user: User = {
                user_id: response.data.user_id,
                email: response.data.email,
                name: response.data.name,
                role_name: response.data.role_name || response.data.role,
                role: response.data.role,
                access_token: response.data.access_token
            };

            setAuthState(user, response.data.access_token);
        }

        return response;
    } catch (error: any) {
        $authError.set(error.message || 'Sign in failed');
        throw error;
    } finally {
        $authLoading.set(false);
    }
}

export async function signup(userData: {
    email: string;
    password: string;
    name: string;
    password_confirmation?: string;
}): Promise<ApiResponse<any>> {
    try {
        $authLoading.set(true);
        $authError.set(null);

        const response = await apiPost<any>('/auth/signup', userData, { auth: false });

        return response;
    } catch (error: any) {
        $authError.set(error.message || 'Sign up failed');
        throw error;
    } finally {
        $authLoading.set(false);
    }
}

export async function forgotPassword(email: string): Promise<ApiResponse<any>> {
    try {
        $authLoading.set(true);
        $authError.set(null);

        const response = await apiPost<any>('/auth/forgot-password', { email }, { auth: false });

        return response;
    } catch (error: any) {
        $authError.set(error.message || 'Forgot password failed');
        throw error;
    } finally {
        $authLoading.set(false);
    }
}

export async function verifyAccount(token: string): Promise<ApiResponse<any>> {
    try {
        $authLoading.set(true);
        $authError.set(null);

        const response = await apiGet<any>(`/auth/verify-account?token=${token}`, { auth: false });

        return response;
    } catch (error: any) {
        $authError.set(error.message || 'Verify account failed');
        throw error;
    } finally {
        $authLoading.set(false);
    }
}

export async function updatePassword(
    password_new: string,
    password_confirmation: string,
    token: string
): Promise<ApiResponse<any>> {
    try {
        $authLoading.set(true);
        $authError.set(null);

        const response = await apiPut<any>(
            `/auth/update-password?token=${token}`,
            { password_new, password_confirmation },
            { auth: false }
        );

        return response;
    } catch (error: any) {
        $authError.set(error.message || 'Update password failed');
        throw error;
    } finally {
        $authLoading.set(false);
    }
}

export async function getProfile(): Promise<ApiResponse<User>> {
    try {
        $authLoading.set(true);
        $authError.set(null);

        const response = await apiGet<User>('/auth/profile');

        if (response.status && response.data) {
            const currentUser = $authUser.get();
            const updatedUser: User = {
                ...response.data,
                access_token: currentUser?.access_token || getTokenFromCookie() || undefined
            };
            $authUser.set(updatedUser);
        }

        return response;
    } catch (error: any) {
        $authError.set(error.message || 'Get profile failed');
        throw error;
    } finally {
        $authLoading.set(false);
    }
}

export async function updateProfile(userData: Partial<User>): Promise<ApiResponse<User>> {
    try {
        $authLoading.set(true);
        $authError.set(null);

        const response = await apiPut<User>('/auth/profile', userData);

        if (response.status && response.data) {
            const currentUser = $authUser.get();
            const updatedUser: User = {
                ...response.data,
                access_token: currentUser?.access_token || getTokenFromCookie() || undefined
            };
            $authUser.set(updatedUser);
        }

        return response;
    } catch (error: any) {
        $authError.set(error.message || 'Update profile failed');
        throw error;
    } finally {
        $authLoading.set(false);
    }
}

export function logout(): void {
    setAuthState(null, null);
    $authError.set(null);
}

export function checkAuth(): void {
    const token = getTokenFromCookie();
    const user = $authUser.get();

    if (token && !user) {
        // Token exists but no user data - need to fetch profile
        getProfile().catch(() => {
            // If fetching fails, clear invalid token
            logout();
        });
    } else if (!token) {
        // No token, clear user data
        setAuthState(null, null);
    }
}
