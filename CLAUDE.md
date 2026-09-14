# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project purpose

A Go backend clone of Uzum (an Uzbek e-commerce marketplace), built as a learning project for **Domain-Driven Design** and **Bounded Context / Layered Architecture** — not intended for production. If asked to explain patterns here, it's worth noting NestJS/TypeORM analogies (Module ≈ bounded context, Service ≈ use case, Entity/Repository map closely), since that's the mental model the project was designed against.

## Common commands

```bash
make up             # docker-compose up -d (postgres, redis)
make down           # docker-compose down
make run            # go run cmd/api/main.go
make test           # go test ./...
make swagger        # swag init -g cmd/api/main.go -o docs
make migrate-up     # golang-migrate up, reads DB_* from .env
make migrate-down   # golang-migrate down
```

- Run a single test: `go test ./internal/cart/domain/... -run TestName -v`
- Full swagger regen command (used by `make swagger`, needs `--parseDependency --parseInternal` to resolve cross-package DTOs): `swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal`
- Server listens on `:8080` (hardcoded in `cmd/api/main.go`, though `APP_PORT` exists in config — check both if changing the port). Swagger UI at `http://localhost:8080/swagger/index.html`.
- Config is loaded from `.env` via `pkg/config` (see `.env.example` for the full list of required vars: DB_*, REDIS_*, JWT_*, AWS_* for S3 uploads).

## Architecture

Each business capability is a **bounded context** under `internal/<context>/`, wired together only in `cmd/api/main.go`. Contexts implemented so far: `identity` (users/auth/RBAC — complete), `catalog` (products/categories/events — mostly done), `cart`, `wishlist`, `ordering` (checkout, admin manual orders, cancellation), `review` (product reviews, gated on a delivered purchase), `gallery` (simple image posts, no cross-context dependencies). `internal/shared/acl`/`internal/shared/event` remain empty placeholder directories — no code lives there yet; `review` is the first context to establish an explicit ACL pattern (a domain-level port + a focused read-only use case on the other side) rather than reaching into another context's infrastructure directly.

Not everything under `internal/` is a bounded context: `internal/admin/dashboard` is a deliberate exception — a thin, read-only reporting layer with no domain/application/infrastructure split and no invariants, allowed to query `catalog` and `ordering`'s tables directly via raw SQL since there's no business logic, only aggregation. Don't use it as a precedent for skipping the four-layer structure anywhere else.

Every context repeats the same four-layer structure:

```
domain/  →  application/ (use cases)  →  interfaces/http/ (handlers, DTOs, router)
                ↑
        infrastructure/ (postgres repos, security, etc. — implements domain ports)
```

- **`domain/`** — pure business logic, zero HTTP/DB/framework imports. Entities own their invariants (e.g. `User.ChangeRole()`), value objects are immutable structs with unexported fields and constructor validation (e.g. `internal/shared/money.Money` — cannot be built or mutated except through its methods, so it can't be deserialized directly from JSON). Repository interfaces (ports) are declared here (`domain/repository.go`), not in infrastructure.
- **`application/`** — one use case per file, one exported struct/constructor per use case (e.g. `NewAddItemUseCase`), not grouped into a single fat service. Use cases depend only on domain interfaces, never on infrastructure or HTTP types directly.
- **`infrastructure/`** — concrete adapters: `postgres/` repos built on `sqlx` with named parameters (`:field`), `security/` for JWT/bcrypt. Swappable without touching domain or application.
- **`interfaces/http/`** — chi handlers, DTOs, per-context `router.go`, swagger annotations, and a `*_errors.go` file per context that maps domain sentinel errors to HTTP status codes via `errors.Is` (see `cart/interfaces/http/cart_errors.go` for the pattern). Domain objects are never serialized directly as DTOs — HTTP DTOs use plain types (e.g. `Amount float64`/`Currency string` instead of `Money`), and use cases construct the domain value objects from them.

**Module wiring**: each context exposes `module.go` at its root with `Config` (dependencies in) and `NewModule(cfg) *Module`. `cmd/api/main.go` constructs shared infra (postgres `*sqlx.DB`, redis client, S3 uploader), then builds each context's module in dependency order, then mounts routers under `/api/v1/...`. Cross-context dependencies are passed as constructor args (e.g. `cart` and `wishlist` both take `identityModule.TokenService`); note that `cart`'s module also directly imports `catalog/infrastructure/postgres` for product lookups rather than going through an ACL — there's no enforced boundary between contexts yet beyond directory structure.

**Routing gotcha — single-Router vs multi-Router (discovered in `ordering`, Stage 4; rediscovered in `review` and `gallery`)**: `chi.Router.Mount(pattern, handler)` panics if `pattern` is registered with the exact same literal string twice. `catalog` is mounted first and claims the **bare `/api/v1`** prefix (its own `interfaces/http/router.go` handles `/products`, `/categories`, `/events`, etc. internally) — so no other context can mount a second router at bare `/api/v1` to get that same "one router, many internal sub-routes" convenience. Mounting at a *more specific* prefix that merely shares a leading path segment with something catalog handles is fine and does **not** panic — chi's radix-tree routing correctly prioritizes the more specific pattern (verified directly: `review`'s router mounted at `/api/v1/products/{id}/reviews` coexists with catalog's bare `/api/v1` mount handling `/api/v1/products/{id}`, and the `{id}` URL param propagates through to the more specific router's handler).

Practical consequence for module design: if a context's routes all fit under **one** prefix nobody else has claimed (e.g. `/api/v1/wishlist`, `/api/v1/cart`), the simple pattern works — `Module.Router chi.Router`, one `NewRouter(...)` constructor, one `r.Mount(...)` call in `main.go`. But a context that needs several logically distinct resource groups at **different** prefixes — typically because it would otherwise want its own internal sub-routing under its own bare top-level prefix the way `catalog` does, but can't since that prefix is taken — must instead expose one named `chi.Router` field per prefix on its `Module` (a **multi-Router** module), each wired to its own `New<Name>Router(...)` constructor and `r.Mount()`ed separately in `main.go`. Every context mounted after `catalog` that needs more than one route group runs into this; there's no way around it short of not using chi's `Mount`.

Current pattern per context (check here before assuming `Module.Router` — the wrong assumption is exactly what causes the chi panic):

| Pattern | Contexts | Prefixes |
|---|---|---|
| **Single-Router** (`Module.Router`, one `r.Mount`) | `identity`, `catalog`, `cart`, `wishlist` | `/api/v1/auth`; bare `/api/v1` (catalog — first to claim it); `/api/v1/cart`; `/api/v1/wishlist` |
| **Multi-Router** (named `chi.Router` fields, each mounted separately) | `ordering` — `CheckoutRouter`, `OrdersRouter`, `AdminOrdersRouter` | `/api/v1/checkout`; `/api/v1/orders`; `/api/v1/admin/orders` |
| | `review` — `ReviewsRouter`, `ProductReviewsRouter`, `AdminReviewsRouter` | `/api/v1/reviews`; `/api/v1/products/{id}/reviews`; `/api/v1/admin/reviews` |
| | `gallery` — `GalleryRouter`, `AdminGalleryRouter` | `/api/v1/gallery`; `/api/v1/admin/gallery` |

`internal/admin/dashboard` isn't a bounded context (no `Module`), but hits the same constraint at the call site: `dashboard.NewRouter(db, tokenService)` is built and `r.Mount("/api/v1/admin/dashboard", ...)`ed directly in `main.go` rather than through a module.

When adding a new context: default to single-Router — it's simpler. Only split into multi-Router once you actually need a second prefix, and when you do, mount each router as its own `r.Mount(...)` line in `main.go` (see the existing calls for the pattern) rather than trying to combine them.

**Auth/RBAC** (`identity` context): JWT carries `userID` + `role`. `middleware.Authenticate` validates the token and injects both into the request context using **typed context keys** (`middleware.ContextKey`, not raw strings — see `identity/interfaces/http/middleware/context_keys.go`) to avoid cross-package key collisions. `middleware.RequireRole` gates role-specific endpoints. Every other context's router applies `middleware.Authenticate(tokenService)` from `identity` — there's no per-context auth reimplementation.

**Stack**: chi (router), sqlx + lib/pq (Postgres), golang-migrate (migrations in `migrations/`, sequentially numbered), golang-jwt/jwt/v4, swaggo/swag (Swagger, generated into `docs/` — don't hand-edit), aws-sdk-go-v2/s3 (media uploads via `internal/shared/media`), google/uuid, go-redis/v9.

**Shared kernel**: `Money` (amount + currency, with `NewMoney`/`Add`/`Amount`/`Currency`) lives in `internal/shared/money` — it used to live in `catalog/domain` and get imported directly by `ordering/domain` and others, which is exactly the kind of incidental cross-context domain coupling this project otherwise avoids. If another genuinely-shared value object shows up (used as a domain type, not just a DTO field, by 2+ contexts), it belongs in `internal/shared/`, not in whichever context happened to define it first.

## Current status (update this after every session)

**Bounded contexts:**

| Context | Status |
|---|---|
| `identity` | Complete — JWT auth, RBAC middleware, typed context keys |
| `wishlist` | Complete and tested |
| `catalog` | Mostly complete — Products (Cloudinary uploads, soft delete, `Money` VO), Categories (image support, soft delete, Cloudinary rollback on failure). See known issues below. |
| `cart` | Complete and tested — domain, application (use cases), infrastructure (repository), and interfaces/http (handlers) all written. One persistent cart per user (1:1, lazy-created); price is read live from `catalog` at read time, never snapshotted; duplicate `AddItem` upserts quantity; promocode/discount deferred to `ordering`. |
| `ordering` | Complete and tested — checkout (from cart) and admin manual orders (cart-less, duplicate `ProductID`s merged before stock reservation) both snapshot name/price from `catalog` and decrement stock via `StockReserver`; payment/delivery status updates (`AdvanceDelivery` is rank-based, forward-only); order cancellation (`Order.Cancel()`, Preparing-only, restores stock) as a dedicated one-way terminal transition kept out of `AdvanceDelivery`'s rank logic. Still uses `cart`'s direct-infra-import pattern for `catalog` (no ACL) via `StockReserver`. |
| `review` | Complete (MVP) and tested — one review per user+product, gated on a delivered purchase. First context with a real ACL boundary: `review/domain.DeliveredPurchaseChecker` port, implemented in `review/infrastructure/postgres` by calling `ordering/application.HasDeliveredProductUseCase` (never `ordering/domain` or its DB directly). `GetUserReviewsUseCase` exists but has no HTTP route yet (only product-reviews and submit are wired). Admin moderation: `DELETE /api/v1/admin/reviews/{id}` (`AdminReviewsRouter`, gated by `middleware.RequireRole(identitydomain.RoleAdmin)`) lets an admin delete any user's review outright — no ownership check in `DeleteReviewUseCase`, access control is enforced entirely at the router. |
| `gallery` | Complete and tested — no cross-context dependencies (simplest context: no pricing/stock). Up to 3 images per post, stored as a `GalleryImage{URL, PublicID}` JSONB column (same pattern as `catalog.Product.Images`) so `PublicID` survives for S3 cleanup on delete; `ImageURLs()` strips `PublicID` before it reaches HTTP DTOs. `Description` is optional/trimmed, no validation error (matches `ordering.CustomerInfo.Note`'s precedent). Upload filenames are generated as `uuid.New().String() + filepath.Ext(header.Filename)` from the start — the Category fix, never the raw-filename bug still open on `catalog`'s Product handler (see known issues). |
| `internal/admin/dashboard` | Complete — deliberately **not** a bounded context (no domain/application/infrastructure split, no module.go): plain handler functions running raw SQL directly against `*sqlx.DB` across `catalog`'s and `ordering`'s tables (justified since there's no business logic, only read-only aggregation). Three admin-only endpoints: summary (product count + units sold + revenue), revenue-history (day/month `date_trunc` buckets), low-stock (top 5). The `payment_status = 'paid' AND delivery_status <> 'cancelled'` filter is repeated in both summary and revenue-history — a paid-then-cancelled order must not count as revenue. Has a real integration test (`internal/admin/dashboard/dashboard_integration_test.go`, white-box `package dashboard`) that seeds via the actual `catalog`/`ordering` repositories against a live Postgres, asserts via before/after deltas (robust against pre-existing dev data), and skips itself if Postgres isn't reachable — the first and only DB-backed test in the repo; every other context tests against in-memory fakes. |

**Known issues (fix before or alongside next context):**

- `ProductCard.vue`: `.toFixed()` bug caused by a `Money` shape mismatch between backend (`{ amount, currency }`) and frontend TS types (`~/types/product.ts`, `ProductCard.vue`) — needs the two files open side by side to diagnose.
- Migrations 4 and 5 are missing `.down.sql` files.
- Repository layer has mixed parameter styles: some files use `database/sql` positional (`$1`), others use `sqlx` named (`:field`). Standardize on named (`:field`) per the architecture convention above.
- `catalog`'s Product handler generates the S3 key as `"product-images/" + uuid.NewString()` — better than the raw-filename bug `Category` originally had, but still drops the file extension (unlike `Category`'s and now `gallery`'s `uuid.New().String() + filepath.Ext(header.Filename)`). Worth aligning for consistency, though it's not the security issue this line used to describe.

**Next up:**

1. `internal/admin/dashboard`'s revenue figures assume a single currency across all orders (no `currency` grouping/check) — fine today since the app only ever uses UZS, but would silently mix currencies if that ever changes.
2. `gallery` has no update path yet (edit description, replace/add/remove one image) — only create/list/delete exist, matching what was asked.
3. Expose `review.GetUserReviewsUseCase` over HTTP (a "my reviews" endpoint) if/when the frontend needs it — deferred since it wasn't part of the original ask.
4. Revisit `cart`'s and `ordering`'s direct import of `catalog/infrastructure/postgres` (via `StockReserver`/product lookups) now that `review` has established what a real ACL boundary looks like in this codebase.
5. `move_to_cart` (or equivalent) as a cross-context use case, going through `cart`'s application-layer interface, not its domain or repository directly.
6. Search + pagination on `GetAll` endpoints across contexts — flagged, not implemented.
7. Redis caching — deferred until real traffic exists; when added, wrap via a `CachedProductRepository` decorator without changing the existing repository interface.
