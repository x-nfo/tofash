/**
 * Authentication Middleware
 * Protects routes that require authentication and handles redirects
 */

import { defineMiddleware } from 'astro:middleware';
import { getTokenFromCookie } from '../stores/auth';

export const onRequest = defineMiddleware((context, next) => {
    const { url, cookies, redirect } = context;

    // Get token from cookie
    const token = cookies.get('access_token')?.value || '';

    // Define protected routes (require authentication)
    const protectedRoutes = ['/admin', '/account'];
    const isProtectedRoute = protectedRoutes.some(route =>
        url.pathname.startsWith(route)
    );

    // Define auth routes (should redirect authenticated users away)
    const authRoutes = ['/auth/signin', '/auth/signup', '/auth/forgot-password', '/auth/verify-account', '/auth/update-password'];
    const isAuthRoute = authRoutes.includes(url.pathname);

    // Redirect to login if trying to access protected route without token
    if (isProtectedRoute && !token) {
        const redirectUrl = encodeURIComponent(url.pathname + url.search);
        return redirect(`/auth/signin?redirect=${redirectUrl}`);
    }

    // Redirect authenticated users away from auth pages
    if (isAuthRoute && token) {
        return redirect('/');
    }

    // Continue to the requested page
    return next();
});
