import { atom, computed } from 'nanostores';
import type { CartItem, ApiResponse } from '../types/api';
import { apiGet, apiPost, apiDelete } from '../lib/api';

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

// Cart State
export const $cart = createPersistentAtom<CartItem[]>('cart', []);
export const $cartLoading = atom<boolean>(false);
export const $cartError = atom<string | null>(null);

// Computed values
export const $cartTotalItems = computed([$cart], (items) =>
    items.reduce((total, item) => total + item.quantity, 0)
);

export const $cartTotalPrice = computed([$cart], (items) =>
    items.reduce((total, item) => {
        const price = item.product?.sale_price || item.product?.reguler_price || 0;
        return total + (price * item.quantity);
    }, 0)
);

// Helper to map Backend Response to Frontend CartItem
interface BackendCartItem {
    id: number; // This is product_id
    product_name: string;
    product_image: string;
    product_status: string; // "active"
    sale_price: number;
    quantity: number;
    unit: string;
    weight: number;
    size: string;
    color: string;
}

// Actions
export async function fetchCarts(): Promise<ApiResponse<BackendCartItem[]>> {
    try {
        $cartLoading.set(true);
        $cartError.set(null);

        const response = await apiGet<BackendCartItem[]>('/carts');

        if (response.status && response.data) {
            const mappedItems: CartItem[] = response.data.map((item) => ({
                product_id: item.id,
                quantity: item.quantity,
                size: item.size || "Standard",
                color: item.color || "Standard",
                product: {
                    id: item.id,
                    product_name: item.product_name,
                    product_image: item.product_image,
                    image: item.product_image,
                    reguler_price: item.sale_price,
                    sale_price: item.sale_price,
                    category_name: "Unknown",
                    stock: 99,
                    weight: item.weight,
                    unit: item.unit,
                }
            }));
            $cart.set(mappedItems);
        }

        return response;
    } catch (error: any) {
        $cartError.set(error.message || 'Failed to fetch cart');
        throw error;
    } finally {
        $cartLoading.set(false);
    }
}

export async function addToCart(
    productId: number,
    quantity: number,
    size: string = "Standard",
    color: string = "Standard",
    sku?: string
): Promise<ApiResponse<any>> {
    try {
        $cartLoading.set(true);
        $cartError.set(null);

        // Optimistic Update (Local)
        const currentCart = $cart.get();
        const existingItemIndex = currentCart.findIndex(
            (i) => i.product_id === productId && i.size === size && i.color === color
        );

        let newCart;
        if (existingItemIndex > -1) {
            newCart = [...currentCart];
            newCart[existingItemIndex].quantity += quantity;
        } else {
            newCart = [...currentCart, {
                product_id: productId,
                quantity,
                size,
                color,
            }];
        }
        $cart.set(newCart);

        // Sync with Backend
        const response = await apiPost<any>('/carts', {
            product_id: productId,
            quantity,
            size,
            color,
            sku: sku || ""
        });

        // Re-fetch to ensure consistency (IDs, etc.)
        await fetchCarts();

        return response;
    } catch (error: any) {
        $cartError.set(error.message || 'Failed to add to cart');
        // Rollback optimistic update on error
        await fetchCarts();
        throw error;
    } finally {
        $cartLoading.set(false);
    }
}

export async function deleteCart(
    productId: number,
    size?: string,
    color?: string
): Promise<ApiResponse<any>> {
    try {
        $cartLoading.set(true);
        $cartError.set(null);

        // Optimistic Update
        const currentCart = $cart.get();
        const newCart = currentCart.filter(
            (i) => !(i.product_id === productId && (!size || i.size === size) && (!color || i.color === color))
        );
        $cart.set(newCart);

        // Sync with Backend
        // WARNING: Backend currently removes ALL items with this ProductID
        // We need to improve backend to support removing specific variant
        const response = await apiDelete<any>(`/carts?product_id=${productId}`);

        // Re-fetch to ensure consistency
        await fetchCarts();

        return response;
    } catch (error: any) {
        $cartError.set(error.message || 'Failed to remove from cart');
        // Rollback optimistic update on error
        await fetchCarts();
        throw error;
    } finally {
        $cartLoading.set(false);
    }
}

export async function deleteAllCart(): Promise<ApiResponse<any>> {
    try {
        $cartLoading.set(true);
        $cartError.set(null);

        // Optimistic Update
        $cart.set([]);

        // Sync with Backend
        const response = await apiDelete<any>('/carts/all');

        return response;
    } catch (error: any) {
        $cartError.set(error.message || 'Failed to clear cart');
        // Rollback optimistic update on error
        await fetchCarts();
        throw error;
    } finally {
        $cartLoading.set(false);
    }
}

export function updateCartItemQuantity(
    productId: number,
    size: string,
    color: string,
    quantity: number
): void {
    const currentCart = $cart.get();
    const newCart = currentCart.map((item) => {
        if (item.product_id === productId && item.size === size && item.color === color) {
            return { ...item, quantity: Math.max(0, quantity) };
        }
        return item;
    });
    $cart.set(newCart);
}

export function clearCart(): void {
    $cart.set([]);
    $cartError.set(null);
}
