import { map, atom } from 'nanostores';
import type { CartItem, User } from '../types/api';

// User Session Store
// We use map for object, or atom if it's a simple value/object replaceable entirely.
// Let's use map for user to update properties easily if needed, or atom for simple session object.
// Atom is safer for "set entire user".
export const $user = atom<User | null>(null);

// Cart Store
export const $cart = atom<CartItem[]>([]);

// Actions

export function setUser(user: User | null) {
    if (user) {
        // Save to localStorage/cookie if needed for persistence, 
        // usually persistence is handled by a separate persistence layer or simpler:
        // localStorage.setItem('user', JSON.stringify(user));
        $user.set(user);
    } else {
        // localStorage.removeItem('user');
        $user.set(null);
    }
}

export function addToCart(item: CartItem) {
    const currentCart = $cart.get();
    const existingItemIndex = currentCart.findIndex(
        (i) => i.product_id === item.product_id && i.size === item.size && i.color === item.color
    );

    if (existingItemIndex > -1) {
        // Update quantity
        const newCart = [...currentCart];
        newCart[existingItemIndex].quantity += item.quantity;
        $cart.set(newCart);
    } else {
        // Add new
        $cart.set([...currentCart, item]);
    }
}

export function removeFromCart(productId: number, size: string, color: string) {
    const currentCart = $cart.get();
    const newCart = currentCart.filter(
        (i) => !(i.product_id === productId && i.size === size && i.color === color)
    );
    $cart.set(newCart);
}

export function clearCart() {
    $cart.set([]);
}
