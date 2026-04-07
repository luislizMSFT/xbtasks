# Phase 4: Work Item Lifecycle Tracking — Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-07
**Phase:** 04-work-item-lifecycle-tracking-prs-pipelines-task-traceability
**Areas discussed:** Completed PR expansion, PR lifecycle management, Graph view repurpose, PR→Task auto-linking, Completed PR retention, PR→ADO work item detection

---

## Completed PR Presentation

| Option | Description | Selected |
|--------|-------------|----------|
| Inline expansion | Completed PR row expands accordion-style | |
| Card layout | Completed PR becomes richer card with sections | |
| Separate "Completed" section | Dedicated area below active PRs | |
| Compact row + connected badges | PR becomes smaller with connected items (user's own description) | ✓ |

**User's choice:** "The PR becomes smaller with less info and has a connect to pipelines or items related to it"
**Notes:** Not a full card or accordion — a condensed trail with connected badges/chips

## Related Item Linking Method

| Option | Description | Selected |
|--------|-------------|----------|
| Auto-link by branch name match | PR branch matches pipeline source branch, task_id/ado_id links to work items | ✓ |
| Manual linking only | User explicitly connects PRs to tasks/pipelines | |
| Auto-link + manual override | Auto-detect from branch/IDs, user can add/remove | |

**User's choice:** Auto-link by branch name match
**Notes:** Already have branch matching and task_id/ado_id data

## Graph View Layout

| Option | Description | Selected |
|--------|-------------|----------|
| Force graph (current) + Timeline tab | Two tabs: relationship graph and horizontal timeline | ✓ |
| Force graph + Dependency chain tab | Graph for relationships, linear dependency chains | |
| Force graph + Timeline + Dependency chain | Three tabs | |

**User's choice:** Force graph + Timeline tab
**Notes:** User specified graph needs polish: "links are too small and thin and should be straighter", wants hierarchy "with top and bottom" instead of messy force graph

## Graph Node Entities

| Option | Description | Selected |
|--------|-------------|----------|
| All entity types | Tasks, ADO work items, PRs, and pipeline runs as nodes | ✓ |
| Task-centric | Tasks as primary, others as satellites | |
| Work item-centric | ADO work items as primary | |

**User's choice:** All entity types

## Graph Interaction

| Option | Description | Selected |
|--------|-------------|----------|
| Click to highlight chain | Click highlights full relationship chain, dims others | |
| Click to expand detail | Click opens side panel | |
| Both | Click highlights, double-click opens detail | |
| Hierarchy + highlight + detail | User's own description: "opens and highlights chain, hierarchy with top and bottom" | ✓ |

**User's choice:** Opens and highlights chain, hierarchical layout instead of messy graph
**Notes:** User explicitly wants hierarchical top-to-bottom layout, not force-directed

## Graph Layout Approach

| Option | Description | Selected |
|--------|-------------|----------|
| Hierarchical top-to-bottom | Work Items → Tasks → PRs → Pipelines | |
| Flexible hierarchy | User chooses top entity | |
| Agent's discretion | Agent decides exact layout | ✓ |

**User's choice:** You decide the exact layout approach

## PR Snooze Behavior

| Option | Description | Selected |
|--------|-------------|----------|
| Hide permanently | Removes until un-hidden | |
| Snooze with duration | Disappears for X time, reappears | ✓ |
| Both hide and snooze | Right-click offers both | |

**User's choice:** Snooze with duration

## Snooze Duration Options

| Option | Description | Selected |
|--------|-------------|----------|
| Quick presets only | 1h / 4h / 1d / 1w buttons | ✓ |
| Presets + custom | Quick buttons plus "Pick time" | |
| Smart snooze | Presets + auto-unsnooze on status change | |

**User's choice:** Quick presets only

## PR → Task Status Update

| Option | Description | Selected |
|--------|-------------|----------|
| Auto-transition | Auto-change to "in_review"/"done" | |
| Suggest status change | Toast nudge, no auto-change | ✓ |
| No auto-transition | Keep manual, just show PR status | |

**User's choice:** Suggest status change — nudge via toast
**Notes:** Consistent with Phase 2 "never auto-mutate" philosophy

## PR → ADO Work Item Detection

| Option | Description | Selected |
|--------|-------------|----------|
| Auto-detect and link | Parse AB#12345, auto-create connection | |
| Detect and suggest | Surface nudge, user confirms | ✓ |
| Don't auto-detect | Rely on existing task_id/ado_id fields | |

**User's choice:** Detect and suggest

## Completed PR Retention

| Option | Description | Selected |
|--------|-------------|----------|
| Time-based | Show 7 days, then archive | |
| Until acknowledged | Stay until dismissed | |
| Always visible | Stay forever | |
| Permanent local cache | Keep in SQLite forever, dashboard shows recent, timeline shows all | ✓ |

**User's choice:** "Keep forever in local db just don't query over and over... should be able to see a timeline of all work done"

## Timeline Layout

| Option | Description | Selected |
|--------|-------------|----------|
| Horizontal timeline | Left-to-right with time axis | ✓ |
| Vertical timeline | Top-to-bottom chronological | |
| Gantt-style | Horizontal bars showing duration | |

**User's choice:** Horizontal timeline

## Timeline Organization

| Option | Description | Selected |
|--------|-------------|----------|
| Per-task rows | Each task gets a lane | |
| All events on one stream | Single timeline mixing all entity events | ✓ |
| Grouped lanes | Separate lanes per entity type | |

**User's choice:** All events on one stream

## Agent's Discretion

- Exact hierarchical graph layout algorithm
- Dashboard section ordering for completed PRs
- Timeline visual design
- Snooze UI placement
- How many days of completed PRs show on dashboard

## Deferred Ideas

None
