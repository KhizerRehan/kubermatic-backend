# AdminGroupBinding — Implementation Progress & Work Plan

Proposal: [admin-group-binding.plan.md](./admin-group-binding.plan.md) · Issue: [kubermatic/kubermatic#14761](https://github.com/kubermatic/kubermatic/issues/14761)
Last updated: 2026-07-14

---

## 1. Current status (verified against working tree)

All on branch `feat/14761-admin-group-binding`, uncommitted:

| Item | File | Status |
|---|---|---|
| CRD Go types + `EvaluateAdminGroups` helper + annotation const | `sdk/apis/kubermatic/v1/admin_group_binding.go` | ✅ done |
| Unit test (table-driven: no/empty/matching/case-sensitive/OR) | `sdk/apis/kubermatic/v1/admin_group_binding_test.go` | ✅ done |
| Scheme registration | `sdk/apis/kubermatic/v1/register.go:110` | ✅ done |
| Deepcopy | `sdk/apis/kubermatic/v1/zz_generated.deepcopy.go` | ✅ done |
| Generated CRD YAML (location `master,seed`) | `pkg/crd/k8c.io/kubermatic.k8c.io_admingroupbindings.yaml` | ✅ done |
| codegen locationMap entry | `hack/update-codegen.sh:70` | ✅ done |
| Test fixture `GenAdminGroupBinding` | `pkg/test/generator/objects.go:310` | ✅ done |
| Reconciler | — | ⏳ |
| Validation webhook | — | ⏳ |
| Operator webhook wiring | — | ⏳ |
| RBAC (alertmanager chart, if needed) | — | ⏳ |
| Docs (EntraID GUID caveat) | — | ⏳ |
| Dashboard API + UI | separate repo | ⏳ (UI mockup only: `admin-administrators.png`) |

⚠️ **Drift to fix:** `pkg/crd/k8c.io/kubermatic.k8c.io_users.yaml` was hand-edited (reworded `groups` description) but the source comment in `sdk/apis/kubermatic/v1/user.go:88` was not — `hack/verify-codegen.sh` will fail. Either update the Go comment to match or revert the YAML hunk.

## 2. Codebases involved

Two repositories:

| Repo | Module | Changes |
|---|---|---|
| `kubermatic/kubermatic` | `sdk/` | `AdminGroupBinding` types, scheme registration, deepcopy, `EvaluateAdminGroups` helper |
| `kubermatic/kubermatic` | `pkg/crd/` | generated CRD YAML (`master,seed`) |
| `kubermatic/kubermatic` | `pkg/ee/` | reconciler (promote/demote + annotation) |
| `kubermatic/kubermatic` | `pkg/webhook/` + operator | validation webhook + `ValidatingWebhookConfiguration` wiring |
| `kubermatic/dashboard` | `modules/api` | v2 CRUD endpoints + provider; `GetAdmins` surfaces annotation; `SetAdmin` demotion guard |
| `kubermatic/dashboard` | `modules/web` | new "Admin Groups" page; "via group X" badge + disabled delete on Administrators page |

| Constraint | Detail |
|---|---|
| Sequencing | SDK/CRD must merge + tag in `kubermatic/kubermatic` **before** the dashboard can re-vendor (`make update-kkp`) — start both PRs early in the 2.31 cycle |
| Auth path | dashboard auth path untouched (reconciler owns all `IsAdmin` writes) — dashboard changes are UI/CRUD only |

## 3. Remaining work — file map (traced from `GroupProjectBinding` wiring)

`GroupProjectBinding` is the exact template; every remaining item below has a GPB counterpart to copy.

### Phase A — this repo

1. **Reconciler** `pkg/ee/admin-group-binding/controller/{controller,reconciler,doc}.go`
   — model on `pkg/ee/group-project-binding/controller/`. Watches `AdminGroupBinding` + `User`; calls `EvaluateAdminGroups`; promote/annotate, demote/un-annotate, never touch un-annotated users. Also drive `Spec.IsGlobalViewer` (respect the isAdmin/isGlobalViewer mutual exclusion in `pkg/validation/user.go:178`).
   — register: `cmd/master-controller-manager/wrappers_ee.go` (`Add(...)`) + no-op in `wrappers_ce.go`. Master-only (Users live on master; sync-controller for seeds only if a seed-side consumer emerges — GPB's sync-controller pattern available at `pkg/ee/group-project-binding/sync-controller/`).
   — fake-client tests modelled on `pkg/ee/group-project-binding/controller/reconciler_test.go`, matrix: promote / demote-with-annotation / untouched-without-annotation / first-user auto-admin untouched / service account never escalated / idempotent under conflict retry.
2. **Validation webhook** `pkg/webhook/admingroupbinding/validation/{validation,wrappers_ee,wrappers_ce}.go` + EE impl `pkg/ee/validation/admingroupbinding/`
   — model on `pkg/webhook/groupprojectbinding/validation/`. Rules: non-empty group (defense-in-depth over the CRD `MinLength=1`), trim/reject whitespace-only, reject duplicate binding for same group. CE wrapper behavior = open question 1 below.
   — register in `cmd/kubermatic-webhook/main.go` (~:234, next to GPB's `builder.WebhookManagedBy`).
3. **Operator webhook wiring**
   — const `AdminGroupBindingAdmissionWebhookName` in `pkg/controller/operator/common/resources.go` (~:92)
   — `AdminGroupBindingValidatingWebhookConfigurationReconciler` in `pkg/controller/operator/master/resources/kubermatic/webhooks.go` (model on GPB's at :312)
   — add to reconcile list `pkg/controller/operator/master/reconciler.go:745` and cleanup list `:248`
4. **CRD deployment:** nothing to do — `pkg/crd/fs.go` uses `//go:embed *`; the YAML + `kubermatic.k8c.io/location: master,seed` annotation are already in place. Installer (`pkg/install/stack/kubermatic-master/stack.go:330`) and seed operator (`pkg/controller/operator/seed/reconciler.go:359`) pick it up automatically.
5. **RBAC:** kubermatic-api SA already wildcards `kubermatic.k8c.io/*` (`pkg/controller/operator/master/resources/kubermatic/api.go:63`); operator SA is cluster-admin. Only add `admingroupbindings` to `charts/mla/alertmanager-proxy/templates/authzserver-clusterrole.yaml:20-33` **if** the authz server ever evaluates bindings directly (not needed with the reconciler approach — it reads `IsAdmin`).
6. **Docs:** EntraID caveat — `groups` claim emits **GUIDs** by default, not display names; matching is exact + case-sensitive.

### Phase B — dashboard `modules/api` (after SDK re-vendor via `make update-kkp`)

1. `AdminGroupBindingProvider` interface + EE impl (mirror group-project-binding provider)
2. New `pkg/handler/v2/admin-group-binding/` handlers + routes + swagger in `routes_v2.go`
3. `GetAdmins` (`pkg/provider/kubernetes/admin.go`): surface `grantedByGroup` from the annotation; `SetAdmin`: block demoting annotated users (must remove binding/group instead); guard against self-demotion via binding delete
4. **No `UserSaver` / auth-path change** (pure-reconciler approach)

### Phase C — dashboard `modules/web`

1. New "Admin Groups" page under Admin Panel → Users (clone `src/app/dynamic/enterprise/group/`), route + nav
2. Administrators page: "via group X" badge + disabled delete for annotation-derived rows (mockup: `admin-administrators.png`)

## 4. Verification

1. Fix the `users.yaml`/`user.go` drift → `hack/update-codegen.sh` then `hack/verify-codegen.sh` clean.
2. `cd sdk && go test ./apis/kubermatic/v1/ -run TestEvaluateAdminGroups` (passing today).
3. Reconciler fake-client test matrix (Phase A.1 above).
4. Webhook tests: empty/whitespace/duplicate group (model on `pkg/ee/validation/groupprojectbinding/groupprojectbinding_test.go`).
5. Manual e2e:
   - create binding for own group → log in → Admin Panel appears; `kubectl get users` ADMIN column true; annotation present
   - delete binding / leave group → demoted on next reconcile; annotation gone
   - manual admin also in mapped group → keeps admin after group removed
   - token without `groups` claim → no change
   - EntraID GUID-vs-display-name scenario

## 5. Edge cases

- **Offline users**: reconciler handles both directions — no login needed for demotion.
- **Self-lockout**: blocked at dashboard API; recovery always possible via `kubectl edit user`.
- **First-user auto-admin**: untouched (no annotation).
- **Explicit + group overlap**: user manually promoted *and* later annotated by group logic — v1 accepts that leaving the group demotes them (accept & document; granular source-tracking deferred).
- **Matching**: exact, case-sensitive (consistent with GroupProjectBinding).
- **Service accounts**: no OIDC login path → no `Spec.Groups` → never escalated.

## 6. Open questions (maintainer sign-off)

1. **CE or EE?** GroupProjectBinding is EE-only (CE webhook wrapper rejects with "disabled for the Community Edition"); plain admin management is CE. Plan is written against the EE template (safe default, whole template exists); flipping to CE = move controller out of `pkg/ee/`, drop the ce/ee wrapper pair.
2. Confirm CRD approach vs config field (customer asked for CRD; recommend confirming on the issue).
3. `isGlobalViewer` on the binding — currently modelled **yes** (types already ship it).
4. EntraID GUID group claims — document only (recommended), or resolve display names?

## 7. Acceptance criteria

- [ ] admin can create a group→admin mapping (CRD + dashboard UI)
- [ ] member of a mapped OIDC group gets KKP admin on login without per-user edit
- [ ] removing the binding or group membership revokes admin (including offline users)
- [ ] manually-granted admins unaffected by group logic
- [ ] group-derived admins visible + distinguishable ("via group X") in the Administrators list
- [ ] no leak of admin group names to non-admins
