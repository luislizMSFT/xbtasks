# Phase 13: Wikis, Bookmarks & ADO Wiki Integration - Research

**Researched:** 2025-07-22
**Domain:** Personal wiki system, project bookmarks, ADO wiki API integration, FTS5 full-text search
**Confidence:** HIGH

## Summary

This phase adds three features to the existing Wails v3 + Vue 3 + Go + SQLite desktop app: (1) a personal wiki with markdown editing and `[[wiki links]]`, (2) project-level bookmarks reusing the ExternalLinks pattern, and (3) an ADO wiki browser with selective import and FTS5 search. The codebase already has strong, reusable patterns for every aspect: `ExternalLinks.vue`/`externallinks.go` for bookmarks, `AdoTreeBrowser.vue`/`pkg/ado/client.go` for ADO browsing, Pinia stores with flight guards for state management, and `modernc.org/sqlite` (which includes FTS5 by default) for full-text search.

The markdown editing approach should use a **simple textarea + `marked` + `dompurify`** strategy rather than a heavyweight WYSIWYG editor. The project already has `dompurify` installed; adding `marked` (~40KB) keeps the bundle minimal and aligns with the "Vue is thin shell" philosophy. The edit/preview toggle maps directly to the existing shadcn-vue `Tabs` component. ADO wiki integration extends the existing `pkg/ado/client.go` with wiki-specific endpoints using the same Bearer token auth. FTS5 is a natural fit since `modernc.org/sqlite` compiles it in by default — external content tables with triggers keep the index in sync automatically.

**Primary recommendation:** Use textarea + marked + dompurify for editing (not a WYSIWYG editor), FTS5 external content tables with triggers for search, and extend the existing ADO client with wiki API endpoints following the same `doRequest`/`decodeResponse` pattern.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **Wiki page model**: Flat list (no sub-pages/nesting), personal to user, stored in SQLite, expected scale 10-30 pages
- **Wiki layout**: Split layout — sidebar page list (left) + title at top + content area below. FTS5-powered search bar at top of sidebar
- **Content editing**: Markdown with live preview (toggle edit/rendered). Sections are just markdown headings
- **`[[Page Name]]` links**: Simple inter-page links for navigation. Backlinks/graph visualization deferred
- **ADO Wiki integration**: Browse → Import pattern (same UX as ADO work item browser). Imported pages are read-only mirrors marked "from ADO". Periodic re-fetch. Architecture should allow more sources later
- **Project bookmarks**: ExternalLinks-style bookmark section on project detail pages. Reuse `ExternalLinks.vue` + `externallinks.go` pattern — extend to work with project_id
- **Full-text search**: FTS5 virtual table indexing all wiki content. Inline search bar at top of wiki sidebar. Results show page title + preview snippet
- **No Notion databases**: Pure wiki pages only — no table views, kanban, galleries, relation properties
- **Read-only imports**: ADO pages cannot be edited locally — always read-only
- **No cross-references to tasks/projects for v1**: Wikis are standalone knowledge

### Claude's Discretion
- Exact markdown editor library choice (e.g., milkdown, tiptap, CodeMirror, or simple textarea + marked)
- FTS5 tokenizer and ranking configuration
- ADO wiki API pagination and rate limiting strategy
- Sync interval for ADO wiki re-fetch
- Sidebar width and responsive behavior
- Empty state design for new wiki pages
- Keyboard shortcuts for edit/preview toggle

### Deferred Ideas (OUT OF SCOPE)
- Backlinks / graph visualization
- Cross-references to tasks/projects
- Additional external sources (beyond ADO wiki)
- Bidirectional wiki sync (editing ADO pages from the app)
- Wiki page templates
</user_constraints>

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| marked | 18.0.0 | Markdown → HTML rendering | Lightweight (~40KB), fast, well-maintained, supports custom extensions for `[[wiki links]]` |
| dompurify | 3.3.3 | HTML sanitization for rendered markdown | **Already installed** — used in ADODiscussion.vue |
| modernc.org/sqlite | 1.48.0 | SQLite with FTS5 support | **Already in go.mod** — FTS5 compiled in by default, no build tags needed |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| highlight.js | 11.11.1 | Code syntax highlighting in markdown preview | Optional — only if code blocks need highlighting in wiki pages |
| @vueuse/core | 14.2.1 (installed) | Debounce for search input, keyboard shortcuts | Already used throughout app (useDebounceFn, useMagicKeys) |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| textarea + marked | Milkdown 7.20 | Full WYSIWYG, but 500KB+ bundle, complex Vue 3 integration, overkill for 10-30 pages |
| textarea + marked | Tiptap 3.22 | Rich editor with Vue bindings, but 300KB+, needs prosemirror chain, heavy for simple wiki |
| textarea + marked | CodeMirror 6.41 | Excellent code editor, but designed for code not prose, poor markdown UX for non-developers |
| marked | markdown-it 14.1 | Similar capability, slightly more extensible plugin system, but marked has simpler API and smaller bundle |

**Decision: textarea + marked + dompurify.** Rationale:
1. The project philosophy is "Vue is thin shell" — all logic in Go, Vue for display only
2. The user expects edit/preview toggle (not inline WYSIWYG), which is exactly what textarea + rendered HTML provides
3. `dompurify` is already installed, `marked` adds only ~40KB
4. Custom `[[wiki links]]` are trivial to implement as a `marked` extension or pre-parse regex
5. 10-30 pages with simple markdown don't justify a 300-500KB editor framework

**Installation:**
```bash
cd frontend && npm install marked @types/marked
```

**No Go dependencies needed** — `modernc.org/sqlite` with FTS5 is already in go.mod.

## Architecture Patterns

### Recommended Project Structure
```
# Backend (Go)
domain/
├── wiki.go                    # WikiPage, WikiSearchResult, ADOWikiPage structs
internal/
├── app/
│   ├── wikiservice.go         # WikiService — CRUD for local wiki pages + search
│   ├── bookmarkservice.go     # BookmarkService — project bookmarks (extends ExternalLinks pattern)
│   └── adowikiservice.go      # ADOWikiService — browse/import/refresh ADO wiki pages
├── db/
│   ├── wiki.go                # Wiki DB operations (pages + FTS5 queries)
│   └── bookmarks.go           # Bookmark DB operations
│   └── migrate.go             # (extend with wiki_pages, wiki_fts, project_bookmarks tables)
pkg/
├── ado/
│   ├── wiki.go                # ADO Wiki REST API functions (ListWikis, GetPageTree, GetPageContent)
│   └── models.go              # (extend with WikiInfo, WikiPage models)

# Frontend (Vue)
frontend/src/
├── api/
│   ├── wiki.ts                # Wails binding wrappers for WikiService
│   ├── bookmarks.ts           # Wails binding wrappers for BookmarkService
│   └── ado-wiki.ts            # Wails binding wrappers for ADOWikiService
├── stores/
│   └── wiki.ts                # Pinia store — wiki pages, selected page, search, ADO wiki tree
├── views/
│   └── WikiView.vue           # Main wiki view (split layout: sidebar + content)
├── components/
│   └── wiki/
│       ├── WikiSidebar.vue    # Page list + search bar + ADO wiki import trigger
│       ├── WikiEditor.vue     # Textarea for markdown editing
│       ├── WikiPreview.vue    # Rendered markdown with wiki link handling
│       ├── WikiPageHeader.vue # Title + edit/preview toggle + metadata
│       └── ADOWikiBrowser.vue # Browse ADO wikis, select pages to import (reuse AdoTreeBrowser pattern)
├── composables/
│   └── useMarkdown.ts         # marked configuration + wiki link extension + DOMPurify sanitization
├── types/
│   └── index.ts               # (extend with WikiPage, WikiSearchResult, ADOWikiPage, ProjectBookmark)
```

### Pattern 1: Wiki Service (Go Backend)
**What:** Go service bound to Wails frontend for wiki CRUD + FTS5 search
**When to use:** All wiki operations go through this service
**Example:**
```go
// internal/app/wikiservice.go — follows existing service patterns
type WikiService struct {
    db *db.DB
}

func NewWikiService(database *db.DB) *WikiService {
    return &WikiService{db: database}
}

func (s *WikiService) CreatePage(title, content string) (domain.WikiPage, error) {
    // Insert into wiki_pages → FTS5 triggers auto-update wiki_fts
    res, err := s.db.Exec(
        `INSERT INTO wiki_pages (title, content) VALUES (?, ?)`,
        title, content,
    )
    // ...
}

func (s *WikiService) Search(query string) ([]domain.WikiSearchResult, error) {
    rows, err := s.db.Query(`
        SELECT wp.id, wp.title, wp.source, wp.is_readonly,
               snippet(wiki_fts, 1, '<mark>', '</mark>', '...', 32) as preview
        FROM wiki_fts
        JOIN wiki_pages wp ON wp.id = wiki_fts.rowid
        WHERE wiki_fts MATCH ?
        ORDER BY rank
        LIMIT 20`, query)
    // ...
}
```

### Pattern 2: Pinia Store (Frontend)
**What:** Wiki store following existing tasks.ts pattern with ref state, computed, async actions
**When to use:** All wiki state management
**Example:**
```typescript
// frontend/src/stores/wiki.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import type { WikiPage, WikiSearchResult } from '@/types'
import * as wikiApi from '@/api/wiki'

export const useWikiStore = defineStore('wiki', () => {
  const pages = ref<WikiPage[]>([])
  const selectedPageId = ref<number | null>(null)
  const loading = ref(false)
  const searchQuery = ref('')
  const searchResults = ref<WikiSearchResult[]>([])
  const editMode = ref(false)

  let fetchInFlight = false  // flight guard pattern from tasks.ts

  const selectedPage = computed(() =>
    pages.value.find(p => p.id === selectedPageId.value) ?? null
  )

  const sortedPages = computed(() =>
    [...pages.value].sort((a, b) => a.title.localeCompare(b.title))
  )

  // Debounced search (200ms, matching ADO tree search pattern)
  const doSearch = useDebounceFn(async (q: string) => {
    if (!q.trim()) { searchResults.value = []; return }
    searchResults.value = (await wikiApi.search(q)) as WikiSearchResult[]
  }, 200)

  async function fetchPages() {
    if (fetchInFlight) return
    fetchInFlight = true
    loading.value = true
    try {
      pages.value = (await wikiApi.listPages()) as WikiPage[]
    } catch {
      pages.value = []
    } finally {
      loading.value = false
      fetchInFlight = false
    }
  }

  // ... createPage, updatePage, deletePage, selectPage, toggleEditMode
  return { pages, selectedPageId, loading, searchQuery, searchResults, editMode, selectedPage, sortedPages, fetchPages, doSearch }
})
```

### Pattern 3: Markdown Rendering with Wiki Links
**What:** Composable that configures `marked` with custom `[[wiki link]]` extension and DOMPurify sanitization
**When to use:** WikiPreview.vue for rendering, also for search snippet display
**Example:**
```typescript
// frontend/src/composables/useMarkdown.ts
import { marked } from 'marked'
import DOMPurify from 'dompurify'

// Pre-parse: convert [[Page Name]] to markdown links before marked processes
function resolveWikiLinks(markdown: string, existingPages: Set<string>): string {
  return markdown.replace(/\[\[([^\]]+)\]\]/g, (match, pageName) => {
    const exists = existingPages.has(pageName.trim())
    const cssClass = exists ? 'wiki-link' : 'wiki-link wiki-link-missing'
    return `<a href="#" class="${cssClass}" data-wiki-page="${pageName.trim()}">${pageName.trim()}</a>`
  })
}

export function useMarkdown() {
  function render(content: string, existingPages: Set<string>): string {
    const withLinks = resolveWikiLinks(content, existingPages)
    const html = marked.parse(withLinks, { async: false }) as string
    return DOMPurify.sanitize(html, {
      ALLOWED_TAGS: ['h1','h2','h3','h4','h5','h6','p','a','strong','em','code','pre','ul','ol','li','blockquote','img','mark','br','hr','table','thead','tbody','tr','th','td'],
      ALLOWED_ATTR: ['href','src','alt','class','data-wiki-page'],
    })
  }

  return { render }
}
```

### Pattern 4: ADO Wiki API Extension
**What:** Extend `pkg/ado/client.go` pattern with wiki-specific API calls
**When to use:** Browsing and importing ADO wiki pages
**Example:**
```go
// pkg/ado/wiki.go
package ado

import "fmt"

// WikiInfo represents an ADO wiki.
type WikiInfo struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Type string `json:"type"` // "projectWiki" or "codeWiki"
    URL  string `json:"url"`
}

// WikiPageInfo represents a page in an ADO wiki.
type WikiPageInfo struct {
    ID          int             `json:"id"`
    Path        string          `json:"path"`
    Order       int             `json:"order"`
    IsParent    bool            `json:"isParentPage"`
    GitItemPath string          `json:"gitItemPath"`
    Content     string          `json:"content"`
    SubPages    []WikiPageInfo  `json:"subPages"`
    URL         string          `json:"url"`
}

// ListWikis fetches all wikis in the configured project.
func ListWikis(c *Client) ([]WikiInfo, error) {
    url := c.apiURL("wiki/wikis?api-version=7.1")
    resp, err := c.doRequest("GET", url, nil, "")
    if err != nil {
        return nil, fmt.Errorf("list wikis: %w", err)
    }
    var result struct {
        Value []WikiInfo `json:"value"`
    }
    if err := decodeResponse(resp, &result); err != nil {
        return nil, fmt.Errorf("parse wikis response: %w", err)
    }
    return result.Value, nil
}

// GetWikiPageTree fetches the full page tree for a wiki.
func GetWikiPageTree(c *Client, wikiID string) (*WikiPageInfo, error) {
    url := c.apiURL(fmt.Sprintf("wiki/wikis/%s/pages?path=/&recursionLevel=full&api-version=7.1", wikiID))
    resp, err := c.doRequest("GET", url, nil, "")
    if err != nil {
        return nil, fmt.Errorf("get wiki page tree: %w", err)
    }
    var page WikiPageInfo
    if err := decodeResponse(resp, &page); err != nil {
        return nil, fmt.Errorf("parse wiki page tree: %w", err)
    }
    return &page, nil
}

// GetWikiPageContent fetches the content of a specific wiki page.
func GetWikiPageContent(c *Client, wikiID, pagePath string) (*WikiPageInfo, error) {
    url := c.apiURL(fmt.Sprintf("wiki/wikis/%s/pages?path=%s&includeContent=true&api-version=7.1",
        wikiID, pagePath))
    resp, err := c.doRequest("GET", url, nil, "")
    if err != nil {
        return nil, fmt.Errorf("get wiki page content: %w", err)
    }
    var page WikiPageInfo
    if err := decodeResponse(resp, &page); err != nil {
        return nil, fmt.Errorf("parse wiki page content: %w", err)
    }
    return &page, nil
}
```

### Pattern 5: Project Bookmarks (Extending ExternalLinks)
**What:** Reuse the ExternalLinks pattern with project_id instead of task_id
**When to use:** Project detail page bookmark section
**Example:**
```go
// internal/app/bookmarkservice.go
type BookmarkService struct {
    db *db.DB
}

func (s *BookmarkService) AddBookmark(projectID int, url, label string) (domain.ProjectBookmark, error) {
    linkType := DetectLinkType(url)  // Reuse existing URL pattern detection
    return s.db.CreateBookmark(projectID, url, label, linkType)
}
```

### Anti-Patterns to Avoid
- **Don't embed a WYSIWYG editor**: Milkdown/Tiptap add 300-500KB and complex prosemirror state management. The user wants edit/preview toggle, not inline editing.
- **Don't use `content` FTS5 tables (non-external)**: External content tables avoid data duplication. Use `content=wiki_pages, content_rowid=id` to point FTS5 at the real table.
- **Don't fetch all ADO wiki page content at once**: Fetch the tree structure first (lightweight), then fetch content only for pages the user selects to import.
- **Don't create a separate search page/view**: Search is inline in the wiki sidebar — filter the page list and show snippets directly.
- **Don't hand-roll markdown rendering**: Use `marked` — it handles edge cases in CommonMark spec that are extremely hard to get right.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Markdown → HTML | Custom parser | `marked` library | CommonMark spec has 600+ edge cases; marked handles them all |
| HTML sanitization | Regex stripping | `dompurify` (already installed) | XSS prevention is security-critical; custom sanitizers always miss edge cases |
| Full-text search | LIKE queries or JS-side filtering | SQLite FTS5 | FTS5 handles tokenization, stemming, ranking, snippets — orders of magnitude faster than LIKE |
| URL type detection | New detection logic | Existing `DetectLinkType()` in `externallinks.go` | Already handles ICM, Grafana, ADO, Wiki patterns with regex |
| Wiki link syntax | Custom markdown processor | Pre-parse regex + marked | Simple `[[...]]` regex before marked processes; don't modify marked internals |
| ADO REST client | New HTTP client | Extend existing `pkg/ado/client.go` | Already has Bearer auth, timeout, JSON parsing, error handling |

**Key insight:** Every building block already exists in the codebase. This phase is primarily about composing existing patterns (ExternalLinks → Bookmarks, ADO tree browser → Wiki browser, SQLite → FTS5) with one new addition (marked for markdown rendering).

## Common Pitfalls

### Pitfall 1: FTS5 Sync Triggers Not Handling Updates Correctly
**What goes wrong:** FTS5 external content tables require both a DELETE of the old content and INSERT of the new content on UPDATE. Missing either causes stale or duplicate search results.
**Why it happens:** SQLite FTS5 doesn't auto-track external content changes — you must manually maintain sync via triggers.
**How to avoid:** Use BEFORE UPDATE (delete old) + AFTER UPDATE (insert new) trigger pair. Always test search after edits.
**Warning signs:** Search finds old content after editing a page, or finds the same page twice.

### Pitfall 2: Wiki Link Rendering Before Marked Processing
**What goes wrong:** If `[[Page Name]]` syntax is passed to `marked` first, it may be interpreted as nested bracket links or interfere with standard markdown link parsing.
**Why it happens:** Marked processes `[text](url)` links, and `[[text]]` can be partially matched.
**How to avoid:** Pre-process wiki links BEFORE passing to `marked`. Replace `[[Page Name]]` with HTML `<a>` tags, then let marked process the rest. DOMPurify allows `data-wiki-page` attribute for click handling.
**Warning signs:** Wiki links rendered as broken markdown or with extra brackets.

### Pitfall 3: ADO Wiki Page Paths Need URL Encoding
**What goes wrong:** ADO wiki page paths contain spaces and special characters (e.g., `/My Page/Sub Page`). Passing them raw to the API URL causes 404 errors.
**Why it happens:** The REST API expects URL-encoded paths in the query parameter.
**How to avoid:** Use `url.QueryEscape()` or `url.PathEscape()` on the path parameter before constructing the API URL.
**Warning signs:** ADO wiki pages with spaces in their names fail to load.

### Pitfall 4: FTS5 MATCH Query Syntax Errors
**What goes wrong:** User search input may contain FTS5 special characters (`*`, `"`, `AND`, `OR`, `NOT`, `NEAR`) that cause SQLite errors.
**Why it happens:** FTS5 MATCH uses its own query syntax. Raw user input can be invalid syntax.
**How to avoid:** Sanitize search input — wrap each word in double quotes, or use the `{column}:{term}` syntax. Simplest: quote the entire query as a phrase search.
**Warning signs:** Search crashes on certain inputs with "fts5: syntax error" from SQLite.

### Pitfall 5: Read-Only Flag Not Enforced at Backend
**What goes wrong:** Imported ADO pages should be read-only, but if the backend doesn't enforce this, the frontend edit button could accidentally modify ADO-sourced content.
**Why it happens:** Frontend-only guards can be bypassed; backend should reject updates to read-only pages.
**How to avoid:** Add `is_readonly` column to wiki_pages table. Backend `UpdatePage()` checks this flag and returns an error for read-only pages.
**Warning signs:** ADO-imported page content gets overwritten by local edits.

### Pitfall 6: Large ADO Wiki Import Blocking UI
**What goes wrong:** Importing an entire ADO wiki with 100+ pages fetches content sequentially, blocking the Wails binding call for minutes.
**Why it happens:** Wails bindings are synchronous from the frontend's perspective.
**How to avoid:** Import in batches — fetch tree structure first (fast), then import page content in background goroutine. Emit Wails events (`wiki:import:progress`, `wiki:import:complete`) for frontend progress tracking. Same pattern as existing `sync:started`/`sync:completed` events.
**Warning signs:** App UI freezes during large wiki imports.

## Code Examples

### Database Schema Extension
```sql
-- wiki_pages table
CREATE TABLE IF NOT EXISTS wiki_pages (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    content     TEXT DEFAULT '',
    source      TEXT DEFAULT 'local' CHECK(source IN ('local','ado')),
    is_readonly INTEGER DEFAULT 0,
    ado_wiki_id TEXT DEFAULT '',     -- ADO wiki identifier (for re-fetch)
    ado_page_path TEXT DEFAULT '',   -- ADO page path (for re-fetch)
    ado_org     TEXT DEFAULT '',     -- ADO org (for multi-org)
    ado_project TEXT DEFAULT '',     -- ADO project
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- FTS5 index (external content table)
CREATE VIRTUAL TABLE IF NOT EXISTS wiki_fts USING fts5(
    title, content,
    content=wiki_pages, content_rowid=id,
    tokenize='unicode61 remove_diacritics 2'
);

-- Sync triggers for FTS5
CREATE TRIGGER IF NOT EXISTS wiki_fts_insert AFTER INSERT ON wiki_pages BEGIN
    INSERT INTO wiki_fts(rowid, title, content)
    VALUES (new.id, new.title, new.content);
END;

CREATE TRIGGER IF NOT EXISTS wiki_fts_delete BEFORE DELETE ON wiki_pages BEGIN
    INSERT INTO wiki_fts(wiki_fts, rowid, title, content)
    VALUES ('delete', old.id, old.title, old.content);
END;

CREATE TRIGGER IF NOT EXISTS wiki_fts_update_before BEFORE UPDATE ON wiki_pages BEGIN
    INSERT INTO wiki_fts(wiki_fts, rowid, title, content)
    VALUES ('delete', old.id, old.title, old.content);
END;

CREATE TRIGGER IF NOT EXISTS wiki_fts_update_after AFTER UPDATE ON wiki_pages BEGIN
    INSERT INTO wiki_fts(rowid, title, content)
    VALUES (new.id, new.title, new.content);
END;

-- Project bookmarks (extends task_links pattern)
CREATE TABLE IF NOT EXISTS project_bookmarks (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    label      TEXT DEFAULT '',
    type       TEXT DEFAULT 'url',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_wiki_pages_source ON wiki_pages(source);
CREATE INDEX IF NOT EXISTS idx_wiki_pages_ado ON wiki_pages(ado_wiki_id, ado_page_path);
CREATE INDEX IF NOT EXISTS idx_project_bookmarks_project ON project_bookmarks(project_id);
```

### Domain Types (Go)
```go
// domain/wiki.go
package domain

import "time"

type WikiPage struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    Source      string    `json:"source"`      // "local" or "ado"
    IsReadonly  bool      `json:"isReadonly"`
    AdoWikiID   string    `json:"adoWikiId"`
    AdoPagePath string    `json:"adoPagePath"`
    AdoOrg      string    `json:"adoOrg"`
    AdoProject  string    `json:"adoProject"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

type WikiSearchResult struct {
    ID       int    `json:"id"`
    Title    string `json:"title"`
    Source   string `json:"source"`
    Preview  string `json:"preview"`  // snippet from FTS5
    IsReadonly bool `json:"isReadonly"`
}

type ProjectBookmark struct {
    ID        int       `json:"id"`
    ProjectID int       `json:"projectId"`
    URL       string    `json:"url"`
    Label     string    `json:"label"`
    Type      string    `json:"type"`
    CreatedAt time.Time `json:"createdAt"`
}
```

### Frontend Types (TypeScript)
```typescript
// Extend frontend/src/types/index.ts

export interface WikiPage {
  id: number
  title: string
  content: string
  source: 'local' | 'ado'
  isReadonly: boolean
  adoWikiId: string
  adoPagePath: string
  adoOrg: string
  adoProject: string
  createdAt: string
  updatedAt: string
}

export interface WikiSearchResult {
  id: number
  title: string
  source: string
  preview: string
  isReadonly: boolean
}

export interface ProjectBookmark {
  id: number
  projectId: number
  url: string
  label: string
  type: string
  createdAt: string
}

export interface ADOWikiInfo {
  id: string
  name: string
  type: string  // "projectWiki" | "codeWiki"
}

export interface ADOWikiPageInfo {
  id: number
  path: string
  order: number
  isParentPage: boolean
  content: string
  subPages: ADOWikiPageInfo[]
}
```

### FTS5 Search with Input Sanitization
```go
// internal/db/wiki.go
func (db *DB) SearchWikiPages(query string) ([]domain.WikiSearchResult, error) {
    // Sanitize: wrap in double quotes to prevent FTS5 syntax errors
    // This treats the entire input as a phrase search
    safeQuery := `"` + strings.ReplaceAll(query, `"`, `""`) + `"`

    rows, err := db.Query(`
        SELECT wp.id, wp.title, wp.source, wp.is_readonly,
               snippet(wiki_fts, 1, '<mark>', '</mark>', '...', 32) as preview
        FROM wiki_fts
        JOIN wiki_pages wp ON wp.id = wiki_fts.rowid
        WHERE wiki_fts MATCH ?
        ORDER BY rank
        LIMIT 20`, safeQuery)
    if err != nil {
        return nil, fmt.Errorf("search wiki pages: %w", err)
    }
    defer rows.Close()

    var results []domain.WikiSearchResult
    for rows.Next() {
        var r domain.WikiSearchResult
        if err := rows.Scan(&r.ID, &r.Title, &r.Source, &r.IsReadonly, &r.Preview); err != nil {
            return nil, err
        }
        results = append(results, r)
    }
    return results, rows.Err()
}
```

### Wiki Link Click Handling (Vue)
```typescript
// In WikiPreview.vue — handle clicks on wiki links
function handleWikiLinkClick(event: MouseEvent) {
  const target = event.target as HTMLElement
  const link = target.closest('a[data-wiki-page]')
  if (link) {
    event.preventDefault()
    const pageName = link.getAttribute('data-wiki-page')
    if (pageName) {
      const page = wikiStore.pages.find(p => p.title === pageName)
      if (page) {
        wikiStore.selectPage(page.id)
      } else {
        // Create-on-click for missing pages
        wikiStore.createPage(pageName, '')
      }
    }
  }
}
```

### ADO Wiki Re-fetch Pattern
```go
// internal/app/adowikiservice.go
func (s *ADOWikiService) RefreshImportedPages() error {
    // Find all ADO-sourced pages grouped by org/project/wiki
    pages, err := s.db.ListADOWikiPages()
    if err != nil {
        return err
    }

    // Group by org+project+wikiID for efficient batching
    groups := groupByWiki(pages)
    for key, group := range groups {
        token, err := s.tokenProv.GetToken()
        if err != nil {
            continue // skip on auth error, try next group
        }
        client := ado.NewClient(key.Org, key.Project, token)
        for _, page := range group {
            content, err := ado.GetWikiPageContent(client, page.AdoWikiID, page.AdoPagePath)
            if err != nil {
                log.Printf("Failed to refresh ADO wiki page %s: %v", page.AdoPagePath, err)
                continue
            }
            s.db.UpdateWikiPageContent(page.ID, content.Content)
        }
    }
    return nil
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Custom markdown parsers | Use `marked` or `markdown-it` | Stable for years | Never hand-roll a markdown parser |
| FTS3/FTS4 | FTS5 (default since SQLite 3.9.0, 2015) | 2015 | Better ranking (BM25), column filters, external content |
| ADO Wiki API v5.x | Wiki API v7.1 (current) | 2023 | pagesBatch endpoint, better pagination |
| WYSIWYG-first editors | Edit/preview toggle (Obsidian, GitHub) | 2020+ | Users prefer seeing raw markdown; simpler, fewer bugs |

**Deprecated/outdated:**
- **FTS3/FTS4**: Superseded by FTS5; no reason to use older versions
- **ADO API v5.x**: Use v7.1 for wiki endpoints; v7.0+ for consistency with existing work item calls

## Open Questions

1. **ADO Wiki re-fetch interval**
   - What we know: ADO wiki content changes infrequently (days/weeks between edits)
   - What's unclear: Exact interval that balances freshness vs. API load
   - Recommendation: Default to 1 hour. Configurable in settings. Manual refresh button always available. This matches the existing sync pattern (15 min for work items, wikis change less frequently)

2. **FTS5 tokenizer choice**
   - What we know: `unicode61` is the default and handles most languages. `porter` adds English stemming.
   - What's unclear: Whether the user's wiki content will be primarily English
   - Recommendation: Use `unicode61 remove_diacritics 2` (good multilingual support). Don't add porter stemming — it complicates exact-match searches and the wiki content is mixed technical/English.

3. **Sidebar responsive behavior**
   - What we know: Desktop app with min-width 1024px. Sidebar needs to coexist with editor/preview.
   - Recommendation: Fixed 280px sidebar width. At minimum window size (1024px), content area gets 744px minus icon sidebar (56px) = ~688px — plenty for markdown editing. No collapsible sidebar needed for v1.

## Sources

### Primary (HIGH confidence)
- **Codebase analysis**: `internal/db/migrate.go`, `internal/app/externallinks.go`, `pkg/ado/client.go`, `pkg/ado/query.go` — direct inspection of existing patterns
- **Codebase analysis**: `frontend/src/stores/tasks.ts`, `frontend/src/stores/ado.ts` — Pinia store patterns
- **Codebase analysis**: `frontend/src/components/tasks/ExternalLinks.vue` — Reusable component pattern for bookmarks
- **modernc.org/sqlite**: FTS5 is compiled in by default in the pure-Go SQLite translation (verified from go.mod v1.48.0)
- **SQLite FTS5 documentation**: https://www.sqlite.org/fts5.html — External content tables, triggers, snippet(), rank

### Secondary (MEDIUM confidence)
- **ADO Wiki REST API**: https://learn.microsoft.com/en-us/rest/api/azure/devops/wiki/ — Wiki pages API v7.1 endpoints, authentication, pagination
- **marked library**: https://marked.js.org/ — Version 18.0.0 confirmed via npm registry

### Tertiary (LOW confidence)
- **ADO Wiki API pagination behavior**: The pagesBatch endpoint's exact continuationToken handling may vary — test with actual ADO wiki data during implementation

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all libraries verified in npm registry, existing deps confirmed in package.json/go.mod
- Architecture: HIGH — all patterns directly derived from existing codebase (ExternalLinks, ADO browser, Pinia stores, migration system)
- Pitfalls: HIGH — FTS5 trigger patterns and MATCH syntax issues are well-documented in SQLite docs
- ADO Wiki API: MEDIUM — API endpoints documented by Microsoft, but exact response shapes should be validated during implementation

**Research date:** 2025-07-22
**Valid until:** 2025-08-22 (stable domain — markdown rendering, SQLite FTS5, and ADO REST API are all mature)
