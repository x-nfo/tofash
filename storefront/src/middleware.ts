
import { defineMiddleware } from "astro:middleware";

export const onRequest = defineMiddleware(async (context, next) => {
    const { request, cookies, redirect } = context;
    const url = new URL(request.url);

    // Get tokens from cookies
    const accessToken = cookies.get("access_token")?.value;
    const userRole = cookies.get("user_role")?.value;

    // 1. Admin Protection
    if (url.pathname.startsWith("/admin")) {
        // Allow public admin assets or login if we had an admin specific login,
        // but typically /admin/login is the entry.
        // If the path IS /admin/login, allow it.
        if (url.pathname === "/admin/login") {
            return next();
        }

        if (!accessToken || userRole !== "admin") {
            return redirect("/login");
        }
    }

    // 2. User Protection (Checkout, Account)
    if (url.pathname.startsWith("/checkout")) {
        if (!accessToken) {
            return redirect("/login");
        }
    }

    // Pass locals if needed (optional for this task but good practice)
    context.locals.user = accessToken ? { role: userRole } : null;

    return next();
});
