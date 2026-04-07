# KKP Proposal: Migrate Kubernetes Dashboard to Headlamp

- **Issue:** [kubermatic/kubermatic#15287](https://github.com/kubermatic/kubermatic/issues/15287)
- **Upstream KEP:** [KEP-5008: Move Headlamp to SIG UI](https://github.com/kubernetes/enhancements/tree/master/keps/sig-ui/5008-headlamp)
- **Target Release:** KKP v2.31
- **Authors:** @KhizerRehan

---

## Table of Contents

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Background: Upstream Kubernetes Dashboard Retirement](#background-upstream-kubernetes-dashboard-retirement)
  - [What is Headlamp?](#what-is-headlamp)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [Design Decision: Seed-Side Deployment](#design-decision-seed-side-deployment)
  - [Alternative Considered: Application Framework](#alternative-considered-application-framework)
  - [Decision Rationale](#decision-rationale)
- [Architecture](#architecture)
  - [Current Architecture (Kubernetes Dashboard)](#current-architecture-kubernetes-dashboard)
  - [Proposed Architecture (Headlamp)](#proposed-architecture-headlamp)
  - [Architecture Comparison](#architecture-comparison)
  - [User Access Flow](#user-access-flow)
- [API Changes](#api-changes)
  - [New Types](#new-types)
  - [New Constants](#new-constants)
  - [Defaulting](#defaulting)
  - [Health Status](#health-status)
  - [Backward Compatibility](#backward-compatibility)
- [Migration Strategy](#migration-strategy)
  - [Phased Rollout](#phased-rollout)
  - [Parallel Coexistence](#parallel-coexistence)
  - [Upgrade Path for Existing Clusters](#upgrade-path-for-existing-clusters)
- [Implementation Details](#implementation-details)
  - [Phase 1: API Types and Constants](#phase-1-api-types-and-constants)
  - [Phase 2: Seed-Side Resources](#phase-2-seed-side-resources)
  - [Phase 3: User-Cluster Resources](#phase-3-user-cluster-resources)
  - [Phase 4: Seed Controller Manager Integration](#phase-4-seed-controller-manager-integration)
  - [Phase 5: User-Cluster Controller Manager Integration](#phase-5-user-cluster-controller-manager-integration)
  - [Phase 6: Image Registration](#phase-6-image-registration)
- [File Change Map](#file-change-map)
  - [New Files](#new-files)
  - [Modified Files](#modified-files)
  - [Auto-Generated Files](#auto-generated-files)
  - [Files NOT Changed](#files-not-changed)
- [RBAC Design](#rbac-design)
- [Risks and Mitigations](#risks-and-mitigations)
- [Pre-Implementation Requirements](#pre-implementation-requirements)
- [Test Plan](#test-plan)
  - [Unit Tests](#unit-tests)
  - [Integration Tests](#integration-tests)
  - [Manual Verification](#manual-verification)
- [Graduation Criteria](#graduation-criteria)
- [Production Readiness](#production-readiness)
  - [Observability](#observability)
  - [Air-Gapped Environments](#air-gapped-environments)
  - [Rollback Strategy](#rollback-strategy)
- [Validated Technical Details](#validated-technical-details)
  - [Server Flags](#server-flags)
  - [Container Details](#container-details)
  - [Live Cluster Test Results](#live-cluster-test-results)
  - [Known Issue: Read-Only Filesystem](#known-issue-read-only-filesystem)
- [Open Questions](#open-questions)
- [Future Work](#future-work)
- [References](#references)
<!-- /toc -->

---

## Summary

This proposal replaces the retired Kubernetes Dashboard with [Headlamp](https://headlamp.dev) in KKP. Headlamp is deployed as a **seed-side per-user-cluster deployment**, mirroring the exact pattern used by the current Kubernetes Dashboard. The old dashboard code remains intact during the transition period to allow parallel coexistence and safe rollback.

---

## Motivation

### Background: Upstream Kubernetes Dashboard Retirement

The Kubernetes Dashboard project (`github.com/kubernetes-retired/dashboard`) has been officially retired and archived by the Kubernetes community. The consequences for KKP are:

- **No security patches** — CVEs discovered in the dashboard will not be fixed
- **No bug fixes** — known issues remain unresolved
- **No Kubernetes compatibility** — newer Kubernetes API changes are not supported
- **Community risk** — continued use of archived software signals neglect to KKP users

The upstream Kubernetes community has chosen Headlamp as the successor project, formalized through [KEP-5008](https://github.com/kubernetes/enhancements/tree/master/keps/sig-ui/5008-headlamp) which moves Headlamp under SIG UI. The Headlamp repository now lives at `github.com/kubernetes-sigs/headlamp`.

### What is Headlamp?

Headlamp is a modern, extensible Kubernetes web UI. Key features relevant to KKP:

- **Active development** with dedicated maintainers and external contributors
- **Plugin system** — extensible via plugins (Cert Manager, KEDA, Flux, etc.)
- **OIDC support** — native authentication integration
- **Out-of-cluster mode** — connects to clusters via kubeconfig (critical for seed-side deployment)
- **Desktop and in-cluster modes** — flexible deployment options
- **Multi-cluster support** — single instance can manage multiple clusters
- **Built-in resource browsing** — no separate metrics-scraper component needed
- **Modern tech stack** — TypeScript, React, Material UI, Go backend

### Goals

1. Deploy Headlamp to every user cluster through KKP, matching current Kubernetes Dashboard functionality
2. Maintain the same security model (dashboard runs on seed, connects via kubeconfig)
3. Provide a smooth migration path with no downtime or breaking changes
4. Reduce user-cluster resource footprint (from 10 resources to 3)
5. Establish foundation for future plugin support (Cert Manager, KEDA, Flux)
6. Establish foundation for future OIDC integration

### Non-Goals

1. **NOT** removing the old Kubernetes Dashboard code in this scope (deferred to KKP v2.32+)
2. **NOT** implementing plugin support in this scope (follow-up work)
3. **NOT** implementing OIDC integration in this scope (follow-up work)
4. **NOT** forcing migration — both dashboards coexist during transition
5. **NOT** changing how users access the dashboard from the KKP UI
6. **NOT** deploying Headlamp via the Application Framework

---

## Proposal

### Design Decision: Seed-Side Deployment

Headlamp is deployed on the **seed cluster** in the per-user-cluster namespace (e.g., `cluster-xyz`), using a kubeconfig secret to connect to the user cluster API server. This mirrors the exact architecture of the current Kubernetes Dashboard.

```
Decision: SEED-SIDE GO RECONCILERS

Deploy Headlamp on the seed cluster using Go reconcilers,
the same proven pattern used by the current Kubernetes Dashboard.

+-- Same deployment location (seed cluster)
+-- Same access model (kubeconfig to user cluster)
+-- Same security boundary (credentials on seed)
+-- Same reconciler integration points
+-- Same health check pattern
```

### Alternative Considered: Application Framework

An earlier design (deprecated) proposed deploying Headlamp via the KKP Application Framework (Helm chart installed directly into user clusters):

```
REJECTED: Application Framework Approach

Seed Cluster                          User Cluster
+----------------------------+       +----------------------------+
| ApplicationDefinition      |       | headlamp namespace         |
|   name: headlamp           |       |   +-- Deployment (Helm)    |
|   method: helm             |       |   +-- Service (Helm)       |
|                            | Helm  |   +-- RBAC (Helm)          |
| ApplicationInstallation    |------>|   +-- OIDC config (Helm)   |
|   (one per user cluster)   |       |                            |
+----------------------------+       +----------------------------+
                                       ^-- Dashboard runs HERE
                                           (different from today)
```

### Decision Rationale

The Application Framework approach was rejected for the following reasons:

| Factor | Application Framework (Rejected) | Seed-Side Reconcilers (Chosen) |
|--------|----------------------------------|-------------------------------|
| **Architecture change** | Moves dashboard from seed to user cluster — fundamentally different deployment model | No architecture change — same deployment model as today |
| **Security model** | Dashboard credentials live in user cluster — weaker security boundary | Credentials stay on seed — same security boundary as today |
| **Radius** | Changes how every existing cluster's dashboard is accessed | Changes only what binary is deployed, not how or where |
| **Framework coupling** | Core platform component depends on optional Application Framework | No new dependencies — uses proven reconciler infrastructure |
| **Migration complexity** | Must migrate resources across two different cluster contexts | Swaps one seed deployment for another in the same namespace |
| **Rollback safety** | Complex rollback across seed and user cluster | Simple: delete headlamp deployment, k8s-dashboard still running |
| **CRD changes** | Requires immediate field rename (breaking change) | Adds new field alongside old one (non-breaking) |
| **Tested?** | Helm chart tested in isolation only | Full seed-side deployment validated against live dev cluster |

**The seed-side approach was chosen because it minimizes blast radius, preserves the security model, and provides the simplest migration path.** The only thing that changes is the container image and its flags — everything else (deployment location, access patterns, credential handling, health monitoring) stays the same.

---

## Architecture

### Current Architecture (Kubernetes Dashboard)

```
SEED CLUSTER (per-user-cluster namespace: cluster-xyz)
+--------------------------------------------------------------+
|                                                              |
|  Deployment: kubernetes-dashboard (2 replicas)               |
|    Image: kubernetesui/dashboard:v2.7.0                      |
|    Port: 9090                                                |
|    Command: /dashboard                                       |
|    Security: runAsUser 1001, readOnlyRootFilesystem          |
|                                                              |
|  Secret: kubernetes-dashboard-kubeconfig                      |
|    Cert user: kubermatic:kubernetes-dashboard                 |
|    Generated via GetInternalKubeconfigReconciler              |
|                                                              |
|  Controllers:                                                |
|    DeploymentReconciler   --> creates deployment             |
|    KubeconfigReconciler   --> creates kubeconfig secret      |
|    HealthCheck            --> monitors deployment readiness   |
+------------------------------+-------------------------------+
                               | kubeconfig (cert auth)
                               v
USER CLUSTER
+--------------------------------------------------------------+
|                                                              |
|  Namespace: kubernetes-dashboard (PSA baseline)              |
|    +-- dashboard-metrics-scraper Deployment (2 replicas)     |
|    |   Image: kubernetesui/metrics-scraper:v1.0.8            |
|    |   Port: 8000                                            |
|    +-- dashboard-metrics-scraper Service (ClusterIP)         |
|    +-- dashboard-metrics-scraper ServiceAccount              |
|    +-- system:kubernetes-dashboard Role                       |
|    +-- system:kubernetes-dashboard RoleBinding                |
|    +-- kubernetes-dashboard-key-holder Secret (JWE key)      |
|    +-- kubernetes-dashboard-csrf Secret                       |
|                                                              |
|  ClusterRole: system:dashboard-metrics-scraper               |
|  ClusterRoleBinding: system:dashboard-metrics-scraper        |
|                                                              |
|  Total: 10 resources across namespace + cluster scope        |
+--------------------------------------------------------------+
```

### Proposed Architecture (Headlamp)

```
SEED CLUSTER (per-user-cluster namespace: cluster-xyz)
+--------------------------------------------------------------+
|                                                              |
|  Deployment: headlamp (2 replicas)                           |
|    Image: ghcr.io/headlamp-k8s/headlamp:v0.26.0             |
|    Command: /headlamp/headlamp-server                        |
|    Args: -kubeconfig /etc/kubernetes/kubeconfig/kubeconfig    |
|          -html-static-dir /headlamp/frontend                 |
|    Port: 4466                                                |
|    Security: runAsUser 100, runAsGroup 101,                   |
|              readOnlyRootFilesystem                           |
|    Resources: 100m/128Mi request, 250m/256Mi limit           |
|    Volumes:                                                  |
|      - headlamp-kubeconfig (Secret, ro)                      |
|      - tmp-volume (EmptyDir, rw)                             |
|                                                              |
|  Secret: headlamp-kubeconfig                                 |
|    Cert user: kubermatic:headlamp                            |
|    Generated via GetInternalKubeconfigReconciler              |
|                                                              |
|  Controllers:                                                |
|    DeploymentReconciler   --> creates deployment             |
|    KubeconfigReconciler   --> creates kubeconfig secret      |
|    HealthCheck            --> monitors deployment readiness   |
+------------------------------+-------------------------------+
                               | kubeconfig (cert auth)
                               v
USER CLUSTER
+--------------------------------------------------------------+
|                                                              |
|  Namespace: headlamp (PSA baseline)                          |
|                                                              |
|  ClusterRole: system:headlamp                                |
|    Read access to cluster-scoped resources                   |
|    (namespaces, nodes, CRDs, workloads, RBAC, metrics)       |
|                                                              |
|  ClusterRoleBinding: system:headlamp                         |
|    Subject: User "kubermatic:headlamp"                       |
|    RoleRef: ClusterRole "system:headlamp"                    |
|                                                              |
|  Total: 3 resources (namespace + ClusterRole + binding)      |
|                                                              |
|  NO deployments, services, secrets, or service accounts      |
+--------------------------------------------------------------+
```

### Architecture Comparison

```
                 BEFORE                              AFTER
          (Kubernetes Dashboard)                  (Headlamp)

Seed:     kubernetes-dashboard deploy        headlamp deploy
          kubernetes-dashboard-kubeconfig     headlamp-kubeconfig
          --------------------------------   --------------------------
          2 resources                        2 resources  (same)

User:     kubernetes-dashboard namespace     headlamp namespace
          metrics-scraper deploy             ClusterRole
          metrics-scraper service            ClusterRoleBinding
          metrics-scraper SA                 --------------------------
          Role + RoleBinding                 3 resources  (70% reduction)
          ClusterRole + ClusterRoleBinding
          key-holder Secret
          csrf Secret
          --------------------------------
          10 resources

Go files: 2 seed + 11 user-cluster          2 seed + 4 user-cluster
          = 13 files total                   = 6 files total  (54% less)
```

| Aspect | Kubernetes Dashboard | Headlamp |
|--------|---------------------|----------|
| **Project status** | Archived (kubernetes-retired) | Active (kubernetes-sigs) |
| **Container image** | `kubernetesui/dashboard:v2.7.0` | `ghcr.io/headlamp-k8s/headlamp:v0.26.0` |
| **Seed deployment** | 2 replicas, port 9090 | 2 replicas, port 4466 |
| **User-cluster deployments** | metrics-scraper (2 replicas) | None |
| **User-cluster secrets** | JWE key-holder + CSRF token | None |
| **User-cluster services** | metrics-scraper ClusterIP | None |
| **User-cluster RBAC** | Role + ClusterRole (metrics only) | ClusterRole only (full browsing) |
| **User-cluster resource count** | 10 | 3 |
| **Plugin support** | Not possible | Native plugin system |
| **OIDC support** | Not implemented | Native support |
| **Health monitoring** | `ExtendedHealth.KubernetesDashboard` | `ExtendedHealth.Headlamp` |
| **Go files to maintain** | 13 | 6 |

**Why is Headlamp simpler on the user-cluster side?**

- **Built-in resource browsing** — Headlamp queries the Kubernetes metrics API directly. No separate metrics-scraper sidecar component needed.
- **Auth via kubeconfig** — Headlamp authenticates using the mounted kubeconfig secret from the seed cluster. No JWE key-holder or CSRF token secrets needed.
- **No in-cluster service** — Since Headlamp runs on the seed, there is no service to expose inside the user cluster.
- **Cluster-wide read access** — A single ClusterRole provides the read access Headlamp needs, replacing both the namespaced Role and the ClusterRole from the old dashboard.

### User Access Flow

```
End User (browser)
    |
    | KKP Dashboard UI exposes proxy URL
    | (same pattern as current kubernetes-dashboard)
    v
KKP API Server / Proxy
    |
    | routes to seed cluster, cluster-xyz namespace
    v
Headlamp Service (seed cluster)
    |
    | port 4466
    v
Headlamp Pod
    |
    | reads mounted kubeconfig secret
    | (cert user: kubermatic:headlamp)
    v
User Cluster API Server
    |
    | authorized via ClusterRoleBinding
    | (User "kubermatic:headlamp" -> ClusterRole "system:headlamp")
    v
Cluster Resources (namespaces, pods, deployments, nodes, etc.)
```

---

## API Changes

### New Types

**File:** `sdk/apis/kubermatic/v1/cluster.go`

```go
// Headlamp contains settings for the Headlamp component
// as part of the cluster control plane.
type Headlamp struct {
    // Controls whether Headlamp is deployed to the user cluster or not.
    // Enabled by default.
    Enabled bool `json:"enabled,omitempty"`
}
```

New field on `ClusterSpec`:

```go
// Headlamp holds the configuration for the Headlamp web UI component.
Headlamp *Headlamp `json:"headlamp,omitempty"`
```

New method on `ClusterSpec`:

```go
func (c ClusterSpec) IsHeadlampEnabled() bool {
    return c.Headlamp == nil || c.Headlamp.Enabled
}
```

**Design note:** `IsHeadlampEnabled()` returns `true` when the field is `nil` (not set). This means Headlamp is **enabled by default** for all clusters, matching the behavior of `IsKubernetesDashboardEnabled()`.

### New Constants

**File:** `pkg/resources/resources.go`

```go
HeadlampDeploymentName          = "headlamp"
HeadlampKubeconfigSecretName    = "headlamp-kubeconfig"
HeadlampCertUsername            = "kubermatic:headlamp"
HeadlampClusterRoleName         = "system:headlamp"
HeadlampClusterRoleBindingName  = "system:headlamp"
```

The namespace constant `"headlamp"` lives in the user-cluster package `constants.go`, matching the kubernetes-dashboard pattern where namespace constants are co-located with the user-cluster resource code.

### Defaulting

**File:** `pkg/defaulting/cluster.go`

```go
// Headlamp is enabled by default.
if spec.Headlamp == nil {
    spec.Headlamp = &kubermaticv1.Headlamp{
        Enabled: true,
    }
}
```

This is placed immediately after the existing `KubernetesDashboard` defaulting block (~line 145). Both defaulting blocks run independently — they do not interact.

### Health Status

**File:** `sdk/apis/kubermatic/v1/cluster_status.go`

```go
// Added to ExtendedClusterHealth:
Headlamp *HealthStatus `json:"headlamp,omitempty"`
```

Health is determined by checking the Headlamp Deployment's readiness in the seed cluster namespace, using the same `resources.HealthyDeployment()` utility used for the kubernetes-dashboard health check.

### Backward Compatibility

The existing `KubernetesDashboard` struct, field, and `IsKubernetesDashboardEnabled()` method are **NOT modified or removed** in this scope. Both fields coexist independently:

```
ClusterSpec:
  kubernetesDashboard:          <-- existing, unchanged
    enabled: true
  headlamp:                     <-- new, added alongside
    enabled: true
```

The old field will be deprecated and removed in a future release (KKP v2.32+) after the old dashboard code is fully removed.

---

## Migration Strategy

### Phased Rollout

```
KKP v2.30 (current)       KKP v2.31 (this proposal)     KKP v2.32+ (future)
+-----------------------+  +-----------------------+     +-----------------------+
|                       |  |                       |     |                       |
| kubernetes-dashboard  |  | kubernetes-dashboard  |     |                       |
| (deployed, active)    |  | (deployed, unchanged) |     | (REMOVED)             |
|                       |  |                       |     |                       |
|                       |  | + headlamp            |     | headlamp              |
|                       |  |   (NEW, deployed)     |     | (sole dashboard)      |
|                       |  |                       |     |                       |
| API:                  |  | API:                  |     | API:                  |
|   kubernetesDashboard |  |   kubernetesDashboard |     |   headlamp            |
|                       |  |   + headlamp (NEW)    |     |   (old field removed) |
+-----------------------+  +-----------------------+     +-----------------------+
```

### Parallel Coexistence

Because Headlamp uses completely different resource names, namespaces, and ports, both dashboards can run side-by-side with zero conflicts:

```
Seed cluster namespace (cluster-xyz):

  kubernetes-dashboard           (existing, port 9090)
  kubernetes-dashboard-kubeconfig (existing)
  headlamp                       (NEW, port 4466)
  headlamp-kubeconfig            (NEW)

User cluster:

  kubernetes-dashboard namespace  (existing)
    +-- all existing resources
  headlamp namespace              (NEW)
    +-- ClusterRole + binding only
```

Benefits of parallel coexistence:
- **Zero-downtime migration** — Headlamp is added, not replacing
- **Rollback safety** — if Headlamp has issues, kubernetes-dashboard is still there
- **Gradual confidence building** — test Headlamp in production alongside the working dashboard
- **No flag day** — frontend team can update at their own pace
- **Independent enable/disable** — each dashboard is controlled by its own API field

### Upgrade Path for Existing Clusters

When KKP is upgraded to v2.31, existing clusters are handled automatically:

```
KKP v2.31 Upgrade
       |
       v
  Defaulting Webhook runs for each Cluster object
       |
       v
  +-- spec.headlamp == nil? --+
  |                            |
 YES                          NO
  |                            |
  v                            v
  Set spec.headlamp =        Use as-is
  { enabled: true }          (user explicitly configured)
  (default for all clusters)
       |
       v
  Seed Controller reconciles
       |
       v
  Creates headlamp Deployment + kubeconfig Secret
  in cluster-xyz namespace on seed
       |
       v
  User-Cluster Controller reconciles
       |
       v
  Creates headlamp Namespace + ClusterRole + ClusterRoleBinding
  in user cluster
       |
       v
  Both dashboards now running in parallel
```

---

## Implementation Details

### Phase 1: API Types and Constants

**Goal:** Define the data model for Headlamp configuration.

| File | Change |
|------|--------|
| `sdk/apis/kubermatic/v1/cluster.go` | Add `Headlamp` struct, `Headlamp` field on `ClusterSpec`, `IsHeadlampEnabled()` method, `Headlamp` field on `ExtendedClusterHealth` |
| `pkg/resources/resources.go` | Add `HeadlampDeploymentName`, `HeadlampKubeconfigSecretName`, `HeadlampCertUsername`, `HeadlampClusterRoleName`, `HeadlampClusterRoleBindingName` |
| `pkg/defaulting/cluster.go` | Add defaulting block for `spec.Headlamp` (enabled=true) |

Then run deepcopy and CRD generation.

### Phase 2: Seed-Side Resources

**Goal:** Create the deployment and deletion reconcilers that manage Headlamp on the seed cluster.

**Create `pkg/resources/headlamp/deployment.go`:**

Mirrors `pkg/resources/kubernetes-dashboard/deployment.go`. Key differences:

```go
// headlampData interface (same pattern as kubernetesDashboardData)
type headlampData interface {
    Cluster() *kubermaticv1.Cluster
    RewriteImage(string) (string, error)
}

// DeploymentReconciler returns a NamedDeploymentReconcilerFactory
func DeploymentReconciler(data headlampData) reconciling.NamedDeploymentReconcilerFactory

// Key configuration:
//   Image:     ghcr.io/headlamp-k8s/headlamp:v0.26.0
//   Command:   ["/headlamp/headlamp-server"]
//   Args:      ["-kubeconfig", "/etc/kubernetes/kubeconfig/kubeconfig",
//               "-html-static-dir", "/headlamp/frontend"]
//   Port:      4466
//   Replicas:  2
//   RunAsUser: 100, RunAsGroup: 101
//   Resources: 100m/128Mi request, 250m/256Mi limit
//   Volumes:   headlamp-kubeconfig (Secret, ro), tmp-volume (EmptyDir, rw)
```

**Create `pkg/resources/headlamp/deletion.go`:**

```go
func ResourcesForDeletion(namespace string) []ctrlruntimeclient.Object {
    return []ctrlruntimeclient.Object{
        &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
            Name: resources.HeadlampDeploymentName, Namespace: namespace,
        }},
        &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
            Name: resources.HeadlampKubeconfigSecretName, Namespace: namespace,
        }},
    }
}
```

### Phase 3: User-Cluster Resources

**Goal:** Create the RBAC and namespace resources deployed in the user cluster.

**Create `pkg/controller/user-cluster-controller-manager/resources/resources/headlamp/`:**

| File | Purpose |
|------|---------|
| `constants.go` | `Namespace = "headlamp"`, `AppName = "headlamp"` |
| `namespace.go` | Creates `headlamp` namespace with PSA baseline labels |
| `clusterrole.go` | Creates `system:headlamp` ClusterRole with read-only access |
| `clusterrolebinding.go` | Binds ClusterRole to User `kubermatic:headlamp` |
| `deletion.go` | Returns list of resources for cleanup (ClusterRole, ClusterRoleBinding, Namespace) |

### Phase 4: Seed Controller Manager Integration

**Goal:** Wire Headlamp into the existing seed controller reconciliation loop.

**File: `pkg/controller/seed-controller-manager/kubernetes/resources.go`**

Three integration points, all following the exact pattern of the existing kubernetes-dashboard wiring:

```go
// 1. Deployment reconciler (~line 440, after kubernetes-dashboard block)
if data.Cluster().Spec.IsHeadlampEnabled() {
    deployments = append(deployments, headlamp.DeploymentReconciler(data))
}

// 2. Kubeconfig secret (~line 555, after kubernetes-dashboard kubeconfig)
if data.Cluster().Spec.IsHeadlampEnabled() {
    creators = append(creators,
        resources.GetInternalKubeconfigReconciler(
            namespace,
            resources.HeadlampKubeconfigSecretName,
            resources.HeadlampCertUsername,
            nil, data, r.log,
        ),
    )
}

// 3. Cleanup when disabled (~line 207, after kubernetes-dashboard cleanup)
if !cluster.Spec.IsHeadlampEnabled() {
    if err := r.ensureHeadlampResourcesAreRemoved(ctx, data); err != nil {
        return nil, err
    }
}
```

**File: `pkg/controller/seed-controller-manager/kubernetes/health.go`**

```go
// (~line 113, after kubernetes-dashboard health check)
if cluster.Spec.IsHeadlampEnabled() {
    status, err := r.headlampHealthCheck(ctx, cluster, ns)
    if err != nil {
        return nil, fmt.Errorf("failed to get headlamp health: %w", err)
    }
    extendedHealth.Headlamp = &status
}
```

### Phase 5: User-Cluster Controller Manager Integration

**Goal:** Wire Headlamp user-cluster resources into the existing reconciliation loop.

**File: `pkg/controller/user-cluster-controller-manager/resources/reconciler.go`**

```go
// Add to reconcileData struct
headlampEnabled bool

// Set during data init (~line 154)
data.headlampEnabled = cluster.Spec.IsHeadlampEnabled()

// Wire into reconcileNamespaces (~line 1014)
if data.headlampEnabled {
    creators = append(creators, headlamp.NamespaceReconciler)
}

// Wire into reconcileClusterRoles (~line 546)
if data.headlampEnabled {
    creators = append(creators, headlamp.ClusterRoleReconciler())
}

// Wire into reconcileClusterRoleBindings (~line 586)
if data.headlampEnabled {
    creators = append(creators, headlamp.ClusterRoleBindingReconciler())
}

// Cleanup when disabled (~line 297)
if !data.headlampEnabled {
    if err := r.ensureHeadlampResourcesAreRemoved(ctx); err != nil {
        return err
    }
}
```

**Methods NOT wired** (intentionally — Headlamp has no user-cluster deployments, services, secrets, roles, role bindings, or service accounts):

- `reconcileDeployments` — no in-cluster deployment
- `reconcileServices` — no in-cluster service
- `reconcileServiceAccounts` — no in-cluster service account
- `reconcileRoles` — using ClusterRole instead of namespaced Role
- `reconcileRoleBindings` — using ClusterRoleBinding instead
- `reconcileSecrets` — no JWE/CSRF secrets needed

### Phase 6: Image Registration

**File: `pkg/install/images/images.go`**

Add Headlamp image to the list of images collected for mirroring and air-gapped preloading.

**File: `hack/versions.yaml`**

Add Headlamp version entry for tracking.

---

## File Change Map

### New Files

```
pkg/resources/headlamp/
+-- deployment.go              Headlamp deployment reconciler (seed-side)
+-- deletion.go                Seed-side resource cleanup list

pkg/controller/user-cluster-controller-manager/resources/resources/headlamp/
+-- constants.go               Namespace and AppName constants
+-- namespace.go               Headlamp namespace with PSA baseline labels
+-- clusterrole.go             system:headlamp ClusterRole (read-only access)
+-- clusterrolebinding.go      Binds ClusterRole to User kubermatic:headlamp
+-- deletion.go                User-cluster resource cleanup list
```

### Modified Files

| File | Change |
|------|--------|
| `sdk/apis/kubermatic/v1/cluster.go` | Add `Headlamp` struct, field, `IsHeadlampEnabled()`, health field |
| `pkg/resources/resources.go` | Add 5 Headlamp constants |
| `pkg/defaulting/cluster.go` | Add Headlamp defaulting (enabled=true) |
| `pkg/controller/seed-controller-manager/kubernetes/resources.go` | Wire deployment reconciler, kubeconfig secret, cleanup |
| `pkg/controller/seed-controller-manager/kubernetes/health.go` | Add `headlampHealthCheck` |
| `pkg/controller/user-cluster-controller-manager/resources/reconciler.go` | Wire namespace, ClusterRole, ClusterRoleBinding, cleanup |
| `pkg/install/images/images.go` | Add Headlamp image to collection |
| `hack/versions.yaml` | Add Headlamp version entry |

### Auto-Generated Files

These are regenerated via `make generate`:

- `pkg/crd/k8c.io/kubermatic.k8c.io_clusters.yaml` — CRD schema includes `headlamp` field
- `pkg/crd/k8c.io/kubermatic.k8c.io_clustertemplates.yaml` — CRD schema includes `headlamp` field
- `sdk/apis/kubermatic/v1/zz_generated.deepcopy.go` — DeepCopy for `Headlamp` struct

### Files NOT Changed

The old Kubernetes Dashboard code is **intentionally left intact** for parallel coexistence:

- `pkg/resources/kubernetes-dashboard/` — seed-side deployment remains
- `pkg/controller/.../kubernetes-dashboard/` — all 11 user-cluster files remain
- `pkg/resources/test/fixtures/deployment-*-kubernetes-dashboard.yaml` — ~66 fixtures remain
- All existing kubernetes-dashboard wiring in `resources.go`, `health.go`, `reconciler.go` — unchanged

---

## RBAC Design

The ClusterRole `system:headlamp` grants **read-only access** to cluster resources that Headlamp needs to display:

```go
Rules: []rbacv1.PolicyRule{
    {
        APIGroups: []string{
            "",                            // core (pods, services, configmaps, etc.)
            "apps",                        // deployments, statefulsets, daemonsets
            "batch",                       // jobs, cronjobs
            "networking.k8s.io",           // ingresses, network policies
            "rbac.authorization.k8s.io",   // roles, bindings (RBAC visibility)
            "storage.k8s.io",              // storageclasses, PVs, PVCs
            "apiextensions.k8s.io",        // CRDs
            "policy",                      // poddisruptionbudgets
            "autoscaling",                 // HPAs
        },
        Resources: []string{"*"},
        Verbs:     []string{"get", "list", "watch"},
    },
    {
        APIGroups: []string{"metrics.k8s.io"},
        Resources: []string{"pods", "nodes"},
        Verbs:     []string{"get", "list", "watch"},
    },
}
```

**Security note:** The `rbac.authorization.k8s.io` group grants read access to roles and bindings. This is intentional — Headlamp's RBAC visibility feature allows users to see who has access to what. If more restrictive access is needed, this API group can be removed.

---

## Risks and Mitigations

| # | Risk | Impact | Likelihood | Mitigation |
|---|------|--------|------------|------------|
| 1 | ClusterRole grants broad read access to user-cluster resources | HIGH | CERTAIN | Review RBAC rules with security team; restrict API groups if needed. Read-only access only — no write/delete. |
| 2 | Headlamp container image UID/GID compatibility | MEDIUM | MEDIUM | Helm chart defaults to UID 100/GID 101. Validated that the container image supports these values. Test on actual clusters. |
| 3 | Plugin support not in initial scope (ticket requirement) | MEDIUM | CERTAIN | Documented as explicit follow-up work. Foundation is laid — Headlamp has a native plugin system. |
| 4 | OIDC integration not in initial scope (ticket requirement) | MEDIUM | CERTAIN | Documented as explicit follow-up work. Headlamp natively supports OIDC configuration. |
| 5 | Headlamp image not mirrored for air-gapped environments | HIGH | CERTAIN | Add `ghcr.io/headlamp-k8s/headlamp` to `quay.io/kubermatic-mirror` pipeline. Must be done before release. |
| 6 | Frontend (KKP Dashboard UI) needs updates for new health field | MEDIUM | CERTAIN | Coordinate with frontend team. Health field is additive — old field continues to work. |
| 7 | Read-only filesystem error on startup | LOW | CERTAIN | Known issue: Headlamp logs `mkdir /home/headlamp/.config: read-only file system` on startup. Does not affect functionality — Headlamp starts and serves requests normally. Mitigated by EmptyDir tmp-volume. |
| 8 | Both dashboards consume seed cluster resources | LOW | CERTAIN | Temporary during transition (v2.31 only). Resource requirements are modest (100m CPU, 128Mi memory per dashboard). |

---

## Pre-Implementation Requirements

| # | Requirement | Priority | Status |
|---|-------------|----------|--------|
| 1 | Verify headlamp container image works with UID 100/GID 101 and readOnlyRootFilesystem | BLOCKER | Validated on dev cluster |
| 2 | Add headlamp image to `quay.io/kubermatic-mirror` | HIGH | Not started |
| 3 | Define and review ClusterRole RBAC rules | HIGH | Draft complete (see RBAC Design section) |
| 4 | Frontend team coordination for new `headlamp` health status field | MEDIUM | Not started |

---

## Test Plan

### Unit Tests

- Golden test fixtures for Headlamp deployment reconciler (same pattern as `deployment-*-kubernetes-dashboard.yaml` fixtures)
- `go test ./pkg/resources/...` — verify deployment reconciler produces correct manifests
- `go test ./pkg/controller/...` — verify controller wiring
- `go test ./sdk/...` — verify API types, deepcopy, defaulting

### Integration Tests

- CRD generation: `make generate` completes without errors
- CRD schema: `kubermatic.k8c.io_clusters.yaml` includes `headlamp` field with correct structure
- DeepCopy: `zz_generated.deepcopy.go` includes `Headlamp` struct

### Manual Verification

```
Testing Flow

+-----------------+     +------------------+     +------------------+
| Unit Tests      |     | CRD Generation   |     | Create Cluster   |
| go test ./...   |---->| make generate    |---->| headlamp.enabled |
| Golden fixtures |     | Verify CRD YAML  |     | = true (default) |
+-----------------+     +------------------+     +------------------+
                                                        |
                                                        v
+-----------------+     +------------------+     +------------------+
| Coexistence     |     | Disable Test     |     | Verify Resources |
| Both dashboards |<----| Set enabled=false|<----| Seed: deploy+sec |
| No conflicts    |     | Verify cleanup   |     | User: CR+CRB+NS  |
+-----------------+     +------------------+     +------------------+
```

Checklist:

- [ ] Create new cluster with defaults -> Headlamp deployment appears in seed `cluster-xyz` namespace
- [ ] Headlamp kubeconfig secret is created in seed namespace
- [ ] ClusterRole `system:headlamp` exists in user cluster
- [ ] ClusterRoleBinding `system:headlamp` exists in user cluster
- [ ] `headlamp` namespace exists in user cluster with PSA baseline labels
- [ ] `ExtendedHealth.Headlamp` reports `HealthStatusUp`
- [ ] Headlamp UI is accessible via `kubectl port-forward deploy/headlamp 4466:4466 -n cluster-xyz`
- [ ] Headlamp UI can browse namespaces, pods, deployments, nodes
- [ ] Set `spec.headlamp.enabled: false` -> all Headlamp resources are deleted (seed and user cluster)
- [ ] Re-enable -> all resources are recreated
- [ ] kubernetes-dashboard continues to work alongside Headlamp (no interference)
- [ ] Headlamp image is available in mirrored registry for air-gap environments

---

## Graduation Criteria

### KKP v2.31 (This Proposal)

- Headlamp deployed alongside Kubernetes Dashboard (parallel coexistence)
- Enabled by default for all clusters
- Health monitoring via `ExtendedHealth.Headlamp`
- Can be disabled per-cluster via `spec.headlamp.enabled: false`

### KKP v2.32 (Future — Deprecation Phase)

- Mark `spec.kubernetesDashboard` field as deprecated in API docs
- Announce deprecation in release notes
- Frontend UI updated to show Headlamp instead of Kubernetes Dashboard
- Plugin support investigation and implementation (Cert Manager, KEDA, Flux)
- OIDC integration investigation and implementation

### KKP v2.33+ (Future — Removal Phase)

- Remove `spec.kubernetesDashboard` field from API
- Remove all Kubernetes Dashboard Go code (seed-side + user-cluster)
- Remove ~66 test fixture files
- Remove dashboard entries from `images.go` and `versions.yaml`
- Migration cleanup controllers removed (no more old resources to clean up)

---

## Production Readiness

### Observability

- **Health check:** `ExtendedHealth.Headlamp` status (Up/Down) visible in KKP API and eventually in KKP Dashboard UI
- **Pod monitoring:** Standard Kubernetes deployment monitoring (replicas, restarts, resource usage)
- **Logs:** Headlamp produces structured JSON logs to stdout

### Environments

Headlamp container image must be mirrored to `quay.io/kubermatic-mirror` before release. This is done via `pkg/install/images/images.go` which collects all images needed for air-gapped installations.

### Rollback Strategy

If Headlamp has issues after upgrade:

1. **Disable per-cluster:** `kubectl patch cluster <id> --type=merge -p '{"spec":{"headlamp":{"enabled":false}}}'`
2. **Kubernetes Dashboard still running** — unaffected by Headlamp addition
3. **No data migration** — Headlamp stores no persistent data
4. **No CRD migration** — old `kubernetesDashboard` field unchanged

---

## Validated Technical Details

### Server Flags

Validated against `ghcr.io/headlamp-k8s/headlamp:v0.26.0` on 2026-03-31:

| Flag | Default | KKP Usage | Notes |
|------|---------|-----------|-------|
| `-kubeconfig` | `""` | `/etc/kubernetes/kubeconfig/kubeconfig` | Path to mounted kubeconfig secret |
| `-in-cluster` | `false` | Not set (use default) | Must be false for seed-side deployment |
| `-port` | `4466` | Use default | Container port |
| `-plugins-dir` | platform-dependent | `/headlamp/plugins` | Static plugins directory |
| `-html-static-dir` | `""` | `/headlamp/frontend` | Frontend static files |
| `-base-url` | `""` | May set for routing | Investigate for KKP proxy |
| `-insecure-ssl` | `false` | May need for self-signed certs | Investigate per-cluster |

### Container Details

| Property | Value |
|----------|-------|
| Binary | `/headlamp/headlamp-server` |
| Image | `ghcr.io/headlamp-k8s/headlamp:v0.26.0` |
| Frontend static dir | `/headlamp/frontend` (built into container) |
| Plugins dir | `/headlamp/plugins` (built into container) |
| Default user | `headlamp` (UID 100, GID 101) |
| Default port | 4466 |

### Live Cluster Test Results

Tested 2026-03-31 on a KKP dev cluster with Headlamp deployed in a seed cluster namespace:

```
+-- Headlamp starts with -kubeconfig flag pointing to user cluster
+-- Proxy setup to external cluster API server works
+-- 34 namespaces, 12 nodes, deployments all browsable
+-- No authentication prompt (kubeconfig provides auth)
+-- Port 4466 serves UI correctly
+-- No -in-cluster flag needed (defaults to false)
+-- Pod logs confirm: "Kubeconfig path: /etc/kubernetes/kubeconfig/kubeconfig"
+-- Pod logs confirm: "Creating Headlamp handler"
```

### Known Issue: Read-Only Filesystem

On startup, Headlamp logs an error when trying to create its plugins config directory:

```json
{"level":"error","error":"mkdir /home/headlamp/.config: read-only file system","message":"creating plugins directory"}
```

This is a non-fatal error. Headlamp starts and serves requests normally. The `readOnlyRootFilesystem: true` security setting prevents writing to the container filesystem, but an `EmptyDir` volume is mounted at `/tmp` for any temporary files Headlamp needs. Plugin support is a follow-up concern that will address this properly.

Additionally, Headlamp logs a frontend file write error:

```json
{"level":"error","error":"open /headlamp/frontend/index.baseUrl.html: read-only file system","message":"writing file"}
```

This occurs because Headlamp tries to write a customized `index.html` with the base URL baked in. When `-base-url` is not set (our case), this is harmless — the default `index.html` works correctly.

---


## Future Work

| Release | Work Item | Description |
|---------|-----------|-------------|
| v2.32 | Remove old dashboard code | Delete `pkg/resources/kubernetes-dashboard/`, `pkg/controller/.../kubernetes-dashboard/`, ~66 test fixtures |
| v2.32 | Deprecate old API field | Mark `spec.kubernetesDashboard` as deprecated, update API docs |
| v2.32 | Frontend updates | Update KKP Dashboard UI to reference Headlamp health status |
| v2.32 | Plugin support | Investigate and implement Cert Manager, KEDA, Flux plugin configuration |
| v2.32 | OIDC integration | Integrate Headlamp OIDC with KKP's OIDC provider setup |
| v2.33 | Remove old API field | Remove `spec.kubernetesDashboard` from ClusterSpec |
| v2.33 | Remove migration cleanup | Remove `ensureKubernetesDashboardResourcesAreRemoved` controllers |

---

## References

- **Upstream issue:** [kubermatic/kubermatic#15287](https://github.com/kubermatic/kubermatic/issues/15287)
- **Upstream KEP:** [KEP-5008: Move Headlamp to SIG UI](https://github.com/kubernetes/enhancements/tree/master/keps/sig-ui/5008-headlamp)
- **Headlamp project:** [github.com/kubernetes-sigs/headlamp](https://github.com/kubernetes-sigs/headlamp)
- **Headlamp documentation:** [headlamp.dev/docs](https://headlamp.dev/docs/latest/)
- **Headlamp Helm chart:** `oci://ghcr.io/headlamp-k8s/charts/headlamp`
- **Retired Kubernetes Dashboard:** [github.com/kubernetes-retired/dashboard](https://github.com/kubernetes-retired/dashboard)
- **KKP Application Framework:** [docs.kubermatic.com/applications](https://docs.kubermatic.com/kubermatic/v2.26/architecture/concept/kkp-concepts/applications/)
