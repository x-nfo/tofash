Act as a Senior Go Backend Engineer.
We are optimizing this Monolith application for a 2 vCPU VPS.
The current code still attempts to connect to Elasticsearch, which causes errors or high memory usage.

Please EXECUTE the following refactoring to completely REMOVE Elasticsearch:

### STEP 1: Update `cmd/api/main.go`
Target: `cmd/api/main.go`
1.  **Remove Initialization**: Delete the lines `config.InitElasticsearch(...)` and the `esClient` variable.
2.  **Update Product Wiring**: Change `productRepo.NewProductRepository(db, esClient)` to `productRepo.NewProductRepository(db)`.
3.  **Update Order Wiring**:
    * Delete `elasticRepo := orderRepo.NewElasticRepository(esClient)`.
    * Update `orderService.NewOrderService(...)` to remove `elasticRepo` from the arguments.

### STEP 2: Refactor Product Repository
Target: `internal/modules/product/repository/product_repository.go`
1.  **Struct**: Remove `esClient` field from `productRepository` struct.
2.  **Constructor**: Update `NewProductRepository` signature to remove `*elasticsearch.Client`.
3.  **Search Method**: Rewrite `SearchProducts` method. Instead of calling Elasticsearch, it should simply call `p.GetAll(ctx, query)` (which already implements SQL `ILIKE` search).
4.  **Delete Method**: Remove the code block that attempts to delete data from Elasticsearch.

### STEP 3: Refactor Order Service
Target: `internal/modules/order/service/order_service.go`
1.  **Struct**: Remove `elasticRepo` field from `orderService` struct.
2.  **Constructor**: Update `NewOrderService` signature to remove `elasticRepo`.
3.  **Methods**: In `GetAll` and `GetAllCustomer`, remove the code block that calls `o.elasticRepo.Search...`. Direct the flow to solely use `o.repo.GetAll(...)`.

### STEP 4: Cleanup
1.  **Delete File**: Delete `internal/modules/order/repository/elastic_repository.go` if it exists.
2.  **Delete Config**: Delete `internal/config/elastic.go` if it exists.

**Goal:** The application must compile and run successfully using ONLY PostgreSQL for data storage and searching.