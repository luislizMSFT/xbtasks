---
status: diagnosed
trigger: "synched is done but then prs appear if it's synched I expect everything to load. When I navigate to azure tab it loads again"
created: 2025-01-29T12:00:00Z
updated: 2025-01-29T12:00:00Z
---

## Current Focus

hypothesis: SyncCluster "Synced" label is driven solely by syncStore.syncing (backend sync lifecycle) which is completely decoupled from frontend data loading (PRs, pipelines, work items). Additionally, AdoView.vue unconditionally re-fetches all data on mount.
test: Code trace of SyncCluster label logic, sync store state management, App.vue prefetch flow, AdoView.vue mount behavior
expecting: Confirmed — two independent root causes
next_action: Return diagnosis

## Symptoms

expected: When SyncCluster shows "Synced", all data (including PRs) should already be loaded. Navigating to ADO tab should not trigger a re-fetch.
actual: "synched is done but then prs appear if it's synched I expect everything to load. When I navigate to azure tab it loads again"
errors: None reported
reproduction: Test 1 in UAT — authenticate, observe SyncCluster showing Synced while PRs are still loading, then navigate to ADO tab and observe it loading again
started: Discovered during UAT (phase 10)

## Eliminated

(none — first hypothesis confirmed)

## Evidence

- timestamp: 2025-01-29T12:00:00Z
  checked: SyncCluster.vue label logic (line 25)
  found: Label is `!adoStore.connected ? 'Offline' : syncStore.syncing ? 'Syncing…' : 'Synced'` — only uses syncStore.syncing and adoStore.connected
  implication: "Synced" means "not currently running a backend sync AND connected" — says nothing about whether frontend data (PRs, pipelines, work items) has loaded

- timestamp: 2025-01-29T12:00:00Z
  checked: sync.ts store — what sets syncing=true
  found: syncing is set true only by (1) sync:started backend Wails event (line 32) and (2) manualSync() call (line 70). It starts as false (line 12).
  implication: syncStore.syncing never reflects frontend data loading — it only tracks the Go SyncService goroutine lifecycle

- timestamp: 2025-01-29T12:00:00Z
  checked: App.vue auth watcher (lines 44-60)
  found: On auth, fires Promise.allSettled([taskStore.fetchTasks, projectStore.fetchProjects, prStore.fetchAll, adoStore.fetchWorkItemTree, adoStore.fetchLinkedAdoIds, adoStore.fetchPipelines, adoStore.fetchSavedQueries]). None of these set syncStore.syncing=true.
  implication: All 7 data fetches happen in parallel but syncStore never knows about them. SyncCluster shows "Synced" as soon as adoStore.connected becomes true (from any ADO API succeeding), even though 6 other fetches may still be in-flight.

- timestamp: 2025-01-29T12:00:00Z
  checked: adoStore.connected — what sets it to true
  found: Set true inside fetchWorkItems (line 80), fetchPipelines (line 119), fetchWorkItemTree (line 158), fetchAll (line 134). Any single successful ADO API call flips connected=true.
  implication: The first ADO call to succeed triggers "Synced" in SyncCluster, regardless of PR store or other stores still loading.

- timestamp: 2025-01-29T12:00:00Z
  checked: prStore (prs.ts) — relationship to syncStore
  found: PR store is completely independent. Has its own loading/connected/error refs. No integration with syncStore whatsoever.
  implication: PRs load independently; SyncCluster has zero visibility into PR loading state.

- timestamp: 2025-01-29T12:00:00Z
  checked: AdoView.vue onMounted (lines 77-84)
  found: Unconditionally calls adoStore.fetchWorkItemTree(), adoStore.fetchLinkedAdoIds(), adoStore.fetchSavedQueries(), prStore.fetchAll() — no guards checking if data is already present.
  implication: First navigation to ADO tab always re-fetches everything, even if App.vue prefetch already loaded the data. This sets loading=true in stores, causing loading UI to flash.

- timestamp: 2025-01-29T12:00:00Z
  checked: DashboardView.vue onMounted (lines 56-61) for comparison
  found: DashboardView guards its fetches: `if (!taskStore.tasks.length) taskStore.fetchTasks()` and `if (!prStore.myPRs.length && !prStore.reviewPRs.length) prStore.fetchAll()`. Comment says "Data is fetched by App.vue on auth — only fetch if stores are empty".
  implication: DashboardView already follows the correct pattern. AdoView does not — it always re-fetches.

- timestamp: 2025-01-29T12:00:00Z
  checked: App.vue router-view wrapper (lines 68-70)
  found: Uses `<keep-alive>` so onMounted only fires on first visit, not on re-navigation. Subsequent visits fire onActivated instead.
  implication: The ADO tab re-fetch is only on first visit, but since App.vue prefetches fire on auth (before any view mounts), the first visit to ADO still duplicates all fetches.

## Resolution

root_cause: |
  Two related issues:
  
  1. **SyncCluster "Synced" label is decoupled from actual data loading.** The label logic in SyncCluster.vue (line 25) uses `syncStore.syncing` to determine "Synced" vs "Syncing…". But `syncStore.syncing` only tracks the Go backend SyncService lifecycle (sync:started/sync:completed Wails events). It has zero awareness of frontend data fetches (PRs via prStore, pipelines via adoStore, work items via adoStore). When auth succeeds, App.vue fires 7 parallel fetches via Promise.allSettled — none of which set syncStore.syncing=true. As soon as any single ADO API call succeeds, adoStore.connected=true and SyncCluster shows "Synced" — while PRs and other data are still loading in the background. PRs then "pop in" after the fact.
  
  2. **AdoView.vue unconditionally re-fetches all data on mount.** Unlike DashboardView (which guards with `if (!prStore.myPRs.length)` checks), AdoView's onMounted fires adoStore.fetchWorkItemTree(), prStore.fetchAll(), etc. without checking if data already exists. This duplicates the App.vue prefetch and causes loading states to flash when first navigating to the ADO tab.
  
  The fundamental design gap: there is no unified "all data loaded" concept. SyncCluster conflates "backend sync process idle" with "all data ready", and individual stores load independently with no coordination.

fix: ""
verification: ""
files_changed: []
