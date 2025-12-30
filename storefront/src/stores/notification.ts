import { atom, computed } from 'nanostores';
import type { ApiResponse } from '../types/api';
import { apiGet, apiPut } from '../lib/api';

// Notification Types
export interface Notification {
    id: number;
    user_id: number;
    title: string;
    message: string;
    type: 'info' | 'success' | 'warning' | 'error';
    is_read: boolean;
    created_at: string;
    updated_at: string;
}

// Notification State
export const $notifications = atom<Notification[]>([]);
export const $notificationLoading = atom<boolean>(false);
export const $notificationError = atom<string | null>(null);

// Computed values
export const $unreadCount = computed([$notifications], (notifications) =>
    notifications.filter(n => !n.is_read).length
);

export const $hasUnread = computed([$unreadCount], (count) => count > 0);

// Actions
export async function fetchNotifications(): Promise<ApiResponse<Notification[]>> {
    try {
        $notificationLoading.set(true);
        $notificationError.set(null);

        const response = await apiGet<Notification[]>('/notifications');

        if (response.status && response.data) {
            $notifications.set(response.data);
        }

        return response;
    } catch (error: any) {
        $notificationError.set(error.message || 'Failed to fetch notifications');
        throw error;
    } finally {
        $notificationLoading.set(false);
    }
}

export async function getNotificationDetail(id: number): Promise<ApiResponse<Notification>> {
    try {
        $notificationLoading.set(true);
        $notificationError.set(null);

        const response = await apiGet<Notification>(`/notifications/${id}`);

        if (response.status && response.data) {
            // Update the notification in the list
            const currentNotifications = $notifications.get();
            const updatedNotifications = currentNotifications.map(n =>
                n.id === id ? response.data : n
            );
            $notifications.set(updatedNotifications);
        }

        return response;
    } catch (error: any) {
        $notificationError.set(error.message || 'Failed to fetch notification detail');
        throw error;
    } finally {
        $notificationLoading.set(false);
    }
}

export async function markAsRead(id: number): Promise<ApiResponse<Notification>> {
    try {
        $notificationLoading.set(true);
        $notificationError.set(null);

        const response = await apiPut<Notification>(`/notifications/${id}`, {
            is_read: true
        });

        if (response.status && response.data) {
            // Update the notification in the list
            const currentNotifications = $notifications.get();
            const updatedNotifications = currentNotifications.map(n =>
                n.id === id ? response.data : n
            );
            $notifications.set(updatedNotifications);
        }

        return response;
    } catch (error: any) {
        $notificationError.set(error.message || 'Failed to mark notification as read');
        throw error;
    } finally {
        $notificationLoading.set(false);
    }
}

export async function markAllAsRead(): Promise<void> {
    try {
        $notificationLoading.set(true);
        $notificationError.set(null);

        const currentNotifications = $notifications.get();
        const unreadIds = currentNotifications
            .filter(n => !n.is_read)
            .map(n => n.id);

        // Mark each unread notification as read
        await Promise.all(unreadIds.map(id => markAsRead(id)));

    } catch (error: any) {
        $notificationError.set(error.message || 'Failed to mark all notifications as read');
        throw error;
    } finally {
        $notificationLoading.set(false);
    }
}

export function clearNotifications(): void {
    $notifications.set([]);
}

// Toast notification for temporary messages
export interface Toast {
    id: string;
    message: string;
    type: 'success' | 'error' | 'info' | 'warning';
    duration?: number;
}

export const $toasts = atom<Toast[]>([]);

export function showToast(message: string, type: Toast['type'] = 'info', duration: number = 3000): void {
    const toast: Toast = {
        id: Date.now().toString(),
        message,
        type,
        duration
    };

    const currentToasts = $toasts.get();
    $toasts.set([...currentToasts, toast]);

    // Auto-remove after duration
    setTimeout(() => {
        removeToast(toast.id);
    }, duration);
}

export function removeToast(id: string): void {
    const currentToasts = $toasts.get();
    $toasts.set(currentToasts.filter(t => t.id !== id));
}

export function clearToasts(): void {
    $toasts.set([]);
}
