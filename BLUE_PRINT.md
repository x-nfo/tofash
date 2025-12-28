Act as a Senior Go Backend Engineer. I am migrating a Microservices architecture into a Modular Monolith.
Current State: I have copy-pasted the code from individual services into a single repo `tofash` under `internal/modules/{module_name}`.

Please execute the following REFACTORING STEPS strictly:

### PHASE 1: FIX IMPORTS & CLEANUP
1.  **Global Import Fix**: Scan all files. Replace all old import paths (e.g., `user-service/...`, `product-service/...`, `gosayu/...`) with the new module path: `tofash/internal/modules/...`.
2.  **Remove HTTP Clients (Internal)**: Delete the folder `internal/modules/*/http_client/` and `internal/adapter/http_client` found in Order/User/Product modules. We will NOT use HTTP for internal communication anymore. (Keep external clients like Midtrans).

### PHASE 2: REFACTOR ORDER SERVICE (Dependency Injection)
Target File: `internal/modules/order/service/order_service.go`
1.  **Update Struct**: Remove `httpClient` field. Add `productSvc` (ProductServiceInterface) and `userSvc` (UserServiceInterface) to the `orderService` struct.
2.  **Update Constructor**: Update `NewOrderService` to accept `productSvc` and `userSvc` as arguments instead of `httpClient`.
3.  **Replace Logic**: In methods like `GetOrderByOrderCode`, `CreateOrder`, and `GetByID`:
    * REMOVE the logic that parses Access Token.
    * REMOVE calls to `o.httpClient...`.
    * REPLACE with direct calls: `o.userSvc.GetByID(ctx, buyerID)` and `o.productSvc.GetByID(ctx, productID)`.
    * Map the returned Entity directly to the response.

### PHASE 3: WIRING IN MAIN.GO
Target File: `cmd/api/main.go`
1.  **Single Entry Point**: Create a robust `main` function.
2.  **Database**: Initialize ONE single `gorm.DB` connection (PostgreSQL) using `config.InitDatabase`.
3.  **Wiring**:
    * Initialize Repositories for User, Product, Order, etc. using the shared `db`.
    * Initialize Services. **CRITICAL**: When initializing `OrderService`, inject the already-initialized `UserService` and `ProductService` instances.
    * Initialize Handlers using the Services.
4.  **Routing**: Setup Echo Groups (`/api/v1/users`, `/api/v1/products`, etc.) and register the handlers.

### PHASE 4: INFRASTRUCTURE
Target File: `docker-compose.yml`
1.  **Simplify**: Keep only 1 App Service, 1 Postgres Service (port 5432), and Redis. Remove all other separate database containers (user_db, product_db, etc.).
2.  **Env Vars**: Ensure the App Service points to the single Postgres instance.

**Constraints:**
* Do not break existing business logic, just change the *communication method* (HTTP -> Direct Call).
* Keep `go.mod` module name as `tofash`.
* Ensure the code compiles after these changes.