# Phase 1: Foundation, Auth & Personal Tasks - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-01
**Phase:** 01-foundation-auth-personal-tasks
**Areas discussed:** Navigation & app shell

---

## Navigation Structure

| Option | Description | Selected |
|--------|-------------|----------|
| Slim icon sidebar (like Linear/Notion) | Icons for sections, content takes most of the screen | |
| Full sidebar with labels (like VS Code) | Section names always visible, collapsible | |
| Top nav bar | Horizontal tabs for major sections, no sidebar | |
| Command palette only (like Raycast) | Minimal chrome, navigate via keyboard shortcuts | |

**User's choice:** Slim icon sidebar with a command palette for powerful search and management (hybrid of options 1 and 4)
**Notes:** User wants both visual navigation AND keyboard-driven power access. ⌘K command palette is a first-class citizen alongside the sidebar.

---

## Window Close Behavior

| Option | Description | Selected |
|--------|-------------|----------|
| Minimize to tray, keep running | Always-available like Slack/1Password | ✓ |
| Close fully | Relaunch when needed, no tray | |
| User configurable | Option in settings for either behavior | |

**User's choice:** Minimize to tray, keep running
**Notes:** None — straightforward pick.

---

## Color Scheme

| Option | Description | Selected |
|--------|-------------|----------|
| Dark mode only | Engineering tool, matches ADO/VS Code dark workflows | |
| Light mode only | — | |
| Follow system preference | Auto dark/light | ✓ |
| User toggle | Dark/light switch in settings | |

**User's choice:** Follow system preference (auto dark/light)
**Notes:** None — respects OS-level appearance setting.

---

## Landing Page

| Option | Description | Selected |
|--------|-------------|----------|
| Dashboard overview | Today's work, recent tasks, quick stats | ✓ |
| Task list | Jump right into tasks — the primary workspace | |
| Last viewed page | Restore where you left off | |

**User's choice:** Dashboard overview (today's work, recent tasks, quick stats)
**Notes:** None.

---

## Agent's Discretion

- Auth flow (OAuth2 PKCE, PAT fallback, token storage)
- Task list layout (table/cards, grouping, filters)
- Task detail experience (inline, side panel, full page)
- Tags UX (free-form vs predefined)
- Dependency visualization

## Deferred Ideas

None — discussion stayed within phase scope.
