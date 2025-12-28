Act as a Senior Go Backend Engineer.
We are transforming this generic e-commerce app into a specialized **FASHION STORE**.
Current State: The `Product` model is too generic (it has a confusing `Variant` int field).

Please EXECUTE the following refactoring steps:

### STEP 1: Update Product Model (Database Schema)
Target File: `internal/modules/product/model/product_model.go`
Modify the `Product` struct to support fashion attributes.
1.  **Remove**: `Variant` (int) field.
2.  **Add** the following fields (with gorm tags):
    * `SKU` (string, unique index) -> e.g., "SHIRT-RED-L"
    * `Size` (string) -> e.g., "S", "M", "L", "42"
    * `Color` (string) -> e.g., "Red", "Navy"
    * `Material` (string) -> e.g., "Cotton"
    * `ImagesJSON` (string, type:text) -> To store multiple image URLs as a JSON string (e.g., `["img1.jpg", "img2.jpg"]`).

### STEP 2: Update Product Entity (Domain Layer)
Target File: `internal/modules/product/entity/product_entity.go`
Update `ProductEntity` struct to match the new business logic.
1.  **Add**: `SKU`, `Size`, `Color`, `Material` (all strings).
2.  **Update**: Change `Image` (string) to `Images` ([]string) to support multiple photos.

### STEP 3: Update Repository Logic (Mapping)
Target File: `internal/modules/product/repository/product_repository.go`
Fix the CRUD methods (`Create`, `Update`, `GetByID`, `GetAll`) to handle the new fields.
* **Crucial**: Implement logic to convert `Entity.Images ([]string)` <-> `Model.ImagesJSON (string)` using `json.Marshal` and `json.Unmarshal`.
    * Before saving to DB: Marshal `[]string` to JSON string.
    * After fetching from DB: Unmarshal JSON string to `[]string`.

### STEP 4: Cleanup & Optimization (Housekeeping)
We are running a Monolith locally, so we don't need microservice artifacts.
1.  **Delete Dockerfiles**: Remove `Dockerfile` inside `internal/modules/*/` (Only keep the root `Dockerfile`).
2.  **Delete HTTP Clients**: Ensure all `http_client` folders inside modules are deleted (except external services like Midtrans).
3.  **Delete Unused Configs**: If there are duplicate config files inside modules, remove them. We only use `internal/config`.

**Constraints:**
* Ensure the code compiles after changes.
* Use `encoding/json` for the image array conversion.