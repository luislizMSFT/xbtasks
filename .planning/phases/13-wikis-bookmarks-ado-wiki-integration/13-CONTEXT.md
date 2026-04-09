# Phase 13: Wikis, Bookmarks & ADO Wiki Integration - Context

**Gathered:** 2026-04-09
**Status:** Ready for planning

<domain>
## Phase Boundary

Add a personal wiki system (Obsidian/Notion hybrid — no Notion databases), project-level bookmarks, and ADO wiki browsing/importing with full-text search. Three features:

1. **Local Wikis** — editable markdown pages with sidebar navigation, live preview, and `[[wiki links]]`
2. **Project Bookmarks** — quick-link bookmark section on projects (reuse ExternalLinks pattern)
3. **ADO Wiki Browser + Import** — browse ADO repo wikis, selectively import pages as read-only mirrors, FTS5 search across all wiki content

</domain>

<decisions>
## Implementation Decisions

### Wiki page model
- Flat list of pages (no sub-pages / nesting) — each page has a title and markdown content with sections (headings)
- Pages are personal to the user — all stored locally in SQLite
- Expected scale: 10-30 pages over time — sidebar list is sufficient (no grid/search-first landing needed)
- No cross-references to tasks or projects for v1 — wikis are standalone knowledge
- Simple `[[Page Name]]` inter-page links for navigation — backlinks/graph visualization deferred to a later phase
- "No Notion databases" — pure wiki pages, no table views / kanban / galleries / relation properties

### Wiki layout
- Split layout: sidebar page list (left) + title at top + content area below (hybrid of Obsidian sidebar + Notion full-width)
- Markdown editing with live preview — toggle between edit mode and rendered view
- Full-view pages — sections (headings) always visible, no collapsing/accordion
- Search bar at top of sidebar — FTS5-powered, filters page list and shows preview snippets of matching content

### Content editing
- Markdown with live preview (toggle edit/rendered)
- Sections are just markdown headings (##, ###) — no special section model
- Standard markdown features: headings, bold, italic, code blocks, lists, links, images

### ADO Wiki integration
- Browse → Import pattern (same UX philosophy as ADO work item browser)
- Connect to configured ADO repos, browse their wiki page tree
- User selects individual pages or entire wikis to import
- Imported pages stored locally as read-only mirrors — clearly marked as "from ADO"
- Periodic re-fetch to keep mirrored content current
- ADO wiki is the first external source — architecture should allow adding more sources later
- Read-only always — imported ADO pages cannot be edited locally

### Project bookmarks
- ExternalLinks-style bookmark section on project detail pages
- Same URL pattern detection as task links (ICM, Grafana, ADO, wiki auto-detect)
- Quick save/remove bookmarks per project
- Reuse existing `ExternalLinks.vue` component + `externallinks.go` pattern — extend to work with project_id

### Full-text search
- FTS5 virtual table indexing all wiki content (local pages + imported ADO wiki pages)
- Search UI: inline search bar at top of wiki sidebar
- Results show page title + preview snippet with match context
- Click result to open the page
- Search scope: local wiki pages + mirrored ADO wiki content

### Claude's Discretion
- Exact markdown editor library choice (e.g., milkdown, tiptap, CodeMirror, or simple textarea + marked)
- FTS5 tokenizer and ranking configuration
- ADO wiki API pagination and rate limiting strategy
- Sync interval for ADO wiki re-fetch
- Sidebar width and responsive behavior
- Empty state design for new wiki pages
- Keyboard shortcuts for edit/preview toggle

</decisions>

<specifics>
## Specific Ideas

- "Like Obsidian + Notion but no databases" — clean wiki experience, not a project management tool within the wiki
- "Title at the top" — prominent page title, content flows below
- "Ownership areas" — some pages represent domains the user owns (e.g., "Service X Runbook"), others are shared topics
- "Evergreen, not tasks that die" — wiki content persists and is maintained over time, unlike tasks that complete and archive
- ADO wiki browse should feel like the existing ADO work item browser — familiar tree, select to import
- "Like VS Code" search — sidebar search with preview snippets, click to navigate
- Bookmarks are a project-level feature for quick links — not part of wiki pages

</specifics>

<canonical_refs>
## Canonical References

### External links (reusable pattern)
- `internal/app/externallinks.go` — URL pattern detection (ICM, Grafana, ADO, wiki), CRUD operations
- `frontend/src/components/tasks/ExternalLinks.vue` — Link management UI with type icons, add/delete
- `internal/db/links.go` — TaskLink DB operations (extend for project bookmarks)
- `domain/task.go:42-50` — TaskLink struct (reuse for project bookmarks + wiki links)

### ADO client (extend for wiki API)
- `pkg/ado/client.go` — REST client with Bearer token auth (extend with wiki endpoints)
- `pkg/ado/models.go` — ADO data models (add WikiPage model)
- `pkg/ado/factory.go` — Client factory with multi-org support

### Store patterns
- `frontend/src/stores/tasks.ts` — Pinia store pattern: ref state, computed, async actions
- `frontend/src/stores/ado.ts` — ADO data fetching + tree browsing pattern (reuse for wiki tree)
- `frontend/src/stores/projects.ts` — Project CRUD (extend with bookmarks)

### UI components (reusable)
- `frontend/src/components/tasks/CommentsSection.vue` — Comment/note editing pattern
- `frontend/src/components/ui/tabs/` — Tab switching (edit/preview toggle)
- `frontend/src/components/ui/card/` — Card layout for page containers
- `frontend/src/components/ui/scroll-area/` — Content scrolling
- `frontend/src/components/CommandPalette.vue` — Search UX reference
- `frontend/src/components/PageHeader.vue` — Consistent page headers

### Database schema
- `internal/db/migrate.go` — Existing migration pattern (extend with wiki + bookmark tables)

### Router
- `frontend/src/router/index.ts` — Route pattern with dynamic imports + meta tags

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `ExternalLinks.vue` + `externallinks.go`: Direct reuse for project bookmarks — just swap task_id for project_id
- `CommentsSection.vue`: Editing pattern (textarea + save) adaptable for wiki content editing
- ADO work item browser (`ADOBrowserView` pattern): Same browse → select → import UX for wiki pages
- shadcn/vue component library: 79 components already installed — Tabs, Card, Input, ScrollArea, Dialog
- `lucide-vue-next` icons: BookOpen, FileText, Search, Link, ExternalLink already available

### Established Patterns
- Pinia stores with `ref()` state + `computed()` derivations + async action functions
- Dynamic import routes with code-splitting
- Wails binding wrappers in `frontend/src/api/*.ts`
- Go services in `internal/app/` with SQLite via `internal/db/`
- URL pattern detection with regex in `externallinks.go`

### Integration Points
- New `/wikis` route in router alongside existing `/tasks`, `/projects`, `/ado`
- New `wikis` Pinia store following `tasks.ts` pattern
- New `wiki_pages` + `wiki_fts` tables in SQLite migration
- Extend `projects` with bookmarks (new `project_bookmarks` table or reuse `task_links` with polymorphic ID)
- Extend ADO client (`pkg/ado/`) with wiki API endpoints
- Sidebar navigation: add "Wikis" entry to existing nav

</code_context>

<deferred>
## Deferred Ideas

- **Backlinks / graph visualization** — `[[wiki links]]` exist for navigation, but automatic backlink tracking and graph view deferred to future phase
- **Cross-references to tasks/projects** — wiki pages referencing tasks inline, deferred to keep v1 simple
- **Additional external sources** — ADO wiki is first; architecture should support more but not implemented yet
- **Bidirectional wiki sync** — editing ADO wiki pages from the app; deferred (read-only import only)
- **Wiki page templates** — pre-filled templates for common page types (runbook, onboarding, etc.)

</deferred>

---

*Phase: 13-add-playgrounds-evergreen-ownership-areas-with-links-notes-bookmarks-and-ado-wiki-integration*
*Context gathered: 2026-04-09*
