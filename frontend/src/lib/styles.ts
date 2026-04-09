/**
 * Shared color/class mappings for status, priority, and ADO work-item styling.
 * Single source of truth — imported by all components that render colored badges, dots, or icons.
 */

import type { Component } from 'vue'
import {
  Bug,
  CheckSquare,
  BookOpen,
  Star,
  Mountain,
  Crown,
  Trophy,
  ListChecks,
  Circle,
} from 'lucide-vue-next'

/** Text color + hover for interactive status icons (e.g. task-row check buttons). */
export function statusColor(status: string): string {
  switch (status) {
    case 'todo': return 'text-zinc-400 hover:text-zinc-500'
    case 'in_progress': return 'text-blue-500 hover:text-blue-600'
    case 'in_review': return 'text-violet-500 hover:text-violet-600'
    case 'done': return 'text-emerald-500 hover:text-emerald-600'
    case 'blocked': return 'text-red-500 hover:text-red-600'
    case 'cancelled': return 'text-zinc-400 hover:text-zinc-500'
    default: return 'text-zinc-400'
  }
}

/** Background color for small status-indicator dots. */
export function statusBgColor(status: string): string {
  switch (status) {
    case 'in_progress': return 'bg-blue-500'
    case 'in_review': return 'bg-violet-500'
    case 'todo': return 'bg-zinc-400'
    case 'blocked': return 'bg-red-500'
    case 'done': return 'bg-emerald-500'
    case 'cancelled': return 'bg-zinc-400'
    default: return 'bg-zinc-400'
  }
}

/** Full bg+text+border classes for status badges. Handles both underscore and space separators. */
export function statusClasses(status: string): string {
  switch (status.toLowerCase().replace(/_/g, ' ')) {
    case 'in progress': return 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/25'
    case 'in review': return 'bg-violet-500/15 text-violet-700 dark:text-violet-400 border-violet-500/25'
    case 'done': return 'bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25'
    case 'blocked': return 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25'
    case 'cancelled': return 'bg-zinc-500/15 text-zinc-500 border-zinc-500/25'
    default: return 'bg-muted text-muted-foreground border-border'
  }
}

/** Full bg+text+border classes for priority badges (high / medium / low). */
export function priorityClasses(priority: string): string {
  switch (priority.toLowerCase()) {
    case 'high': return 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25'
    case 'medium': return 'bg-yellow-500/15 text-yellow-700 dark:text-yellow-400 border-yellow-500/25'
    case 'low': return 'bg-muted text-muted-foreground border-border'
    default: return 'bg-muted text-muted-foreground border-border'
  }
}

/** Text color for P0–P3 priority labels. */
export function priorityColor(priority: string): string {
  switch (priority) {
    case 'P0': return 'text-red-600 dark:text-red-400'
    case 'P1': return 'text-orange-600 dark:text-orange-400'
    case 'P2': return 'text-amber-600 dark:text-amber-400'
    case 'P3': return 'text-zinc-500'
    default: return 'text-zinc-500'
  }
}

/** Background color for priority dots (P0-P3). Used on dashboard task rows. */
export function priorityBgColor(priority: string): string {
  switch (priority) {
    case 'P0': return 'bg-red-500'
    case 'P1': return 'bg-orange-500'
    case 'P2': return 'bg-amber-500'
    case 'P3': return 'bg-zinc-400'
    default: return 'bg-zinc-400'
  }
}

/** Text color for ADO work-item types. Matches ADO's native color coding. */
export function adoTypeColor(type: string): string {
  switch (type.toLowerCase()) {
    case 'bug': return 'text-red-500'
    case 'task': return 'text-amber-500'
    case 'task group': return 'text-orange-500'
    case 'user story': return 'text-purple-500'
    case 'feature': return 'text-green-500'
    case 'scenario': return 'text-purple-400'
    case 'deliverable': return 'text-blue-500'
    case 'epic': return 'text-orange-500'
    default: return 'text-muted-foreground'
  }
}

/** Hex color for ADO work-item types. Use in SVG fills/strokes where Tailwind classes don't apply. */
export function adoTypeHex(type?: string): string {
  switch ((type || '').toLowerCase()) {
    case 'bug': return '#ef4444'
    case 'task': return '#f59e0b'
    case 'task group': return '#f97316'
    case 'user story': return '#a855f7'
    case 'feature': return '#22c55e'
    case 'scenario': return '#c084fc'
    case 'deliverable': return '#3b82f6'
    case 'epic': return '#f97316'
    default: return '#71717a'
  }
}

/** Lucide icon component for ADO work-item types. */
export function adoTypeIcon(type?: string): Component {
  switch ((type || '').toLowerCase()) {
    case 'bug': return Bug
    case 'task': return CheckSquare
    case 'task group': return ListChecks
    case 'user story': return BookOpen
    case 'feature': return Star
    case 'scenario': return Crown
    case 'deliverable': return Trophy
    case 'epic': return Mountain
    default: return Circle
  }
}

/** Full bg+text+border classes for ADO state badges (New, Active, Resolved, Closed, Removed). */
export function adoStateClasses(state: string): string {
  switch (state) {
    case 'Active': return 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/25'
    case 'New': return 'bg-muted text-muted-foreground border-border'
    case 'Resolved': return 'bg-yellow-500/15 text-yellow-700 dark:text-yellow-400 border-yellow-500/25'
    case 'Closed': return 'bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25'
    case 'Removed': return 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25'
    default: return 'bg-muted text-muted-foreground border-border'
  }
}

/** Full bg+text+border classes for ADO numeric priority badges (1=critical … 4=low). */
export function adoPriorityClasses(p: number): string {
  switch (p) {
    case 1: return 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25'
    case 2: return 'bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/25'
    case 3: return 'bg-yellow-500/15 text-yellow-700 dark:text-yellow-400 border-yellow-500/25'
    default: return 'bg-muted text-muted-foreground border-border'
  }
}

/** Full bg+text+border classes for pull-request status badges. */
export function prStatusClasses(status: string): string {
  switch (status) {
    case 'active': return 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/20'
    case 'draft': return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
    case 'completed': return 'bg-violet-500/15 text-violet-600 dark:text-violet-400 border-violet-500/20'
    case 'abandoned': return 'bg-red-500/15 text-red-600 dark:text-red-400 border-red-500/20'
    default: return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
  }
}

/** Text color for external link types (icm, grafana, ado, wiki, url). */
export function linkTypeColor(type: string): string {
  switch (type) {
    case 'icm': return 'text-red-500'
    case 'grafana': return 'text-green-500'
    case 'ado': return 'text-blue-500'
    case 'wiki': return 'text-purple-500'
    case 'url':
    default: return 'text-muted-foreground'
  }
}
