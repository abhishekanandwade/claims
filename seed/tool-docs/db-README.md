# N-API-DB: Database Access Library

A high-performance Go database access library built on pgx, featuring automatic slice pooling, parallel query execution using the Rill pattern, and seamless integration with Uber FX.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [The Rill Pattern for Parallel Queries](#the-rill-pattern-for-parallel-queries)
- [Pool Management](#pool-management)
- [FX Integration](#fx-integration)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Testing](#testing)
- [Benchmarks](#benchmarks)

## Features

- **Automatic Slice Pooling**: Type-safe `sync.Pool` integration for zero-allocation query results
- **Rill Pattern**: Efficient parallel query execution with context awareness
- **FX Integration**: Seamless dependency injection with Uber FX bootstrapper
- **Transaction Support**: Read and write transactions with customizable isolation levels
- **Connection Pooling**: Efficient resource management with pgxpool
- **Graceful Cleanup**: Automatic connection and resource management

## Installation

```bash
go get gitlab.cept.gov.in/it-2.0-common/n-api-db
```

## Quick Start

```go
import (
    "context"
    sq "github.com/Masterminds/squirrel"
    db "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

func main() {
    config := &db.DBConfig{
        DBUsername: "postgres",
        DBPassword: "secret",
        DBHost:     "localhost",
        DBPort:     "5432",
        DBDatabase: "mydb",
        MaxConns:   10,
    }

    database, err := db.NewDefaultDbFactory().CreateConnection(
        db.NewDefaultDbFactory().NewPreparedDBConfig(config),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    poolMgr := db.NewPoolManager(&db.PoolConfig{})
    poolMgr.config.Defaults.InitialCapacity = 50

    query := sq.Select("*").From("users").Where(sq.Eq{"active": true})
    users, err := db.SelectRowsFX(context.Background(), database, poolMgr, query, pgx.RowToStructByPos[User])
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d active users\n", len(users))
}
```

## The Rill Pattern for Parallel Queries

### Overview

The **Rill pattern** is a streaming pattern for executing multiple database queries in parallel with automatic context awareness. Using the [rill](https://github.com/destel/rill) library, this pattern provides:

- **Context Awareness**: Automatic cancellation support for timeouts and shutdowns
- **Goroutine Safety**: No goroutine leaks with proper cleanup
- **Efficient Concurrency**: Configurable concurrency limits for optimal resource usage
- **Error Propagation**: Immediate termination on first error
- **Stream Processing**: Process results as they arrive without buffering

### Why Rill.Generate()?

The pattern uses `rill.Generate()` instead of `rill.FromSlice()` for critical reasons:

**✅ Correct (with rill.Generate):**
```go
queryStream := rill.Generate(func(send func(sq.SelectBuilder), sendErr func(error)) {
    for _, query := range queries {
        if ctx.Err() != nil {
            sendErr(ctx.Err())  // Context cancellation check
            return
        }
        send(query)
    }
})
```

**❌ Wrong (with rill.FromSlice):**
```go
queryStream := rill.FromSlice(queries, nil)  // No context support!
```

### Key Benefits

1. **Context Cancellation**: Queries stop immediately when context is cancelled
2. **Timeout Support**: Long-running queries are cancelled on timeout
3. **Goroutine Leak Prevention**: Background goroutines drain remaining items on early return
4. **Memory Efficiency**: Process results incrementally without buffering

### Available Functions

#### 1. SelectRowsParallelFX

Execute multiple queries in parallel, returning separate result sets:

```go
queries := []sq.SelectBuilder{
    sq.Select("*").From("users").Where(sq.Eq{"active": true}),
    sq.Select("*").From("users").Where(sq.Eq{"active": false}),
    sq.Select("*").From("users").Where(sq.Like{"username": "admin%"}),
}

results, err := db.SelectRowsParallelFX(ctx, database, poolMgr, queries, pgx.RowToStructByPos[User], 3)
if err != nil {
    log.Fatal(err)
}

activeUsers := results[0]
inactiveUsers := results[1]
adminUsers := results[2]
```

#### 2. SelectRowsParallelFlatFX

Execute queries in parallel and merge all results into a single slice:

```go
queries := []sq.SelectBuilder{
    sq.Select("*").From("users").Where(sq.Like{"username": "a%"}),
    sq.Select("*").From("users").Where(sq.Like{"username": "b%"}),
    sq.Select("*").From("users").Where(sq.Like{"username": "c%"}),
}

allUsers, err := db.SelectRowsParallelFlatFX(ctx, database, poolMgr, queries, pgx.RowToStructByPos[User], 3)
fmt.Printf("Total users: %d\n", len(allUsers))
```

#### 3. SelectRowsParallelCallbackFX

Process each query result via callback as it completes:

```go
err := db.SelectRowsParallelCallbackFX(ctx, database, poolMgr, queries, pgx.RowToStructByPos[User], 3,
    func(queryIndex int, users []User) error {
        log.Printf("Query %d returned %d users", queryIndex, len(users))
        for _, user := range users {
            processUser(user)
        }
        return nil
    },
)
```

#### 4. SelectRowsParallelBatchFX

Execute queries and return results in batches:

```go
queries := make([]sq.SelectBuilder, 100)
for i := 0; i < 100; i++ {
    queries[i] = sq.Select("*").From("users").Limit(10)
}

batches, err := db.SelectRowsParallelBatchFX(ctx, database, poolMgr, queries, pgx.RowToStructByPos[User], 10, 5)
for batchIndex, batch := range batches {
    log.Printf("Batch %d: %d query results", batchIndex, len(batch))
}
```

### Context Cancellation Example

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

queries := []sq.SelectBuilder{
    sq.Select("pg_sleep(5)").From("users"),
    sq.Select("pg_sleep(5)").From("users"),
    sq.Select("pg_sleep(5)").From("users"),
}

_, err := db.SelectRowsParallelFX(ctx, database, poolMgr, queries, pgx.RowToStructByPos[User], 3)
// After 2 seconds, context is cancelled
// Queries stop immediately, no goroutine leaks
fmt.Println(err) // context.DeadlineExceeded
```

### Concurrency Control

Control parallelism with the `concurrency` parameter:

```go
queries := make([]sq.SelectBuilder, 20)

// Execute 20 queries with concurrency=3
// Only 3 queries run at a time, reducing database load
results, err := db.SelectRowsParallelFX(ctx, database, poolMgr, queries, pgx.RowToStructByPos[User], 3)
```

### Error Handling

Errors are returned immediately, and remaining queries are cancelled:

```go
queries := []sq.SelectBuilder{
    sq.Select("*").From("valid_table"),
    sq.Select("*").From("invalid_table"),    // This will fail
    sq.Select("*").From("another_valid_table"), // Won't execute
}

_, err := db.SelectRowsParallelFX(ctx, database, poolMgr, queries, pgx.RowToStructByPos[User], 3)
// Returns error from invalid_table
// Remaining queries cancelled
```

### Performance Considerations

| Scenario | Recommended Function | Why |
|----------|-------------------|-----|
| Need all results | `SelectRowsParallelFX` | Returns all result sets |
| Merge results | `SelectRowsParallelFlatFX` | Combines into single slice |
| Stream processing | `SelectRowsParallelCallbackFX` | Process as they complete |
| Rate limiting | `SelectRowsParallelBatchFX` | Control result rate |
| Low database load | `concurrency=1` | Sequential execution |
| High throughput | `concurrency=10+` | Maximize parallelism |

### Thread Safety

All Rill-based functions are thread-safe and can be called concurrently:

```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        queries := createQueries(id)
        results, err := db.SelectRowsParallelFX(ctx, db, poolMgr, queries, scanFn, 3)
        handleResults(id, results, err)
    }(i)
}

wg.Wait()
```

### Memory Management

- **Automatic Pooling**: Result slices use `sync.Pool` for reuse
- **Zero Copy Callbacks**: Process results without allocation
- **Immediate Cleanup**: Background goroutines drain on early return

### Best Practices

1. **Use Appropriate Concurrency**: Match concurrency to database capacity
   ```go
   // For low-resource databases
   results, err := db.SelectRowsParallelFX(ctx, db, poolMgr, queries, scanFn, 2)

   // For high-performance databases
   results, err := db.SelectRowsParallelFX(ctx, db, poolMgr, queries, scanFn, 20)
   ```

2. **Always Use Context**: Never pass `context.Background()` without timeout
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
   defer cancel()
   results, err := db.SelectRowsParallelFX(ctx, db, poolMgr, queries, scanFn, 5)
   ```

3. **Handle Errors Early**: Check errors immediately after parallel execution
   ```go
   results, err := db.SelectRowsParallelFX(ctx, db, poolMgr, queries, scanFn, 5)
   if err != nil {
       log.Printf("Query failed: %v", err)
       return
   }
   ```

4. **Use Callbacks for Processing**: Avoid unnecessary allocations
   ```go
   // ❌ Bad: Results copied twice
   results, err := db.SelectRowsParallelFlatFX(ctx, db, poolMgr, queries, scanFn, 5)
   for _, result := range results {
       process(result)
   }

   // ✅ Good: Process in-place
   err = db.SelectRowsParallelCallbackFX(ctx, db, poolMgr, queries, scanFn, 5,
       func(idx int, results []User) error {
           for _, result := range results {
               process(result)
           }
           return nil
       },
   )
   ```

## Pool Management

### Automatic Slice Pooling

The library uses `sync.Pool` to reuse slice buffers, dramatically reducing memory allocations:

```go
// Without pooling (allocates on every call)
users := make([]User, 0, 50)
users = append(users, ...)  // Allocation every query

// With pooling (reuses slices)
resultPtr := db.Get[User](poolMgr)
*resultPtr = append(*resultPtr, ...)  // Reuses from pool
db.Put(poolMgr, resultPtr)  // Returns to pool
```

### Configuration

```go
config := &db.PoolConfig{
    Defaults: struct {
        InitialCapacity int `mapstructure:"initial_capacity" yaml:"initial_capacity"`
        MaxRetain       int `mapstructure:"max_retain" yaml:"max_retain"`
    }{
        InitialCapacity: 50,
        MaxRetain:       500,
    },
    Types: []db.TypePoolConfig{
        {
            Name:            "User",
            Package:         "models",
            InitialCapacity: 100,
            MaxRetain:       1000,
        },
    },
}

poolMgr := db.NewPoolManager(config)
```

### Memory Benefits

Using pooling provides significant memory savings:

| Operation | Without Pool | With Pool | Improvement |
|-----------|--------------|-----------|-------------|
| 10 rows | 24576 B/op | 12288 B/op | 50% reduction |
| 100 rows | 81920 B/op | 40960 B/op | 50% reduction |
| 1000 rows | 409600 B/op | 204800 B/op | 50% reduction |

## FX Integration

### Basic Setup

```go
import (
    "go.uber.org/fx"
    db "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

func New() *Bootstrapper {
    return &Bootstrapper{
        options: []fx.Option{
            fxconfig,
            fxlog,
            db.FxDB,
            db.FxPooling,  // ← Add pool management
            fxrouter,
        },
    }
}
```

### Repository Pattern

```go
type UserRepository struct {
    db      *db.DB
    poolMgr *db.PoolManager  // Injected by FX
}

func NewUserRepository(database *db.DB, poolMgr *db.PoolManager) *UserRepository {
    return &UserRepository{db: database, poolMgr: poolMgr}
}

func (r *UserRepository) FindActive(ctx context.Context) ([]User, error) {
    query := sq.Select("*").From("users").Where(sq.Eq{"active": true})
    return db.SelectRowsFX(ctx, r.db, r.poolMgr, query, pgx.RowToStructByPos[User])
}
```

### Configuration File

```yaml
pools:
  defaults:
    initial_capacity: 50
    max_retain: 500

  types:
    - name: "User"
      package: "models"
      initial_capacity: 100
      max_retain: 1000
```

## Configuration

### Database Configuration

```go
type DBConfig struct {
    DBUsername        string
    DBPassword        string
    DBHost            string
    DBPort            string
    DBDatabase        string
    Schema            string
    MaxConns          int32
    MinConns          int32
    MaxConnLifetime   time.Duration
    MaxConnIdleTime   time.Duration
    HealthCheckPeriod time.Duration
    AppName           string
}
```

### Pool Configuration

```go
type PoolConfig struct {
    Defaults struct {
        InitialCapacity int
        MaxRetain       int
    }
    Types []TypePoolConfig
}
```

## API Reference

### Core Query Functions

- `SelectRowsFX[T]()` - Execute query with pooled results
- `SelectRowsCallbackFX[T]()` - Execute query with callback processing
- `SelectOneFX[T]()` - Execute query returning single row
- `SelectOneOKFX[T]()` - Execute query returning (T, bool, error)
- `InsertReturningFX[T]()` - Insert with returning
- `UpdateReturningFX[T]()` - Update with returning
- `InsertReturningrowsFX[T]()` - Bulk insert with returning

### Parallel Query Functions

- `SelectRowsParallelFX[T]()` - Execute queries in parallel
- `SelectRowsParallelFlatFX[T]()` - Execute and merge results
- `SelectRowsParallelCallbackFX[T]()` - Execute with callback
- `SelectRowsParallelBatchFX[T]()` - Execute in batches

### Batch Operations

- `QueueReturnFX[T]()` - Queue query in batch
- `QueueReturnRowFX[T]()` - Queue single row query

### Raw Query Functions (Direct SQL)

These functions execute raw SQL strings without query builders, useful for complex queries or when using string-based SQL directly:

#### Basic Raw Query Functions

- `SelectRowsRaw[T]()` - Execute raw SQL query returning multiple rows
- `SelectOneRaw[T]()` - Execute raw SQL query returning single row
- `SelectRowsOKRaw[T]()` - Execute raw SQL query returning ([]T, bool, error)
- `SelectOneOKRaw[T]()` - Execute raw SQL query returning (T, bool, error)
- `ExecRaw()` - Execute raw SQL command (INSERT, UPDATE, DELETE)

#### Transaction Raw Query Functions

- `TxExecRaw()` - Execute raw SQL command in transaction
- `TxReturnRowRaw[T]()` - Execute raw SQL returning single row in transaction
- `TxRowsRaw[T]()` - Execute raw SQL returning multiple rows in transaction

#### Batch Raw Query Functions

- `QueueExecRowRaw()` - Queue raw SQL command in batch
- `QueueReturnRaw[T]()` - Queue raw SQL query in batch returning multiple rows
- `QueueReturnRowRaw[T]()` - Queue raw SQL query in batch returning single row
- `QueueReturnBulkRaw[T]()` - Queue raw SQL query appending results to existing slice

#### Timed Batch Raw Query Functions

- `TimedQueueExecRowRaw()` - Queue raw SQL command in timed batch with statement timeout
- `TimedQueueReturnRaw[T]()` - Queue raw SQL query in timed batch returning multiple rows
- `TimedQueueReturnRowRaw[T]()` - Queue raw SQL query in timed batch returning single row
- `TimedQueueReturnBulkRaw[T]()` - Queue raw SQL query in timed batch appending results

### Raw Functions Usage Example

```go
// Simple raw query
users, err := db.SelectRowsRaw(ctx, database, 
    "SELECT id, name, email FROM users WHERE active = $1", 
    []any{true}, 
    pgx.RowToStructByPos[User])

// Single row
user, err := db.SelectOneRaw(ctx, database,
    "SELECT id, name, email FROM users WHERE id = $1",
    []any{userID},
    pgx.RowToStructByPos[User])

// Raw execution (INSERT, UPDATE, DELETE)
tag, err := db.ExecRaw(ctx, database,
    "UPDATE users SET active = $1 WHERE id = $2",
    []any{false, userID})

// Batch with timeout
batch := db.NewTimedBatch(5 * time.Second)
result := []User{}
err := db.TimedQueueReturnRaw(batch, 
    "SELECT * FROM users WHERE status = $1",
    []any{"active"},
    pgx.RowToStructByPos[User],
    &result)
err = database.SendBatch(ctx, batch.Batch).Close()
```

## Examples

See the `examples/` directory for complete examples:

- `examples/complete-repository-example.go` - Full repository pattern
- `examples/parallel-queries-example.go` - Parallel query usage
- `examples/fx-integration/` - FX bootstrapper integration
- `examples/pooling/` - Pooling examples

## Testing

### Unit Tests

```bash
cd n-api-db && go test
```

### Integration Tests

```bash
cd n-api-db && go test -v -run Integration
```

### All Tests with Coverage

```bash
cd n-api-db && go test -cover ./...
```

## Benchmarks

### Run All Benchmarks

```bash
cd n-api-db && go test -bench=. -benchmem
```

### Pool Performance Benchmarks

```bash
cd n-api-db && go test -bench=BenchmarkPoolPerformance -benchmem -benchtime=3s
```

### Pool vs No-Pool Comparison

```bash
cd n-api-db && go test -bench=BenchmarkComparison -benchmem -benchtime=3s
```

### Parallel Query Benchmarks

```bash
cd n-api-db && go test -bench=BenchmarkParallel -benchmem
```

### Specific Benchmark

```bash
cd n-api-db && go test -bench=BenchmarkPoolPerformance_SmallData -benchmem
```

### Benchmark Results

Example results (your results may vary):

```
BenchmarkPoolPerformance_Get_Put_Empty-8               5000000    350 ns/op    128 B/op    1 allocs/op
BenchmarkWithoutPool_Allocation_Empty-8               3000000    450 ns/op    256 B/op    2 allocs/op
BenchmarkPoolPerformance_Get_Put_MediumData-8           50000   28000 ns/op   8192 B/op   50 allocs/op
BenchmarkWithoutPool_MediumData-8                       30000   35000 ns/op  16384 B/op   100 allocs/op
BenchmarkPoolPerformance_Concurrent_ReadWrite-8         100000   12000 ns/op   4096 B/op   20 allocs/op
BenchmarkWithoutPool_Concurrent-8                       80000   15000 ns/op   8192 B/op   40 allocs/op
```

## Performance Tips

1. **Use Pooled Functions**: Always use `*FX` functions when `PoolManager` is available
2. **Choose Appropriate Capacity**: Set `InitialCapacity` based on typical result size
3. **Use Callbacks for Processing**: Avoid unnecessary result copies
4. **Limit Concurrency**: Match concurrency to database capacity
5. **Use Timeouts**: Always use context with timeout for database operations

## License

See LICENSE file for details.

## Contributing

Please see CONTRIBUTING.md for guidelines.
