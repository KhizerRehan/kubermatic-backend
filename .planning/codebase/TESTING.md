# Testing Patterns

**Analysis Date:** 2026-03-31

## Test Framework

**Runner:**
- Standard Go `testing` package with `go test` command
- Config: Makefile targets in `/Makefile`
- Build tags used for test variants: `unit`, `integration`, `e2e`, plus edition tags

**Assertion Library:**
- No external assertion library (no `testify`, `assert`, `require` imports)
- Manual assertion patterns using `if !reflect.DeepEqual(...)` or `fmt.Sprintf` for error messages

**Run Commands:**
```bash
make test                      # Run all tests (unit tests + build tests)
make run-tests                 # Run unit tests specifically
make test-integration          # Run integration tests
make build-tests               # Compile e2e and integration tests without running
```

**Under the hood (from `hack/run-tests.sh`):**
```bash
go_test unit_tests -tags "unit,${KUBERMATIC_EDITION}" -timeout 30m -race -v ./pkg/... ./cmd/... ./codegen/...
cd sdk && go_test unit_tests -tags "unit,${KUBERMATIC_EDITION}" -timeout 30m -race -v ./...
```

**Race Detection:**
- Enabled by default: `-race` flag in all test runs
- Timeout: 30 minutes per test suite
- Verbose output: `-v` flag always enabled

## Test File Organization

**Location:**
- Co-located with source code in same package
- Test files use `_test.go` suffix
- 237 test files found in `pkg/` directory (as of analysis date)

**Naming:**
- Test files: `{module}_test.go` (e.g., `mutator_test.go`, `validation_test.go`)
- Test functions: `Test{FunctionName}` (e.g., `TestMutator`, `TestHandle`, `TestValidate`)
- Test cases within function: Named in struct field `name` within test table

**Structure - Example from `pkg/webhook/cluster/mutation/mutator_test.go`:**
```
mutator_test.go
├── Package declaration (mutation)
├── Import block (standard → kubernetes → internal)
├── Package-level var block for fixtures
│   └── config, seed, defaultPatches, etc.
└── func TestMutator(t *testing.T)
    └── Table-driven test (slice of struct with test cases)
```

## Test Structure

**Suite Organization - Table-Driven Pattern:**

```go
func TestMutator(t *testing.T) {
	tests := []struct {
		name                   string
		oldCluster             *kubermaticv1.Cluster
		newCluster             *kubermaticv1.Cluster
		defaultClusterTemplate *kubermaticv1.ClusterTemplate
		wantAllowed            bool
		wantPatches            []jsonpatch.JsonPatchOperation
	}{
		{
			name: "Create cluster sets default component settings",
			newCluster: rawClusterGen{
				Name: "foo",
				CloudSpec: kubermaticv1.CloudSpec{...},
				// ... fields
			}.Do(),
			wantAllowed: true,
			wantPatches: []jsonpatch.JsonPatchOperation{...},
		},
		// More test cases...
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Test execution
		})
	}
}
```

**Key Patterns:**
- Table-driven tests with `t.Run()` for named subtests
- Struct-based test case definitions with descriptive `name` fields
- Expected result fields prefixed with `want*` (e.g., `wantAllowed`, `wantPatches`)
- Actual result variables typically lowercase (e.g., `got`, `actual`, `result`)

**Setup Pattern:**
- Package-level var blocks initialize shared fixtures (config, seed, default patches)
- Fixtures reused across tests to avoid duplication
- Test-specific objects constructed inline within test cases

**Teardown Pattern:**
- No explicit teardown visible
- Fake clients cleaned up implicitly (no resource cleanup required)
- Focus on isolated unit tests without state management

**Assertion Pattern - Manual Comparison:**
```go
// DeepEqual checking
if !reflect.DeepEqual(expected, actual) {
    t.Fatalf("expected %v, got %v", expected, actual)
}

// Semantic equality (for K8s objects)
if !apiequality.Semantic.DeepEqual(expected, actual) {
    t.Errorf("semantic equality failed")
}

// Diff utilities from test/diff package
actualDiff := diff.ObjectDiff(expectedObj, actualObj)
if actualDiff != "" {
    t.Fatalf("Objects differ:\n%s", actualDiff)
}

// Explicit comparisons
if got != want {
    t.Errorf("got %v, want %v", got, want)
}
```

## Test Data & Fixtures

**Test Object Builders:**
- Builder pattern for complex K8s objects: `rawClusterGen` struct in test files
- Builder methods: `Build()` returns value, `BuildPtr()` returns pointer
- Shorthand: `Do()` method as alias for `Build()`

**Example from `validation_test.go`:**
```go
type rawClusterGen struct {
	Name                string
	Namespace           string
	CloudProviderName   string
	Labels              map[string]string
	ExposeStrategy      string
	// ... more fields
}

func (r rawClusterGen) BuildPtr() *kubermaticv1.Cluster {
	c := r.Build()
	return &c
}

func (r rawClusterGen) Build() kubermaticv1.Cluster {
	// Construct and return cluster object
}
```

**Package-Level Fixtures - in `mutator_test.go`:**
```go
var (
	config = kubermaticv1.KubermaticConfiguration{}
	seed   = kubermaticv1.Seed{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-seed",
			Namespace: "kubermatic",
		},
		// ... detailed setup
	}
	defaultPatches = []jsonpatch.JsonPatchOperation{
		// ... 8 default patches
	}
)
```

**Fixture Location:**
- Fixtures defined at package level (top of test file)
- Reused across multiple test cases
- Separated into logical groups (e.g., `defaultPatches`, `defaultNetworkingPatches`)

## Mocking Strategy

**Mocking Framework:**
- No external mocking library used (no `mock`, `gomock`, `mockery`)
- Fake client from controller-runtime: `ctrlruntimefakeclient`

**Fake Client Pattern:**
```go
// From pkg/test/fake/builder.go
func NewClientBuilder() *ctrlruntimefakeclient.ClientBuilder {
	return ctrlruntimefakeclient.
		NewClientBuilder().
		WithScheme(NewScheme()).
		WithStatusSubresource(
			&appskubermaticv1.ApplicationInstallation{},
			&kubermaticv1.Addon{},
			// ... more types
		)
}
```

**Mocking Patterns:**
1. **Fake Kubernetes Client:** `ctrlruntimefakeclient` for object storage
2. **Getter Functions:** Provider-style getters passed as constructor parameters
3. **Stub Fields:** `disableProviderValidation` and `disableProviderMutation` flags in test hooks

**Provider Disabling in Tests:**
- Validation/mutation classes have boolean flags to disable cloud provider checks:
  ```go
  type validator struct {
      // ... other fields
      disableProviderValidation bool  // Only for unit tests
  }

  type Mutator struct {
      // ... other fields
      disableProviderMutation bool  // Only for unit tests
  }
  ```
- Prevents actual HTTP calls to cloud providers during testing
- Set explicitly in test setup when needed

**What to Mock:**
- Kubernetes API client (use fake client)
- Cloud provider credentials validation (disable with flags)
- HTTP-based external service calls (stub in tests)

**What NOT to Mock:**
- Business logic that's being tested
- Core object transformations and validations
- Kubernetes object semantics (use real K8s types)

## Special Test Patterns

**Diff Utilities for Comparison:**
From `pkg/test/diff/diff.go`:
```go
// Deep structural equality
func DeepEqual(expected, actual interface{}) bool

// Kubernetes semantic equality (ignores internal fields)
func SemanticallyEqual(expected, actual interface{}) bool

// Object diff in YAML format
func ObjectDiff(expected, actual interface{}) string

// String-based diff
func StringDiff(expected, actual string) string

// Set difference for comparable types
func SetDiff[T ordered](expected, actual sets.Set[T]) string
```

**Usage in Tests:**
```go
actualDiff := diff.ObjectDiff(expectedPatches, actualPatches)
if actualDiff != "" {
    t.Errorf("patches differ:\n%s", actualDiff)
}
```

**JSON Patch Assertions:**
Tests verify patch operations generated by mutation webhooks:
```go
wantPatches: []jsonpatch.JsonPatchOperation{
	jsonpatch.NewOperation("replace", "/spec/exposeStrategy", "Tunneling"),
	jsonpatch.NewOperation("add", "/spec/clusterNetwork/ipFamily", "IPv4"),
	// ... more operations
}
```

## Coverage

**Requirements:** Not enforced by default in code (no coverage gate visible)

**View Coverage:**
```bash
go test -v -coverprofile=coverage.out ./pkg/...
go tool cover -html=coverage.out
```

**Integration Tests:**
- Separate from unit tests
- Run via `make test-integration`
- Use integration build tag
- More extensive setup, may use real resources

**E2E Tests:**
- E2E test compilation verified without execution: `go test -tags "e2e,$(KUBERMATIC_EDITION)" -run nope ./pkg/test/e2e/...`
- Actual e2e tests invoked separately via CI/prow configuration

## Test Types & Scope

**Unit Tests:**
- Focus: Individual function behavior
- Scope: Single package or domain
- Example: `TestMutator` validates cluster mutation logic
- Example: `TestHandle` tests validation webhook request handling
- Isolation: Use fake clients, disable external calls
- Execution: `make run-tests` - all unit tests in 30m timeout

**Integration Tests:**
- Focus: Component interaction (e.g., webhook + client + seed/config)
- Scope: Multiple packages working together
- Location: Marked with `integration` build tag
- Execution: `make test-integration`

**E2E Tests:**
- Focus: End-to-end cluster lifecycle operations
- Scope: Full system including external resources
- Location: `pkg/test/e2e/` directory
- Build tag: `e2e`
- Execution: Separate from standard test runs (in CI pipeline)

## Common Testing Patterns

**Async Testing:**
- Context passed to functions: `ctx context.Context`
- No explicit goroutine testing patterns visible
- Timeout handled at `go test` level: `-timeout 30m`

**Error Testing:**
Pattern for testing error conditions:
```go
{
	name: "invalid cluster name format",
	cluster: rawClusterGen{
		Name: "INVALID!!!",  // Invalid name
		// ... other fields
	}.Build(),
	wantAllowed: false,  // Expect validation to fail
},
```

**Parametric Testing:**
- Achieved via table-driven tests
- Each table entry is a complete test scenario
- Uses `t.Run(test.name, func(t *testing.T) { ... })` for subtests

**Dependency Injection:**
- Constructor pattern: `NewValidator(client, seedGetter, configGetter, features, caBundle)`
- Dependencies passed at construction time
- No globals or singletons in test context

## Test File Examples

**Key Test Files:**
- `pkg/webhook/cluster/mutation/mutator_test.go` - Cluster mutation tests (237 lines analyzed)
- `pkg/webhook/cluster/validation/validation_test.go` - Validation webhook tests
- `pkg/webhook/seed/validation_test.go` - Seed validation tests
- `pkg/webhook/user/validation/validation_test.go` - User validation tests

**Test Naming Convention Examples:**
```
TestMutator - Main test function
├── "Create cluster sets default component settings"
├── "Create cluster sets default cni settings"
├── "Update cluster retains user-specified settings"
└── ... (more subtests)

TestHandle - Validation handler tests
├── "Delete cluster success"
├── "Create cluster with Tunneling expose strategy succeeds"
└── ... (more subtests)
```

## Build Tags for Testing

**Standard Tags:**
- `unit` - Default unit tests
- `integration` - Integration tests requiring component interaction
- `e2e` - End-to-end tests

**Edition Tags:**
- `ce` - Community Edition
- `ee` - Enterprise Edition
- `cloud` - Cloud builds

**Feature Tags:**
- `dualstack` - IPv4/IPv6 dual-stack tests
- `kubevirt` - KubeVirt provider tests
- `ipam` - IPAM functionality tests
- `logout` - Session logout tests
- `create` - Object creation tests

**Build Command Examples:**
```bash
# Unit tests only
go test -tags "unit,ce" ./pkg/...

# Integration tests
go test -tags "integration,ce" ./pkg/...

# E2E tests (compile-only check)
go test -tags "e2e,ce" -run nope ./pkg/test/e2e/...
```

---

*Testing analysis: 2026-03-31*
