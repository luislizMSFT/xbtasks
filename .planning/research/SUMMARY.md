# Research Summary: Team ADO Tool

## Key Findings

### Stack
**Next.js 15 + TypeScript + PostgreSQL + Entra ID** is the recommended stack. A unified TypeScript codebase (frontend + API) is more productive than a Go+TypeScript split, despite xl's existing Go ADO code. The ADO integration patterns from xl (state mapping, WIQL, push/pull) port cleanly to TypeScript. PostgreSQL replaces SQLite for multi-user concurrency. MSAL replaces PATs for team-grade auth.

### Table Stakes Features
1. Unified work item view (all ADO items in one place)
2. Personal task management (CRUD, priority, status — proven in xl)
3. Bidirectional task ↔ ADO linking (the core differentiator)
4. PR dashboard with comment threads
5. Microsoft SSO authentication
6. Search and filtering across all entities

### Architecture
BFF monolith pattern — one Next.js app serving UI + API. Hybrid ADO sync (webhooks for real-time + polling for reconciliation). PostgreSQL stores personal data + cached ADO data. Clear ownership model: ADO owns organizational fields, local owns personal fields.

### Critical Watch-Outs
1. **Scope creep** — The #1 risk. Don't rebuild ADO. Layer on top, deep-link back.
2. **Adoption** — Dogfood first. Day-one value: auto-import assigned items. Fewer steps than ADO for common tasks.
3. **ADO rate limits** — Cache aggressively, batch operations, use webhooks over polling.
4. **Sync conflicts** — Field-level ownership rules. ADO wins organizational fields, local wins personal.
5. **Auth complexity** — MSAL token refresh, multi-tenant edge cases, service principal for background jobs.

## Recommended Build Order

| Order | Component | Dependencies |
|-------|-----------|-------------|
| 1 | Project setup + Auth (Entra ID) | None |
| 2 | Personal tasks CRUD + database schema | Auth |
| 3 | ADO work item fetching + caching | Auth |
| 4 | Task ↔ ADO bidirectional linking + sync | Tasks + ADO items |
| 5 | PR dashboard + comment threads | Auth + ADO API |
| 6 | Search, filtering, polish | All data in DB |

## Decisions for Roadmap

| Decision | Recommendation | Confidence |
|----------|---------------|-----------|
| Frontend framework | Next.js 15 with App Router | High |
| Backend | Next.js API routes (monolith) | High |
| Database | PostgreSQL via Prisma | High |
| Auth | Entra ID + MSAL + NextAuth.js | High |
| ADO sync | Webhooks + polling hybrid | High |
| Deployment | Azure App Service | Medium |
| Real-time updates | SSE or polling (not WebSocket) | Medium |

---
*Synthesized: 2026-03-31*
*Sources: STACK.md, FEATURES.md, ARCHITECTURE.md, PITFALLS.md*
