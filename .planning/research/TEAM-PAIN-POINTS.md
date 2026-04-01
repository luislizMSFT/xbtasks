# Team Pain Points: ADO & Todo Management

Source: Team retrospective synthesis (monthly retros, recurring themes)

## Recurring Struggles

### 1. Last-minute ADO updates driven by external planning cycles
- ADO updates happen late and in bulk when planning/review deadlines arrive
- Planning inputs change late, pushing ADO hygiene to crunch time
- **Impact:** Incomplete fields, misaligned scope signals, higher cognitive load before reviews

### 2. ADO items drifting out of sync with reality
- Work changes shape mid-sprint but ADO updates lag behind
- Especially true for non-code work (planning, infra, cross-team dependencies)
- **Impact:** ADO becomes a reporting artifact, not a working system. Retros spend time re-explaining what actually happened.

### 3. To-dos live outside ADO longer than intended
- Short-term or "small" todos live in chat, heads, or files
- Only move into ADO when urgent or visible
- Root cause: creating/grooming ADO items feels too heavyweight for quick todos
- **Impact:** Lost or forgotten actions, repeated retro items ("did we actually do this?")

### 4. Unclear or inconsistent "done" criteria on ADO work
- Expectations for "done" aren't explicit, especially for cross-team or enabling work
- Confusion about whether "done" = code merged, or includes docs/validation/coordination
- **Impact:** Items linger in partially-done states, action items recur across retros

### 5. Retrospective action items don't consistently turn into owned ADO work
- Retros identify issues well, but conversion to durable ADO items with clear owners is uneven
- Same classes of issues reappear in future retros
- **Impact:** Retros feel insightful but less impactful over time

## Root Cause Pattern

| Pattern               | Root Cause                                 |
|-----------------------|--------------------------------------------|
| Late ADO cleanup      | External planning cadence, shifting inputs |
| Drift from reality    | Fast-moving work, slow hygiene             |
| Missing todos         | ADO friction for small/short-term work     |
| "Done" ambiguity      | Implicit vs explicit exit criteria         |
| Repeated retro themes | Weak retro-to-ADO execution loop           |

## Core Insight

The underlying issue isn't tooling or discipline — it's **mismatch between how work actually evolves and how ADO expects work to be shaped upfront**. The team adapts well in execution, but ADO updates and retros lag just behind that adaptation cycle.

## How This Informs the Tool

This tool should directly address these pain points:
- **Pain #3 (todos outside ADO):** Make task creation ultra-low-friction. Personal tasks are first-class, promotable to ADO when ready.
- **Pain #1 (last-minute ADO updates):** Dashboard shows what's out of sync, making continuous hygiene easier than bulk cleanup.
- **Pain #2 (drift from reality):** Bidirectional linking keeps local work and ADO aligned. Visual indicators show staleness.
- **Pain #4 (done ambiguity):** Task subtasks and checklist-style completion criteria make "done" explicit.
- **Pain #5 (retro→ADO loop):** Quick capture from any context → personal task → promote to ADO when ready.
