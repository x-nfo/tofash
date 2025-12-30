import { atom, computed } from 'nanostores';
import type { ApiResponse, Pagination } from '../types/api';
import { apiGet, apiPost } from '../lib/api';

// Order Types
export interface Order {
    order_id: string;
    user_id: number;
    total_amount: number;
    status: string;
    payment_status: string;
    shipping_address: string;
    created_at: string;
    updated_at: string;
    items?: OrderItem[];
}

export interface OrderItem {
    id: number;
    order_id: string;
    product_id: number;
    product_name: string;
    quantity: number;
    price: number;
    product_image?: string;
}

export interface CreateOrderData {
    items: Array<{
        product_id: number;
        quantity: number;
        size?: string;
        color?: string;
    }>;
    shipping_address: string;
    payment_method: string;
}

// Order State
export const $orders = atom<Order[]>([]);
export const $order = atom<Order | null>(null);
export const $orderLoading = atom<boolean>(false);
export const $orderError = atom<string | null>(null);

// Pagination State
export const $orderPagination = atom<Pagination>({
    page: 1,
    total_count: 0,
    per_page: 10,
    total_pages: 0
});

// Query params interface
export interface OrderQueryParams {
    page?: number;
    limit?: number;
    status?: string;
    orderBy?: string;
}

// Computed values
export const $hasMoreOrders = computed([$orderPagination], (pag) =>
    pag.page < pag.total_pages
);

export const $totalOrders = computed([$orderPagination], (pag) =>
    pag.total_count
);

// Actions
export async function fetchOrders(params: OrderQueryParams = {}): Promise<ApiResponse<Order[]>> {
    try {
        $orderLoading.set(true);
        $orderError.set(null);

        const { page, limit, status, orderBy } = params;

        const query = new URLSearchParams();
        if (page) query.append('page', page.toString());
        if (limit) query.append('limit', limit.toString());
        if (status) query.append('status', status);
        if (orderBy) query.append('orderBy', orderBy);

        const queryString = query.toString();
        const endpoint = `/orders${queryString ? `?${queryString}` : ''}`;

        const response = await apiGet<Order[]>(endpoint);

        if (response.status) {
            $orders.set(response.data || []);

            // Update pagination if provided
            if (response.pagination) {
                $orderPagination.set(response.pagination);
            }
        }

        return response;
    } catch (error: any) {
        $orderError.set(error.message || 'Failed to fetch orders');
        throw error;
    } finally {
        $orderLoading.set(false);
    }
}

export async function getOrderDetail(orderId: string): Promise<ApiResponse<Order>> {
    try {
        $orderLoading.set(true);
        $orderError.set(null);

        const response = await apiGet<Order>(`/orders/${orderId}`);

        if (response.status && response.data) {
            $order.set(response.data);
        }

        return response;
    } catch (error: any) {
        $orderError.set(error.message || 'Failed to fetch order detail');
        throw error;
    } finally {
        $orderLoading.set(false);
    }
}

export async function createOrder(data: CreateOrderData): Promise<ApiResponse<Order>> {
    try {
        $orderLoading.set(true);
        $orderError.set(null);

        const response = await apiPost<Order>('/orders', data);

        if (response.status && response.data) {
            $order.set(response.data);
        }

        return response;
    } catch (error: any) {
        $orderError.set(error.message || 'Failed to create order');
        throw error;
    } finally {
        $orderLoading.set(false);
    }
}

export function clearOrder(): void {
    $order.set(null);
}

export function clearOrders(): void {
    $orders.set([]);
    $order.set(null);
    $orderPagination.set({
        page: 1,
        total_count: 0,
        per_page: 10,
        total_pages: 0
    });
}
