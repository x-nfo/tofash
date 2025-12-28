export interface Product {
    id: number;
    product_name: string; // Mapped from backend 'product_name'
    product_image?: string; // Mapped from backend 'product_image' (List)
    image?: string; // Mapped from backend 'image' (Detail)
    reguler_price: number; // Mapped from backend 'reguler_price'
    sale_price: number;
    category_name: string;

    // Optional/Detail fields (might not be in list view)
    description?: string;
    sku?: string;
    size?: string;
    color?: string;
    images?: string[];
    stock?: number;
    category_slug?: string;
    unit?: string;
    weight?: number;
}

export interface CartItem {
    product_id: number;
    quantity: number;
    size: string;
    color: string;
    product?: Product;
}

export interface User {
    user_id: number;
    email: string;
    name: string;
    role_name: string;
    role?: string; // Some responses might use 'role' instead of 'role_name'
    access_token?: string;
}

export interface ApiResponse<T> {
    data: T;
    message: string;
    status: boolean;
}
