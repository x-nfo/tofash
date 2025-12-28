import { defineMiddleware } from "astro:middleware";

export const onRequest = defineMiddleware(async (context, next) => {
    // 1. Get tokens from cookies
    const token = context.cookies.get("access_token")?.value;
    const role = context.cookies.get("user_role")?.value;
    const path = context.url.pathname;

    // 2. PROTECT ADMIN ROUTES
    // If user tries to access /admin* without being an 'admin', redirect to login
    if (path.startsWith("/admin")) {
        if (!token || role !== "admin") {
            return context.redirect("/login");
        }
    }

    // 3. PROTECT CHECKOUT (Optional but recommended)
    // If user tries to checkout without login, redirect to login
    if (path.startsWith("/checkout")) {
        if (!token) {
            return context.redirect("/login");
        }
    }

    return next();
});
