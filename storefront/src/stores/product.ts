import { atom, computed } from 'nanostores';
import type { Product, ApiResponse, Pagination } from '../types/api';
import { apiGet } from '../lib/api';

// Product State
export const $products = atom<Product[]>([]);
export const $product = atom<Product | null>(null);
export const $productLoading = atom<boolean>(false);
export const $productError = atom<string | null>(null);

// Pagination State
export const $pagination = atom<Pagination>({
    page: 1,
    total_count: 0,
    per_page: 10,
    total_pages: 0
});

// Computed values
export const $hasMoreProducts = computed([$pagination], (pag) =>
    pag.page < pag.total_pages
);

// Query params interface
export interface ProductQueryParams {
    search?: string;
    page?: number;
    limit?: number;
    price?: string;
    orderBy?: string;
    category?: string;
    status?: string;
}

// Actions
export async function fetchProductsHome(): Promise<ApiResponse<Product[]>> {
    try {
        $productLoading.set(true);
        $productError.set(null);

        const response = await apiGet<Product[]>('/products/home');

        if (response.status && response.data) {
            $products.set(response.data);
        }

        return response;
    } catch (error: any) {
        $productError.set(error.message || 'Failed to fetch products');
        throw error;
    } finally {
        $productLoading.set(false);
    }
}

export async function fetchProductsShop(params: ProductQueryParams = {}): Promise<ApiResponse<Product[]>> {
    try {
        $productLoading.set(true);
        $productError.set(null);

        const { search, page, limit, price, orderBy, category } = params;

        const query = new URLSearchParams();
        if (search) query.append('search', search);
        if (page) query.append('page', page.toString());
        if (limit) query.append('limit', limit.toString());
        if (price) query.append('price', price);
        if (orderBy) query.append('orderBy', orderBy);
        if (category) query.append('category', category);

        const queryString = query.toString();
        const endpoint = `/products/shop${queryString ? `?${queryString}` : ''}`;

        const response = await apiGet<Product[]>(endpoint);

        if (response.status) {
            $products.set(response.data || []);

            // Update pagination if provided
            if (response.pagination) {
                $pagination.set(response.pagination);
            }
        }

        return response;
    } catch (error: any) {
        $productError.set(error.message || 'Failed to fetch products');
        throw error;
    } finally {
        $productLoading.set(false);
    }
}

export async function fetchProductDetailHome(productId: number): Promise<ApiResponse<Product>> {
    try {
        $productLoading.set(true);
        $productError.set(null);

        const response = await apiGet<Product>(`/products/home/${productId}`);

        if (response.status && response.data) {
            $product.set(response.data);
        }

        return response;
    } catch (error: any) {
        $productError.set(error.message || 'Failed to fetch product detail');
        throw error;
    } finally {
        $productLoading.set(false);
    }
}

export async function fetchProductsAdmin(params: ProductQueryParams = {}): Promise<ApiResponse<Product[]>> {
    try {
        $productLoading.set(true);
        $productError.set(null);

        const { search, page, limit, status, orderBy, category } = params;

        const query = new URLSearchParams();
        if (search) query.append('search', search);
        if (page) query.append('page', page.toString());
        if (limit) query.append('limit', limit.toString());
        if (status) query.append('status', status);
        if (orderBy) query.append('orderBy', orderBy);
        if (category) query.append('category', category);

        const queryString = query.toString();
        const endpoint = `/admin/products${queryString ? `?${queryString}` : ''}`;

        const response = await apiGet<Product[]>(endpoint);

        if (response.status) {
            $products.set(response.data || []);

            // Update pagination if provided
            if (response.pagination) {
                $pagination.set(response.pagination);
            }
        }

        return response;
    } catch (error: any) {
        $productError.set(error.message || 'Failed to fetch admin products');
        throw error;
    } finally {
        $productLoading.set(false);
    }
}

export async function fetchProductDetail(id: number): Promise<ApiResponse<Product>> {
    try {
        $productLoading.set(true);
        $productError.set(null);

        const response = await apiGet<Product>(`/admin/products/${id}`);

        if (response.status && response.data) {
            $product.set(response.data);
        }

        return response;
    } catch (error: any) {
        $productError.set(error.message || 'Failed to fetch product detail');
        throw error;
    } finally {
        $productLoading.set(false);
    }
}

export function clearProduct(): void {
    $product.set(null);
}

export function clearProducts(): void {
    $products.set([]);
    $product.set(null);
    $pagination.set({
        page: 1,
        total_count: 0,
        per_page: 10,
        total_pages: 0
    });
}
