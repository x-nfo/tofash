/**
 * Toast Notification System
 * Simple utility for showing success/error messages
 */

export type ToastType = 'success' | 'error' | 'info' | 'warning';

interface ToastOptions {
    message: string;
    type?: ToastType;
    duration?: number;
}

const TOAST_CONTAINER_ID = 'toast-container';
const DEFAULT_DURATION = 5000;

/**
 * Create toast container if it doesn't exist
 */
function ensureToastContainer(): HTMLElement {
    let container = document.getElementById(TOAST_CONTAINER_ID);

    if (!container) {
        container = document.createElement('div');
        container.id = TOAST_CONTAINER_ID;
        container.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      z-index: 9999;
      display: flex;
      flex-direction: column;
      gap: 10px;
      max-width: 400px;
    `;
        document.body.appendChild(container);
    }

    return container;
}

/**
 * Get toast styles based on type
 */
function getToastStyles(type: ToastType): string {
    const styles = {
        success: {
            bg: '#10B981',
            icon: '✓'
        },
        error: {
            bg: '#EF4444',
            icon: '✕'
        },
        info: {
            bg: '#3B82F6',
            icon: 'ℹ'
        },
        warning: {
            bg: '#F59E0B',
            icon: '⚠'
        }
    };

    return styles[type].bg;
}

/**
 * Create toast element
 */
function createToastElement(message: string, type: ToastType): HTMLElement {
    const toast = document.createElement('div');
    const bgColor = getToastStyles(type);

    toast.style.cssText = `
    background-color: ${bgColor};
    color: white;
    padding: 16px 20px;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    display: flex;
    align-items: center;
    gap: 12px;
    animation: slideIn 0.3s ease-out;
    cursor: pointer;
    min-height: 60px;
  `;

    toast.innerHTML = `
    <span style="font-weight: bold; font-size: 18px;">${type === 'success' ? '✓' : type === 'error' ? '✕' : 'ℹ'}</span>
    <span style="flex: 1;">${message}</span>
    <span style="opacity: 0.7; font-size: 14px;">✕</span>
  `;

    // Add slide-in animation
    const style = document.createElement('style');
    style.textContent = `
    @keyframes slideIn {
      from {
        transform: translateX(100%);
        opacity: 0;
      }
      to {
        transform: translateX(0);
        opacity: 1;
      }
    }
    @keyframes slideOut {
      from {
        transform: translateX(0);
        opacity: 1;
      }
      to {
        transform: translateX(100%);
        opacity: 0;
      }
    }
  `;
    document.head.appendChild(style);

    // Click to dismiss
    toast.addEventListener('click', () => {
        removeToast(toast);
    });

    return toast;
}

/**
 * Remove toast with animation
 */
function removeToast(toast: HTMLElement): void {
    toast.style.animation = 'slideOut 0.3s ease-out';
    setTimeout(() => {
        toast.remove();
    }, 300);
}

/**
 * Show toast notification
 */
export function showToast(options: ToastOptions): void {
    const { message, type = 'info', duration = DEFAULT_DURATION } = options;

    const container = ensureToastContainer();
    const toast = createToastElement(message, type);

    container.appendChild(toast);

    // Auto-dismiss after duration
    if (duration > 0) {
        setTimeout(() => {
            if (toast.parentElement) {
                removeToast(toast);
            }
        }, duration);
    }
}

/**
 * Show success toast
 */
export function showSuccessToast(message: string, duration?: number): void {
    showToast({ message, type: 'success', duration });
}

/**
 * Show error toast
 */
export function showErrorToast(message: string, duration?: number): void {
    showToast({ message, type: 'error', duration });
}

/**
 * Show info toast
 */
export function showInfoToast(message: string, duration?: number): void {
    showToast({ message, type: 'info', duration });
}

/**
 * Show warning toast
 */
export function showWarningToast(message: string, duration?: number): void {
    showToast({ message, type: 'warning', duration });
}
