# AdminGroupBinding — Map OIDC Groups to KKP Administrators

**Issue:** [kubermatic/kubermatic#14761](https://github.com/kubermatic/kubermatic/issues/14761) — _[feature-request] Allow mapping groups to be KKP Administrators_
**Milestone:** KKP 2.31 (due 2026-08-24) · **Labels:** customer-request, kind/feature, sig/api, sig/cluster-management, sig/ui
**Classification:** `feat(auth)` / `feat(admin)`

> Implementation status, phased work plan, verification, and open questions live in [admin-group-binding.progress.md](./admin-group-binding.progress.md).

---

## 1. What are we trying to achieve?

Let a KKP installation designate **OIDC / EntraID security groups** as administrators, so admin rights follow group membership in the identity provider automatically — granted when a member logs in, revoked when they leave the group — with no per-user edits.

## 2. What is the issue?

Today KKP admin is a **per-user boolean**: `User.Spec.IsAdmin` (`sdk/apis/kubermatic/v1/user.go:82`, JSON `admin`). Every admin must be flagged individually (dashboard Admin Panel → Administrators, or `kubectl edit user`). In dynamic organizations this means:

- every new teammate → manual promote; every departure → manual demote (churn, and a security gap if forgotten)
- the OIDC `groups` claim is **already stored** into `User.Spec.Groups` on every login (by the dashboard's `UserSaver` middleware) but currently **drives nothing** — the plumbing exists, unused

```mermaid
flowchart TD
    subgraph NOW [" NOW — manual per user "]
        A[Admin opens Settings → Administrators] --> B[Add alice@corp.com]
        B --> C["User CR: Spec.IsAdmin = true"]
        D[teammate joins] -.repeat manually.-> B
        E[member leaves] -.remove manually.-> B
    end
    subgraph AFTER [" AFTER — group driven "]
        F[Create AdminGroupBinding once<br/>group: platform-admins] --> G[controller reconciles]
        H[user logs in → Spec.Groups filled from OIDC token] --> G
        G --> I["Spec.IsAdmin set/cleared automatically"]
    end
```

## 3. What solution did we pick?

A new **cluster-scoped CRD `AdminGroupBinding`** (`kubermatic.k8c.io/v1`), mirroring the existing EE `GroupProjectBinding`, plus a **reconciler** that writes `User.Spec.IsAdmin`. Provenance is tracked with an **annotation on the User object** — no new field on the `User` CRD.

```yaml
apiVersion: kubermatic.k8c.io/v1
kind: AdminGroupBinding
metadata:
  name: platform-admins
spec:
  group: "platform-admins"   # exact OIDC group name, case-sensitive, required
  isAdmin: true              # default true
  isGlobalViewer: false      # optional; mirrors User.Spec.IsGlobalViewer
```

Provenance annotation (already defined, `sdk/apis/kubermatic/v1/admin_group_binding.go:33`):

```
admin.kubermatic.k8c.io/granted-by-group: <binding name>
```

**Why a CRD (not a `KubermaticSetting` field):**
- customer explicitly asked for a CRD ("admin groups should be a CRD, like admin users are")
- direct precedent: `GroupProjectBinding` already maps group → project role; this is the global-admin sibling
- a `globalsettings` field **leaks admin group names**: `/ws/admin/settings` streams full settings to every authenticated user

**Why an annotation (not a new `UserStatus.IsGroupAdmin` field):**
- zero API change to the `User` CRD → no deepcopy/CRD regen, no SDK re-vendor needed just for provenance, dashboard consumers keep working untouched
- carries **more** information than a bool: names *which* binding granted admin ("via group X" badge in UI)
- a new property would have to be handled in the dashboard repo too (SDK re-vendor, API surface, UI model) — annotation is readable through existing object metadata

## 4. How does it solve the issue?

Declare intent **once** in an `AdminGroupBinding`; thereafter membership is driven by the IdP token. The controller keeps `User.Spec.IsAdmin` — the single flag every existing consumer already reads — in sync with group membership:

```mermaid
flowchart TD
    L[OIDC login] -->|"UserSaver writes Spec.Groups (exists today, dashboard repo)"| U[User CR updated]
    B[AdminGroupBinding created/changed/deleted] --> R[admin-group-binding-controller]
    U -->|watch| R
    R -->|"EvaluateAdminGroups(bindings, user.Spec.Groups)"| DEC{groups match?}
    DEC -->|"match & not admin"| P["PROMOTE: Spec.IsAdmin=true<br/>+ annotation granted-by-group"]
    DEC -->|"no match & annotation present"| D["DEMOTE: Spec.IsAdmin=false<br/>remove annotation"]
    DEC -->|"annotation absent"| N["NEVER TOUCH<br/>(manual admin / first-user auto-admin / service accounts)"]
    P --> C[existing consumers unchanged:<br/>dashboard UserInfo & /me · alertmanager-authz · Grafana MLA]
    D --> C
```

Three-state reconcile logic:

| User state | Controller behavior |
|---|---|
| annotation present, still in a matching group | leave admin |
| annotation present, no longer in any matching group | **demote** + remove annotation |
| annotation **absent** | never touch — protects manual grants, first-user auto-admin, service accounts |

Because the controller writes `Spec.IsAdmin` (single source of truth), **no existing reader changes**: alertmanager-authorization-server (`cmd/alertmanager-authorization-server/main.go:234`), Grafana MLA (`pkg/controller/seed-controller-manager/mla/user_grafana_controller.go:305`, `org_grafana_controller.go:300`), dashboard `UserInfo` / `/me`.

## 5. Approaches

Both approaches share the same mechanism: a reconciler in this repo watches `AdminGroupBinding` + `User` and writes `User.Spec.IsAdmin` in both directions (promote and demote). The dashboard already writes `User.Spec.Groups` at login, so every login triggers a reconcile — the dashboard auth path never changes. What differs between the two approaches is how the system remembers **why** a user is an admin, which is what protects manually-granted admins from ever being demoted by group logic.

### Approach 1 — new `IsGroupAdmin` property on the `User` CRD

Add `UserStatus.IsGroupAdmin bool`. The reconciler sets it alongside `Spec.IsAdmin`; a user with `IsGroupAdmin: false` is never touched.

```mermaid
flowchart TD
    L[OIDC login] -->|UserSaver writes Spec.Groups| U[User CR]
    B[AdminGroupBinding changed] --> R[reconciler]
    U -->|watch| R
    R --> DEC{groups match binding?}
    DEC -->|match & not admin| P["Spec.IsAdmin = true<br/>Status.IsGroupAdmin = true"]
    DEC -->|no match & IsGroupAdmin=true| D["Spec.IsAdmin = false<br/>Status.IsGroupAdmin = false"]
    DEC -->|IsGroupAdmin=false| N["never touch<br/>(manual admin)"]
    P --> API[User CRD schema change:<br/>deepcopy + CRD YAML regen]
    API --> DASH["dashboard repo must follow:<br/>SDK re-vendor · API model · UI model"]
```

Cost: it is a `User` CRD schema change. Deepcopy and CRD YAML must be regenerated here, and the dashboard repo has to re-vendor the SDK and plumb the new property through its API and UI models before it can use it. The field is also just a boolean — it says *that* a group granted admin, not *which* group.

### Approach 2 — annotation on the `User` object (chosen)

No schema change. The reconciler stamps `admin.kubermatic.k8c.io/granted-by-group: <binding name>` (already defined, `sdk/apis/kubermatic/v1/admin_group_binding.go:33`) when promoting, and only demotes users carrying the annotation.

```mermaid
flowchart TD
    L[OIDC login] -->|UserSaver writes Spec.Groups| U[User CR]
    B[AdminGroupBinding changed] --> R[reconciler]
    U -->|watch| R
    R --> DEC{groups match binding?}
    DEC -->|match & not admin| P["Spec.IsAdmin = true<br/>+ annotation granted-by-group: &lt;binding&gt;"]
    DEC -->|no match & annotation present| D["Spec.IsAdmin = false<br/>remove annotation"]
    DEC -->|annotation absent| N["never touch<br/>(manual admin / first-user auto-admin / service accounts)"]
    P --> META["no User CRD change —<br/>annotation is plain object metadata"]
    META --> DASH["dashboard reads annotation as-is:<br/>'via group X' badge, no SDK/model change"]
```

The annotation needs no API change in either repo and carries more information: it names the granting binding, which is exactly what the UI needs for a "via group X" badge and what the controller checks before demoting. The trade-off is ergonomics — a typed bool is nicer in Go than a string lookup — but that is not worth a cross-repo schema change.

**Decision: Approach 2.** Zero `User` CRD change, richer provenance, dashboard keeps working untouched.

### Rejected earlier: config field instead of a CRD

The first draft put an `adminGroups` string list on the `KubermaticSetting` singleton (plus Approach 1's `IsGroupAdmin` field). Less code — the existing `GET/PATCH /api/v1/admin/settings` API would carry the field for free — but `/ws/admin/settings` streams the full settings object to every authenticated user, so admin group names would leak to non-admins. The customer also explicitly asked for a CRD, matching how `GroupProjectBinding` already works.

### Rejected earlier: promote-at-login

The dashboard's `UserSaver` middleware would flip `IsAdmin` during the login request itself, with a small reconciler only for demoting offline users. Promotion would be visible in the same login request instead of seconds later, but the admin-granting logic ends up split across two repos and two mechanisms, and `UserSaver` has a 1-minute write throttle that would need a bypass. The reconciler alone covers both directions, including offline users.

