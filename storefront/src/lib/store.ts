import { atom } from 'nanostores';
import type { CartItem, User, Product, ApiResponse } from '../types/api';
import { api, apiPost, apiDelete, apiGet } from './api';

// User Session Store
export const $user = atom<User | null>(null);

// Cart Store
export const $cart = atom<CartItem[]>([]);

// Actions

export function setUser(user: User | null) {
    if (user) {
        $user.set(user);
        // Auto-fetch cart on login
        fetchCart();
    } else {
        $user.set(null);
        $cart.set([]); // Clear cart on logout
    }
}

// Helper to map Backend Response to Frontend CartItem
// Backend returns a flat structure, Frontend needs nested Product
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

export async function fetchCart() {
    const user = $user.get();
    if (!user) return;

    try {
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
                    product_image: item.product_image, // List usually has 'product_image'
                    image: item.product_image, // Detail, map both just in case
                    reguler_price: item.sale_price, // Assuming sale_price is the transactional price
                    sale_price: item.sale_price,
                    category_name: "Unknown", // Not returned by backend cart
                    stock: 99,
                    weight: item.weight,
                    unit: item.unit,
                } as Product
            }));
            $cart.set(mappedItems);
        }
    } catch (error) {
        console.error("Failed to fetch cart:", error);
    }
}

export async function addToCart(item: CartItem) {
    const user = $user.get();

    // Optimistic Update (Local)
    const currentCart = $cart.get();
    const existingItemIndex = currentCart.findIndex(
        (i) => i.product_id === item.product_id && i.size === item.size && i.color === item.color
    );

    let newCart;
    if (existingItemIndex > -1) {
        newCart = [...currentCart];
        newCart[existingItemIndex].quantity += item.quantity;
    } else {
        newCart = [...currentCart, item];
    }
    $cart.set(newCart);

    // Sync with Backend if Logged In
    if (user) {
        try {
            await apiPost('/carts', {
                product_id: item.product_id,
                quantity: item.quantity,
                size: item.size,
                color: item.color,
                sku: item.product?.sku || ""
            });
            // Re-fetch to ensure consistency (IDs, etc.)
            await fetchCart();
        } catch (error) {
            console.error("Failed to add to cart backend:", error);
            // Rollback could be implemented here
        }
    }
}

export async function removeFromCart(productId: number, size: string, color: string) {
    const user = $user.get();

    // Optimistic Update
    const currentCart = $cart.get();
    const newCart = currentCart.filter(
        (i) => !(i.product_id === productId && i.size === size && i.color === color)
    );
    $cart.set(newCart);

    // Sync with Backend
    if (user) {
        try {
            // WARNING: Backend currently removes ALL items with this ProductID
            // We need to improve backend to support removing specific variant
            await apiDelete(`/carts?product_id=${productId}`);
            await fetchCart();
        } catch (error) {
            console.error("Failed to remove from cart backend:", error);
        }
    }
}

export function clearCart() {
    $cart.set([]);
    // Todo: Call backend remove all if needed
}
