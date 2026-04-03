# Team ADO Tool — Project Reference

## Problem
Medium team (6-15 engineers) at Xbox Services needs better ADO work item hygiene. Current tools:
- Luis uses `xl` TUI (personal, won't be shared)
- Team lives in VS Code primarily
- ADO at `microsoft.visualstudio.com/Xbox` is the official tracker
- No good lightweight workflow for: local todos → refined → push to ADO

## Vision
A **team tool** that lets engineers manage local todos with low friction, sync to ADO when ready, and do code reviews — with multiple entry points (VS Code extension primary, TUI optional, desktop app aspirational).

## Existing Assets to Leverage

### xl (Go) — `/Users/luisliz/knowledge/xl/`
- `pkg/ado/client.go` — ADO REST API client (PAT + Azure CLI auth)
- `pkg/ado/sync.go` — Bidirectional status sync (local ↔ ADO)
- `pkg/db/` — SQLite abstraction with CRUD
- `cmd/ado.go` — CLI commands: `ado sync`, `ado status`, `ado link`
- Config: `XL_ADO_ORG`, `XL_ADO_PROJECT`, `XL_ADO_PAT`
- **ADO org**: `microsoft` / **project**: `Xbox`

### reviewme (TypeScript) — `/Users/luisliz/knowledge/reviewme/`
- VS Code extension for diff-based code review (Phase 1 POC complete)
- `GitService` — clean git abstraction (branches, diffs, file content)
- Git diff parsing (`parseNameStatus`, `parseNumStat`)
- Virtual document provider for side-by-side diffs
- Stack: TypeScript, VS Code API, `simple-git`, esbuild
- No ADO integration yet (planned Phase 5)
- No SQLite yet (planned Phase 3)

### Agency CLI + MCPs
- `agency copilot --mcp ado` — ADO MCP for natural language work item management
- Known MCPs: ado, enghub, es-chat, msft-learn, bluebird
- Team can use Copilot Chat with MCP as zero-UI interface

## Architecture Decision: The Language Question

**Reality:** Two working codebases in two languages. Both do what they're good at.

| Layer | Language | Rationale |
|-------|----------|-----------|
| ADO sync engine | **Go** | Already built in xl, proven, ADO REST client works |
| VS Code extension | **TypeScript** | Required by VS Code, reviewme code exists |
| Review/diff UI | **TypeScript** | reviewme already does this well |
| Desktop app | **TBD** | Aspirational — .NET MAUI considered for native feel |

**Key open question:** How Go and TypeScript connect:
1. **Go as CLI/MCP server** — TypeScript calls via subprocess/stdio (simplest)
2. **Rewrite ADO in TypeScript** — clean split, duplicates some work
3. **Go as local HTTP service** — TypeScript calls REST (cleanest interface, more parts)

## Workflow: What Team Members Experience

```
1. "Hey Copilot, what's on my sprint?"
   → MCP pulls from ADO, shows current items

2. Create local todo (VS Code sidebar, TUI, or Copilot chat)
   → Saved to local SQLite, no ADO overhead

3. "Push my local todos to ADO"
   → Syncs to ADO, creates work items, links to iteration

4. Pull from ADO → local view stays current

5. Code review in VS Code → reviewme handles diff rendering
```

## Proposed Local Schema (per-developer, lightweight)

```sql
todos(
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  status TEXT DEFAULT 'todo',     -- todo | in_progress | done | blocked
  priority TEXT,                   -- P0-P3
  assignee TEXT,
  ado_id INTEGER,                  -- NULL until pushed
  ado_url TEXT,
  ado_state TEXT,                  -- ADO work item state
  sync_status TEXT DEFAULT 'local_only',  -- local_only | synced | modified_locally | modified_remote | conflict
  area_path TEXT,
  iteration_path TEXT,
  tags TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now')),
  pushed_at TEXT
);
```

## Build Priority
1. **MCP server** — immediate value, team uses via Copilot with zero UI learning
2. **VS Code extension** — integrate todo sidebar + ADO sync into reviewme or new ext
3. **TUI** — power user option
4. **Desktop app** — aspirational, .NET MAUI for native feel on Mac/Windows

## Constraints
- No Electron
- No web-wrapper UIs (no Tauri web frontend)
- Must work on Mac and Windows
- ADO is source of truth for official items
- Local todos are staging area, not permanent store
- Team won't adopt a TUI as primary interface

## Related Files
- xl system design: `projects/xl-knowledge-system/system-design.md`
- reviewme architecture: `reviewme/.planning/`
- ADO config in memory: org=`microsoft`, project=`Xbox`
