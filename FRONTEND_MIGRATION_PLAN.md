# Frontend Migration Plan: Astro Storefront

## Overview

This document outlines the plan to migrate and enhance the Astro Storefront project by implementing features from the reference Nuxt.js microservices project (`refrensi_fe/micro-sayur-nuxt-main`).

## Current State Analysis

### Reference Project (micro-sayur-nuxt)
- **Framework**: Nuxt.js 3 with Vue 3
- **State Management**: Pinia
- **UI Libraries**: Bootstrap 5, Tailwind CSS, Swiper, Dropzone, Quill Editor
- **Architecture**: Microservices (separate API endpoints per service)
- **Features**: Complete e-commerce functionality

### Current Storefront (Astro)
- **Framework**: Astro 5
- **State Management**: NanoStores
- **UI Libraries**: Tailwind CSS
- **Architecture**: Monolithic (single API base URL)
- **Features**: Basic cart, basic admin page

## Comparison Matrix

| Feature | Reference (Nuxt) | Current (Astro) | Gap |
|---------|-------------------|-------------------|-----|
| Authentication | ✅ Complete | ❌ None | 100% |
| Cart Management | ✅ Complete | ⚠️ Basic | 60% |
| Product Catalog | ✅ Complete | ❌ None | 100% |
| Order Management | ✅ Complete | ❌ None | 100% |
| Payment Integration | ✅ Midtrans | ❌ None | 100% |
| Admin Dashboard | ✅ Complete | ⚠️ Basic | 70% |
| Customer Dashboard | ✅ Complete | ❌ None | 100% |
| Notifications | ✅ Complete | ❌ None | 100% |
| File Upload | ✅ Dropzone | ❌ None | 100% |
| Rich Text Editor | ✅ Quill | ❌ None | 100% |
| Image Gallery | ✅ Swiper | ❌ None | 100% |

## Migration Strategy

### Phase 1: Foundation Setup (Week 1-2)

#### 1.1 State Management Upgrade
**Current**: NanoStores (basic)
**Target**: Enhanced NanoStores with persistence

**Tasks**:
- [ ] Install `@nanostores/persistent` (already in package.json)
- [ ] Create store structure similar to Pinia
- [ ] Implement auth store with token management
- [ ] Implement cart store with persistence
- [ ] Implement product store
- [ ] Implement order store
- [ ] Implement notification store

**Files to Create**:
```
storefront/src/stores/
├── auth.ts
├── cart.ts
├── product.ts
├── order.ts
├── notification.ts
└── index.ts
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/stores/`

#### 1.2 API Client Setup
**Current**: No centralized API client
**Target**: Centralized API client with interceptors

**Tasks**:
- [ ] Create API client utility
- [ ] Implement request interceptor (add auth token)
- [ ] Implement response interceptor (handle errors)
- [ ] Add retry logic
- [ ] Add loading states

**Files to Create**:
```
storefront/src/lib/
├── api-client.ts
├── api-endpoints.ts
└── types.ts
```

#### 1.3 Configuration Management
**Current**: No environment-based config
**Target**: Environment-based API URLs

**Tasks**:
- [ ] Create `.env.example` file
- [ ] Update `astro.config.mjs` with runtime config
- [ ] Add environment variables for API URLs
- [ ] Add Midtrans client key configuration

**Files to Modify**:
- `storefront/astro.config.mjs`
- `storefront/.env.example`

### Phase 2: Authentication Module (Week 2-3)

#### 2.1 Authentication Pages
**Tasks**:
- [ ] Create sign-in page (`/auth/signin`)
- [ ] Create sign-up page (`/auth/signup`)
- [ ] Create forgot-password page (`/auth/forgot-password`)
- [ ] Create verify-account page (`/auth/verify-account`)
- [ ] Create update-password page (`/auth/update-password`)

**Files to Create**:
```
storefront/src/pages/auth/
├── signin.astro
├── signup.astro
├── forgot-password.astro
├── verify-account.astro
└── update-password.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/pages/auth/`

#### 2.2 Auth Store Implementation
**Tasks**:
- [ ] Implement `signin` action
- [ ] Implement `signup` action
- [ ] Implement `logout` action
- [ ] Implement `forgotPassword` action
- [ ] Implement `verifyAccount` action
- [ ] Implement `updatePassword` action
- [ ] Implement `getProfile` action
- [ ] Implement `updateProfile` action
- [ ] Add token persistence (cookies/localStorage)
- [ ] Add user data persistence

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/stores/auth.js`

#### 2.3 Auth Components
**Tasks**:
- [ ] Create LoginModal component
- [ ] Create SignupForm component
- [ ] Create ForgotPasswordForm component
- [ ] Create UpdatePasswordForm component
- [ ] Create ProfileDropdown component

**Files to Create**:
```
storefront/src/components/auth/
├── LoginModal.astro
├── SignupForm.astro
├── ForgotPasswordForm.astro
└── UpdatePasswordForm.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/components/modals/`

#### 2.4 Middleware
**Tasks**:
- [ ] Create auth middleware
- [ ] Protect authenticated routes
- [ ] Redirect unauthenticated users
- [ ] Check token validity

**Files to Create**:
```
storefront/src/middleware/
└── auth.ts
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/middleware/auth.js`

### Phase 3: Product Module (Week 3-4)

#### 3.1 Product Pages
**Tasks**:
- [ ] Create product listing page (`/products`)
- [ ] Create product detail page (`/products/[id]`)
- [ ] Create category filtering
- [ ] Implement search functionality
- [ ] Implement pagination

**Files to Create**:
```
storefront/src/pages/products/
├── index.astro
└── [id].astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/pages/shop/`

#### 3.2 Product Components
**Tasks**:
- [ ] Create ProductCard component
- [ ] Create ProductGallery component
- [ ] Create ProductFilter component
- [ ] Create ProductSearch component
- [ ] Create CategoryMenu component

**Files to Create**:
```
storefront/src/components/product/
├── ProductCard.astro
├── ProductGallery.astro
├── ProductFilter.astro
├── ProductSearch.astro
└── CategoryMenu.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/components/home/`

#### 3.3 Product Store
**Tasks**:
- [ ] Implement `fetchProducts` action
- [ ] Implement `fetchProductDetail` action
- [ ] Implement `fetchCategories` action
- [ ] Implement search functionality
- [ ] Implement filter functionality

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/stores/product.js`

#### 3.4 UI Libraries Integration
**Tasks**:
- [ ] Install and configure Swiper for image gallery
- [ ] Install and configure rating component
- [ ] Add product image lazy loading

**Dependencies to Add**:
```json
{
  "swiper": "^11.2.6",
  "@types/swiper": "^11.0.0"
}
```

### Phase 4: Cart Module (Week 4-5)

#### 4.1 Cart Pages
**Tasks**:
- [ ] Enhance existing cart page (`/cart`)
- [ ] Add cart item quantity controls
- [ ] Add cart item removal
- [ ] Add cart total calculation
- [ ] Add proceed to checkout button

**Files to Modify**:
- `storefront/src/pages/cart.astro`

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/pages/shop/cart.vue`

#### 4.2 Cart Components
**Tasks**:
- [ ] Create CartSidebar component
- [ ] Create CartItem component
- [ ] Create CartSummary component
- [ ] Create AddToCartButton component

**Files to Create**:
```
storefront/src/components/cart/
├── CartSidebar.astro
├── CartItem.astro
├── CartSummary.astro
└── AddToCartButton.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/components/modals/CartSidebar.vue`

#### 4.3 Cart Store Enhancement
**Tasks**:
- [ ] Implement `addToCart` action
- [ ] Implement `deleteCart` action
- [ ] Implement `deleteAllCart` action
- [ ] Implement `fetchCarts` action
- [ ] Add cart persistence (localStorage/cookies)
- [ ] Add cart total calculation

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/stores/cart.js`

### Phase 5: Order Module (Week 5-6)

#### 5.1 Order Pages
**Tasks**:
- [ ] Create checkout page (`/checkout`)
- [ ] Create order success page (`/checkout/success`)
- [ ] Create order list page (`/orders`)
- [ ] Create order detail page (`/orders/[id]`)

**Files to Create**:
```
storefront/src/pages/
├── checkout.astro
├── checkout/
│   └── success.astro
└── orders/
    ├── index.astro
    └── [id].astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/pages/shop/`

#### 5.2 Order Components
**Tasks**:
- [ ] Create CheckoutForm component
- [ ] Create OrderSummary component
- [ ] Create OrderItem component
- [ ] Create OrderStatus component
- [ ] Create ShippingForm component

**Files to Create**:
```
storefront/src/components/order/
├── CheckoutForm.astro
├── OrderSummary.astro
├── OrderItem.astro
├── OrderStatus.astro
└── ShippingForm.astro
```

#### 5.3 Order Store
**Tasks**:
- [ ] Implement `fetchOrders` action
- [ ] Implement `createOrder` action
- [ ] Implement `getDetailOrders` action
- [ ] Implement pagination
- [ ] Implement filtering by status

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/stores/orders.js`

### Phase 6: Payment Module (Week 6-7)

#### 6.1 Payment Integration
**Tasks**:
- [ ] Integrate Midtrans Snap.js
- [ ] Create payment page
- [ ] Implement payment webhook handler
- [ ] Add payment status tracking

**Files to Create**:
```
storefront/src/pages/
└── payment.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/nuxt.config.ts` (line 67)

#### 6.2 Payment Components
**Tasks**:
- [ ] Create PaymentMethod component
- [ ] Create PaymentCard component
- [ ] Create PaymentStatus component

**Files to Create**:
```
storefront/src/components/payment/
├── PaymentMethod.astro
├── PaymentCard.astro
└── PaymentStatus.astro
```

#### 6.3 Payment Store
**Tasks**:
- [ ] Implement `createPayment` action
- [ ] Implement `getPaymentStatus` action
- [ ] Implement `getPaymentDetail` action

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/stores/payment.js`

### Phase 7: Admin Dashboard (Week 7-8)

#### 7.1 Admin Pages
**Tasks**:
- [ ] Enhance admin dashboard (`/admin`)
- [ ] Create admin products page (`/admin/products`)
- [ ] Create admin categories page (`/admin/categories`)
- [ ] Create admin orders page (`/admin/orders`)
- [ ] Create admin customers page (`/admin/customers`)
- [ ] Create admin roles page (`/admin/roles`)

**Files to Create**:
```
storefront/src/pages/admin/
├── index.astro (enhance existing)
├── products/
│   ├── index.astro
│   ├── create.astro
│   └── edit/
│       └── [id].astro
├── categories/
│   ├── index.astro
│   ├── create.astro
│   └── edit/
│       └── [id].astro
├── orders/
│   ├── index.astro
│   └── show/
│       └── [id].astro
├── customers/
│   ├── index.astro
│   ├── create.astro
│   └── edit/
│       └── [id].astro
└── roles/
    ├── index.astro
    └── create.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/pages/dashboard/`

#### 7.2 Admin Components
**Tasks**:
- [ ] Create SidebarNav component
- [ ] Create Navbar component
- [ ] Create NotificationDropdown component
- [ ] Create ProfileDropdown component
- [ ] Create DataTable component
- [ ] Create Pagination component

**Files to Create**:
```
storefront/src/components/admin/
├── SidebarNav.astro
├── Navbar.astro
├── NotificationDropdown.astro
├── NotificationItem.astro
├── ProfileDropdown.astro
├── DataTable.astro
└── Pagination.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/components/admin/`

#### 7.3 Admin Layouts
**Tasks**:
- [ ] Create admin layout
- [ ] Create dashboard layout

**Files to Create**:
```
storefront/src/layouts/
├── admin.astro
└── dashboard.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/layouts/`

#### 7.4 Admin Middleware
**Tasks**:
- [ ] Create admin middleware
- [ ] Protect admin routes
- [ ] Check admin role
- [ ] Redirect unauthorized users

**Files to Create**:
```
storefront/src/middleware/
└── admin.ts
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/middleware/admin.js`

### Phase 8: Customer Dashboard (Week 8-9)

#### 8.1 Customer Pages
**Tasks**:
- [ ] Create customer dashboard (`/account`)
- [ ] Create customer profile page (`/account/setting`)
- [ ] Create customer orders page (`/account/orders`)

**Files to Create**:
```
storefront/src/pages/account/
├── index.astro
├── setting.astro
└── orders.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/pages/account/`

#### 8.2 Customer Components
**Tasks**:
- [ ] Create ProfileForm component
- [ ] Create OrderHistory component
- [ ] Create OrderTracking component

**Files to Create**:
```
storefront/src/components/account/
├── ProfileForm.astro
├── OrderHistory.astro
└── OrderTracking.astro
```

### Phase 9: Advanced Features (Week 9-10)

#### 9.1 File Upload
**Tasks**:
- [ ] Install Dropzone or similar library
- [ ] Create ImageUpload component
- [ ] Implement drag-and-drop upload
- [ ] Add image preview
- [ ] Add progress indicator

**Dependencies to Add**:
```json
{
  "dropzone": "^6.0.0-beta.2"
}
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/package.json`

#### 9.2 Rich Text Editor
**Tasks**:
- [ ] Install Quill or similar editor
- [ ] Create RichTextEditor component
- [ ] Implement toolbar
- [ ] Add image upload integration

**Dependencies to Add**:
```json
{
  "quill": "^2.0.3"
}
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/package.json`

#### 9.3 Notifications
**Tasks**:
- [ ] Create notification store
- [ ] Implement WebSocket connection
- [ ] Create NotificationBell component
- [ ] Create NotificationList component
- [ ] Add notification settings

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/stores/notification.js`

#### 9.4 Common Components
**Tasks**:
- [ ] Create Header component
- [ ] Create Footer component
- [ ] Create LoadingSpinner component
- [ ] Create ErrorAlert component
- [ ] Create SuccessAlert component
- [ ] Create Modal component
- [ ] Create Toast component

**Files to Create**:
```
storefront/src/components/common/
├── Header.astro
├── Footer.astro
├── LoadingSpinner.astro
├── ErrorAlert.astro
├── SuccessAlert.astro
├── Modal.astro
└── Toast.astro
```

**Reference**: `refrensi_fe/micro-sayur-nuxt-main/components/common/`

### Phase 10: Testing & Optimization (Week 10-11)

#### 10.1 Testing
**Tasks**:
- [ ] Unit tests for stores
- [ ] Component tests
- [ ] Integration tests
- [ ] E2E tests with Playwright

#### 10.2 Performance Optimization
**Tasks**:
- [ ] Implement code splitting
- [ ] Add lazy loading for images
- [ ] Optimize bundle size
- [ ] Add caching strategy
- [ ] Implement service worker

#### 10.3 SEO & Analytics
**Tasks**:
- [ ] Add meta tags
- [ ] Add Open Graph tags
- [ ] Add structured data
- [ ] Integrate analytics (Google Analytics, etc.)

## API Integration

### Current Backend API Structure

Based on the Go backend analysis, the API endpoints are:

```
/api/v1/
├── auth/
│   ├── POST /register
│   ├── POST /login
│   ├── POST /forgot-password
│   ├── GET /verify-account
│   ├── GET /profile
│   └── PUT /profile
├── products/
│   ├── GET /products
│   ├── GET /products/home
│   ├── GET /products/:id
│   ├── GET /categories
│   └── POST /carts (cart operations)
├── orders/
│   ├── POST /orders
│   ├── GET /orders
│   └── GET /orders/:orderID
├── payments/
│   ├── POST /payments
│   ├── GET /payments
│   ├── GET /payments/:id
│   └── POST /midtrans/webhook
└── notifications/
    ├── GET /notifications
    ├── GET /notifications/:id
    └── PUT /notifications/:id
```

### API Client Configuration

**Environment Variables**:
```env
# API Configuration
API_BASE_URL=http://localhost:8081/api/v1

# Midtrans Configuration
MIDTRANS_CLIENT_KEY=SB-Mid-client-XXXX
MIDTRANS_ENVIRONMENT=sandbox
```

## Technical Decisions

### 1. State Management
**Decision**: Keep NanoStores (Astro-native) instead of migrating to Pinia

**Rationale**:
- NanoStores is lightweight and Astro-native
- No need for additional dependencies
- Better performance for Astro's island architecture
- Can implement similar patterns to Pinia

### 2. Component Architecture
**Decision**: Use Astro Islands for interactive components

**Rationale**:
- Better performance (partial hydration)
- Smaller bundle size
- SEO-friendly
- Progressive enhancement

### 3. Styling
**Decision**: Keep Tailwind CSS, add Bootstrap-like utilities

**Rationale**:
- Tailwind is already configured
- Can create utility classes similar to Bootstrap
- Smaller bundle size than full Bootstrap
- More flexible customization

### 4. Routing
**Decision**: Use Astro's file-based routing

**Rationale**:
- Native to Astro
- Automatic route generation
- Better performance
- Easier maintenance

### 5. API Architecture
**Decision**: Single API base URL (monolithic backend)

**Rationale**:
- Backend is monolithic (not microservices)
- Simpler configuration
- Fewer network requests
- Better error handling

## Dependencies to Add

```json
{
  "dependencies": {
    // Existing
    "astro": "^5.16.6",
    "@astrojs/node": "^9.5.1",
    "@nanostores/persistent": "^1.2.0",
    "nanostores": "^1.1.0",
    "tailwindcss": "^4.1.18",

    // New additions
    "swiper": "^11.2.6",
    "quill": "^2.0.3",
    "dropzone": "^6.0.0-beta.2",
    "dayjs": "^1.11.13",
    "@heroicons/vue": "^2.2.0",
    "@iconify/vue": "^4.3.0"
  },
  "devDependencies": {
    "@types/swiper": "^11.0.0",
    "@types/quill": "^2.0.0",
    "@types/dropzone": "^6.0.0"
  }
}
```

## File Structure

### Final Structure
```
storefront/
├── src/
│   ├── components/
│   │   ├── admin/
│   │   │   ├── SidebarNav.astro
│   │   │   ├── Navbar.astro
│   │   │   ├── NotificationDropdown.astro
│   │   │   ├── ProfileDropdown.astro
│   │   │   ├── DataTable.astro
│   │   │   └── Pagination.astro
│   │   ├── auth/
│   │   │   ├── LoginModal.astro
│   │   │   ├── SignupForm.astro
│   │   │   ├── ForgotPasswordForm.astro
│   │   │   └── UpdatePasswordForm.astro
│   │   ├── cart/
│   │   │   ├── CartSidebar.astro
│   │   │   ├── CartItem.astro
│   │   │   ├── CartSummary.astro
│   │   │   └── AddToCartButton.astro
│   │   ├── common/
│   │   │   ├── Header.astro
│   │   │   ├── Footer.astro
│   │   │   ├── LoadingSpinner.astro
│   │   │   ├── ErrorAlert.astro
│   │   │   ├── SuccessAlert.astro
│   │   │   ├── Modal.astro
│   │   │   └── Toast.astro
│   │   ├── order/
│   │   │   ├── CheckoutForm.astro
│   │   │   ├── OrderSummary.astro
│   │   │   ├── OrderItem.astro
│   │   │   ├── OrderStatus.astro
│   │   │   └── ShippingForm.astro
│   │   ├── payment/
│   │   │   ├── PaymentMethod.astro
│   │   │   ├── PaymentCard.astro
│   │   │   └── PaymentStatus.astro
│   │   ├── product/
│   │   │   ├── ProductCard.astro
│   │   │   ├── ProductGallery.astro
│   │   │   ├── ProductFilter.astro
│   │   │   ├── ProductSearch.astro
│   │   │   └── CategoryMenu.astro
│   │   └── account/
│   │       ├── ProfileForm.astro
│   │       ├── OrderHistory.astro
│   │       └── OrderTracking.astro
│   ├── layouts/
│   │   ├── admin.astro
│   │   ├── dashboard.astro
│   │   └── default.astro
│   ├── lib/
│   │   ├── api-client.ts
│   │   ├── api-endpoints.ts
│   │   └── types.ts
│   ├── middleware/
│   │   ├── auth.ts
│   │   └── admin.ts
│   ├── pages/
│   │   ├── admin/
│   │   │   ├── index.astro
│   │   │   ├── products/
│   │   │   ├── categories/
│   │   │   ├── orders/
│   │   │   ├── customers/
│   │   │   └── roles/
│   │   ├── account/
│   │   │   ├── index.astro
│   │   │   ├── setting.astro
│   │   │   └── orders.astro
│   │   ├── auth/
│   │   │   ├── signin.astro
│   │   │   ├── signup.astro
│   │   │   ├── forgot-password.astro
│   │   │   ├── verify-account.astro
│   │   │   └── update-password.astro
│   │   ├── cart.astro (enhance existing)
│   │   ├── checkout.astro
│   │   ├── checkout/
│   │   │   └── success.astro
│   │   ├── orders/
│   │   │   ├── index.astro
│   │   │   └── [id].astro
│   │   ├── products/
│   │   │   ├── index.astro
│   │   │   └── [id].astro
│   │   └── index.astro
│   └── stores/
│       ├── auth.ts
│       ├── cart.ts
│       ├── product.ts
│       ├── order.ts
│       ├── notification.ts
│       └── index.ts
├── public/
│   ├── images/
│   ├── css/
│   └── js/
├── astro.config.mjs
├── package.json
└── .env.example
```

## Implementation Priority

### Critical Path (Must Have)
1. ✅ Authentication (Phase 2)
2. ✅ Product Catalog (Phase 3)
3. ✅ Cart Management (Phase 4)
4. ✅ Order Management (Phase 5)
5. ✅ Payment Integration (Phase 6)

### Important Path (Should Have)
6. ⚠️ Admin Dashboard (Phase 7)
7. ⚠️ Customer Dashboard (Phase 8)
8. ⚠️ Notifications (Phase 9.3)

### Nice to Have (Could Have)
9. 💡 File Upload (Phase 9.1)
10. 💡 Rich Text Editor (Phase 9.2)
11. 💡 Advanced UI Components (Phase 9.4)

## Risk Assessment

### Technical Risks

| Risk | Impact | Probability | Mitigation |
|-------|---------|-------------|------------|
| NanoStores limitations | High | Medium | Test thoroughly, consider alternative if needed |
| Astro Islands complexity | Medium | High | Start simple, increment complexity |
| API integration issues | High | Medium | Comprehensive error handling |
| Performance issues | Medium | Medium | Code splitting, lazy loading |
| Browser compatibility | Low | Low | Test on multiple browsers |

### Timeline Risks

| Risk | Impact | Probability | Mitigation |
|-------|---------|-------------|------------|
| Scope creep | High | High | Strict phase boundaries |
| Delayed backend API | High | Medium | Mock API for development |
| Resource constraints | Medium | Low | Prioritize critical path |

## Success Criteria

### Phase 1: Foundation
- [ ] All stores created and tested
- [ ] API client working
- [ ] Configuration management in place

### Phase 2: Authentication
- [ ] User can sign in
- [ ] User can sign up
- [ ] User can reset password
- [ ] User can verify account
- [ ] User can update profile
- [ ] Auth middleware protecting routes

### Phase 3: Products
- [ ] Product listing works
- [ ] Product detail works
- [ ] Categories filter works
- [ ] Search works
- [ ] Pagination works

### Phase 4: Cart
- [ ] Add to cart works
- [ ] Remove from cart works
- [ ] Update quantity works
- [ ] Cart total calculates correctly
- [ ] Cart persists across sessions

### Phase 5: Orders
- [ ] Create order works
- [ ] Order list displays
- [ ] Order detail shows
- [ ] Order status updates
- [ ] Pagination works

### Phase 6: Payment
- [ ] Midtrans integration works
- [ ] Payment page loads
- [ ] Payment processes
- [ ] Webhook handles callbacks
- [ ] Payment status updates

### Phase 7: Admin
- [ ] Admin can manage products
- [ ] Admin can manage categories
- [ ] Admin can manage orders
- [ ] Admin can manage customers
- [ ] Admin can manage roles

### Phase 8: Customer Dashboard
- [ ] Customer can view profile
- [ ] Customer can update profile
- [ ] Customer can view order history
- [ ] Customer can track orders

## Next Steps

1. **Review and Approve**: Review this plan with stakeholders
2. **Setup Development Environment**: Ensure all dependencies are installed
3. **Start Phase 1**: Begin with foundation setup
4. **Weekly Reviews**: Review progress weekly and adjust plan as needed
5. **Testing**: Test each phase before moving to next

## Resources

### Reference Documentation
- Astro Docs: https://docs.astro.build
- NanoStores Docs: https://github.com/nanostores/nanostores
- Tailwind Docs: https://tailwindcss.com/docs
- Swiper Docs: https://swiperjs.com/
- Quill Docs: https://quilljs.com/docs/

### Backend API Documentation
- Current API endpoints: Analyzed from Go backend code
- API base URL: `http://localhost:8081/api/v1`

### Reference Project
- Location: `refrensi_fe/micro-sayur-nuxt-main/`
- Framework: Nuxt.js 3 + Vue 3
- State: Pinia
- UI: Bootstrap 5 + Tailwind CSS

---

**Document Version**: 1.0
**Created**: 2025-12-29
**Status**: Ready for Implementation
**Estimated Timeline**: 10-12 weeks
