# Coding Conventions

**Analysis Date:** 2026-03-31

## Language & Environment

**Primary Language:** Go 1.25.7

**Code Organization:** Standard Go project layout with module path `k8c.io/kubermatic/v2`

## Naming Patterns

**Files:**
- `*_test.go` for unit tests (co-located with source code)
- `*_ce.go` and `*_ee.go` for Community Edition vs Enterprise Edition variants
- PascalCase for public types and functions
- camelCase for unexported (internal) functions and variables

**Functions:**
- Public functions: PascalCase (e.g., `NewValidator`, `ValidateCreate`, `Build`)
- Private functions: camelCase (e.g., `buildValidationDependencies`, `validateProjectRelation`)
- Constructor functions follow `NewTypeName` pattern (e.g., `NewValidator`, `NewMutator`, `NewClientBuilder`)
- Public methods: PascalCase on receiver (e.g., `(v *validator) ValidateCreate`, `(r rawClusterGen) Build`)

**Types & Structs:**
- PascalCase for all types: `Mutator`, `validator`, `AdmissionHandler`, `Cluster`, `Seed`
- Private receiver types use lowercase first letter: `validator`, `rawClusterGen`
- Struct fields are PascalCase and exported when used across packages

**Variables & Constants:**
- Constants: PascalCase for exported (e.g., `PreventDeletionAnnotation`, `DefaultEtcdClusterSize`)
- Constants: camelCase for unexported (e.g., `defaultingTemplateName`, `addonPrefix`)
- Package-level variables: var blocks used for shared test fixtures
- Const blocks used for compile-time constants

## Code Style

**Formatting:**
- Tool: `gofmt` (Go's built-in formatter) - enforced via `.golangci.yml`
- Line length: No explicit limit enforced in linter, but conventional Go style applies
- Indentation: Go standard (tab characters)

**Linting:**
- Tool: `golangci-lint` (v2 configuration)
- Config: `.golangci.yml`
- Key enabled linters:
  - `errcheck` - catches unchecked errors
  - `govet` - catches common Go mistakes
  - `staticcheck` - comprehensive static analysis
  - `gocyclo` - cyclomatic complexity checking (high-complexity functions explicitly allowed for reconcilers)
  - `errorlint` - error wrapping best practices
  - `unconvert` - removes unnecessary type conversions

**Exception Handling:**
- Disabled check: `ST1005` - allows error strings starting with capitalized type names (K8s convention)
- High complexity allowed for specific functions:
  - `reconcile()` methods on reconcilers
  - `InitializeCloudProvider()` for cloud providers
  - `main()` functions
  - Test initialization functions like `initTestEndpoint`

## Import Organization

**Order:**
1. Standard library imports (e.g., `context`, `fmt`, `testing`)
2. Third-party dependencies (e.g., `go.uber.org/zap`, `encoding/json` packages)
3. Kubernetes ecosystem (e.g., `k8s.io`, `sigs.k8s.io`)
4. Internal project imports (e.g., `k8c.io/kubermatic/v2/...`)

**Import Aliases (Required via golangci-lint importas linter):**
- `k8c.io/kubermatic/sdk/v2/api/v1` → `apiv1`
- `k8c.io/kubermatic/sdk/v2/api/v2` → `apiv2`
- `k8c.io/kubermatic/sdk/v2/apis/kubermatic/v1` → `kubermaticv1`
- `k8c.io/kubermatic/sdk/v2/apis/apps.kubermatic/v1` → `appskubermaticv1`
- `k8c.io/kubermatic/v2/pkg/util/errors` → `utilerrors`
- Kubernetes: `k8s.io/api/{group}/{version}` → `{group}{version}` (e.g., `corev1`, `appsv1`)
- `k8s.io/apimachinery/pkg/apis/meta/v1` → `metav1`
- `k8s.io/apimachinery/pkg/api/errors` → `apierrors`
- `k8s.io/apimachinery/pkg/util/errors` → `kerrors`
- `sigs.k8s.io/controller-runtime/pkg/client` → `ctrlruntimeclient`
- `sigs.k8s.io/controller-runtime/pkg/client/fake` → `ctrlruntimefakeclient`
- All import aliases are enforced with `no-unaliased: true` setting

**File Header:**
- Apache License 2.0 comment block (required)
- Copyright notice with year and contributor attribution
- Located at top of every source file

## Error Handling

**Patterns:**
- Explicit error checking: `if err != nil { return nil, err }`
- Kubernetes API errors wrapped with type-specific functions:
  - `field.InternalError(nil, err)` for internal validation errors
  - `field.Error` return type for webhook validation
  - `apierrors.IsNotFound(err)` for checking specific K8s error types
- Error wrapping: `fmt.Errorf("message: %w", err)` for contextual error messages
- Multiple error aggregation: `errs.ToAggregate()` for field validation lists
- HTTP status codes in admission webhooks: `http.StatusBadRequest`, `http.StatusInternalServerError`

**Error Recovery:**
- Webhook handlers return specific HTTP status codes and error messages
- Validation errors use `*field.Error` type for structured error reporting
- No panic recovery patterns visible - let panics propagate to handlers

## Logging

**Framework:** `go.uber.org/zap` (SugaredLogger)

**Logger Injection:**
- Passed as constructor parameter to handlers: `log *zap.SugaredLogger`
- Used primarily in admission webhook handlers for request/response logging
- Log fields: error context, operation type, request UID

**Patterns:**
- Errors logged with context: `h.log.Error(mutateErr, "cluster mutation failed")`
- No structured logging visible in validation/mutation logic itself
- Logging is minimal - mainly used at handler boundaries

## Comments

**When to Comment:**
- Function/type comments: Always for public (exported) functions and types
- Comment format for public items: `// TypeName describes what it does` (starts with name)
- Comment format for private items: `// lowercase describes what it does`
- Inline comments: Used for non-obvious logic or implementation details
- Block comments: Used for explaining test expectations and complex rules

**Examples from codebase:**
```go
// validator for validating Kubermatic Cluster CRD.
type validator struct { ... }

// NewValidator returns a new cluster validator.
func NewValidator(...) *validator { ... }

// disableProviderValidation is only for unit tests, to ensure no
// provider would phone home to validate dummy test credentials
disableProviderValidation bool
```

**JSDoc/Comments:**
- No GoDoc format enforced at package level
- Comments are descriptive and human-readable
- Long explanations broken into multiple comment lines with `//`

## Function Design

**Size Guidelines:**
- Most functions: 10-50 lines
- Large functions: 1500+ lines are not uncommon (e.g., `resources.go`, `reconciler.go`)
  - These are consolidated resource builders or complex reconcilers
  - Eligible for complexity exclusions in linter
- Strategy: Functions grouped by domain concern rather than split for size alone

**Parameters:**
- Constructor functions typically take: logger, client, config dependencies
- Receiver types follow Go method convention: `func (r *ReceiverType) Method(...)`
- Context always first parameter for async operations: `func (...) Method(ctx context.Context, ...)`
- Error returns at end: `(...) (result Type, err error)`

**Return Values:**
- Multiple returns: `(result, error)` or `(result, *field.Error)` for validation
- Pointers returned for mutable types: `*Mutator`, `*validator`, `*Cluster`
- Nil checking for optional values: `if oldCluster != nil { ... }`
- Error types structured: validation uses `*field.Error`, webhooks use `error`

**Builder Pattern:**
- Test objects use builder pattern: `rawClusterGen` struct with `Build()` and `BuildPtr()` methods
- Fluent builder methods return receiver for chaining (not visible in this codebase but possible)
- Test helpers use `Do()` suffix: `rawClusterGen{...}.Do()` as shorthand for `Build()`

## Module Design

**Exports:**
- Package-level functions: `New*` constructors, public methods
- Private implementation details: Receiver types with lowercase names
- Strategic exports: Only what's needed by other packages

**Barrel Files/Init:**
- No barrel index files observed
- Package organization by domain (webhook, controller, provider, resources)
- Each domain has multiple specialized packages

**Package Structure Examples:**
- `pkg/webhook/cluster/validation/` - Cluster validation webhooks
- `pkg/webhook/cluster/mutation/` - Cluster mutation webhooks
- `pkg/provider/cloud/` - Cloud provider implementations
- `pkg/test/fake/` - Test utilities and fixtures
- `pkg/test/diff/` - Diff utilities for test assertions

## Code Examples

**Function Signature Pattern:**
```go
func (v *validator) ValidateCreate(ctx context.Context, cluster *kubermaticv1.Cluster) (admission.Warnings, error) {
    // Implementation
    return nil, errs.ToAggregate()
}
```

**Type Definition Pattern:**
```go
type Mutator struct {
    client       ctrlruntimeclient.Client
    seedGetter   provider.SeedGetter
    configGetter provider.KubermaticConfigurationGetter
    caBundle     *x509.CertPool
    disableProviderMutation bool
}

func NewMutator(client ctrlruntimeclient.Client, ...) *Mutator {
    return &Mutator{
        client: client,
        // ... field initialization
    }
}
```

**Constructor Convention:**
- Name: `New{TypeName}`
- Receiver parameter handling: Fields set in consistent order matching struct definition
- Returns pointer to struct: `*TypeName`

---

*Convention analysis: 2026-03-31*
