# CLAUDE.md

## Project Overview

Insurance claims processing system for India Post Life Insurance (PLI). This is a **specification and seed repository** containing database schemas, API specifications, business requirements, and framework documentation for building a Go microservice.

**Stack:** Go 1.25+ | PostgreSQL 16 | N-API Framework (Uber FX) | OpenAPI 3.0.3

## Repository Structure

```
claims/
├── seed/
│   ├── analysis/                    # Business requirements (70+ rules, 53 functional reqs)
│   │   └── Phase3_Claims_Analysis.md
│   ├── db/                          # PostgreSQL DDL & patches
│   │   ├── claims_database_schema.sql
│   │   ├── claims_schema_enhancement_patch.sql
│   │   ├── performance_optimization_patch.sql
│   │   └── README.md
│   ├── srs/Team_3-Claim/           # SRS/FRS documents (death, maturity, survival, AML, ombudsman)
│   ├── swagger/
│   │   └── claims_api_swagger_complete.yaml  # 130+ endpoints, 16 API categories
│   ├── template/
│   │   └── template.md             # N-API code generation reference
│   └── tool-docs/
│       └── db-README.md            # N-API-DB library docs (pooling, rill pattern)
```

## Architecture & Patterns

The target application follows **Clean Architecture / Hexagonal Pattern**:

```
main.go                    # Entry point
bootstrap/bootstrapper.go  # FX dependency injection
core/
  domain/                  # Pure business entities (no framework deps)
  port/                    # Request/response interfaces
handler/                   # HTTP handlers, routes, request/response DTOs
repo/postgres/             # Repository implementations
db/                        # SQL migrations
configs/                   # Multi-environment YAML configs
```

### Key Conventions

- **Domain types:** `{Resource}` (e.g., `Claim`, `Investigation`)
- **Repositories:** `{Resource}Repository` with `New{Resource}Repository` constructor
- **Handlers:** `{Resource}Handler` embedding `*serverHandler.Base`
- **Routes:** `/{resources}` (plural, lowercase), versioned under `/v1/`
- **JSON/DB fields:** `snake_case` | **Go fields:** `PascalCase`
- **SQL:** Squirrel query builder with `sq.Dollar` placeholder format
- **DI:** Uber FX modules (`FxRepo`, `FxHandler`, `Fxvalidator`)
- **Validation:** Struct tags (`validate:"required"`) with auto-generated validators via `govalid`

### REST API Pattern

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/v1/{resources}` | Create |
| GET | `/v1/{resources}` | List (paginated) |
| GET | `/v1/{resources}/:id` | Get by ID |
| PUT | `/v1/{resources}/:id` | Update |
| DELETE | `/v1/{resources}/:id` | Soft delete |

### Response Format

```json
{
  "status_code": 200,
  "success": true,
  "message": "operation successful",
  "data": {}
}
```

List responses add: `total_records_count`, `returned_records_count`, `skip`, `limit`, `order_by`, `sort_type`.

## Build & Development Commands

```bash
go mod tidy              # Sync dependencies
go build -o app main.go  # Build
go test ./...            # Run all tests
go test -bench=. -benchmem  # Benchmarks
gofmt -w .               # Format code
golangci-lint run        # Lint
cd handler && govalid    # Regenerate validators
```

## Database

**PostgreSQL 16** with 14 core tables, yearly range partitioning, 115+ indexes, 12 triggers, 6 RLS policies, and a materialized view (`mv_daily_claim_stats`).

Key tables: `claims`, `claim_documents`, `investigations`, `appeals`, `aml_alerts`, `claim_payments`, `claim_history`, `claim_sla_tracking`, `ombudsman_complaints`, `policy_bond_tracking`, `freelook_cancellations`.

Query timeouts: Low (2s) for simple queries, Med (5s) for complex/aggregation queries.

## Business Domain

**Claim types:** Death, Maturity, Survival Benefit, Free Look Cancellation.

**Key workflows:** Claim registration → document verification → investigation (conditional) → calculation → approval → disbursement.

**Integrations:** Policy Service, Customer Service, Banking/NEFT, ECMS (documents), Notification (SMS/Email/WhatsApp), AML regulatory filing, DigiLocker, PFMS.

## Configuration

Multi-environment YAML configs: `config.yaml` (base), `config.dev.yaml`, `config.test.yaml`, `config.sit.yaml`, `config.staging.yaml`, `config.prod.yaml`, `config.training.yaml`.

## Code Generation Guidelines

When generating new resources, follow the template in `seed/template/template.md`. Each resource needs:
1. Domain struct in `core/domain/`
2. Port interfaces in `core/port/`
3. Repository in `repo/postgres/` using Squirrel
4. Handler with routes in `handler/`
5. Request/response DTOs with validation tags
6. FX module registration in `bootstrap/bootstrapper.go`
7. SQL migration in `db/`
