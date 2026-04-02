# Coding Conventions

**Analysis Date:** 2025-06-26

## Naming Patterns

**Go Files:**
- Package names: lowercase single word — `app`, `auth`, `config`, `db`, `domain`
- Service structs: PascalCase with `Service` suffix — `TaskService`, `ProjectService`, `AuthService`, `ConfigService`, `DependencyService`
- Constructors: `New` prefix — `NewTaskService()`, `NewAuthService()`
- Methods: PascalCase (exported) — `Create()`, `GetByID()`, `List()`, `SetStatus()`, `ListFiltered()`
- Private functions: camelCase — `scanTasks()`, `hasCircularDep()`, `setDefaults()`, `ensureConfigDir()`
- Constants: PascalCase for exported, camelCase for unexported — `AppName`, `serviceName`, `callbackTimeout`
- Domain type constants: PascalCase with type prefix — `ProjectStatusActive`, `ADOPullRequestStatusDraft`

**Vue Files:**
- Components: PascalCase single-file — `TaskRow.vue`, `TaskDetail.vue`, `CommandPalette.vue`, `PageHeader.vue`
- Views: PascalCase with `View` suffix — `TasksView.vue`, `DashboardView.vue`, `LoginView.vue`, `SettingsView.vue`
- Layouts: PascalCase with description — `AppShell.vue`
- Composables: camelCase with `use` prefix — `useTheme.ts`
- Stores: camelCase noun — `tasks.ts`, `projects.ts`, `auth.ts`
- Store hook functions: `use` prefix + PascalCase + `Store` — `useTaskStore`, `useProjectStore`, `useAuthStore`
- UI components (shadcn-vue): PascalCase matching component name — `Button.vue`, `Badge.vue`, `Card.vue`

**TypeScript:**
- Interfaces: PascalCase — `Task`, `Project`, `User`, `PR`, `Action`, `ConfigState`
- Functions: camelCase — `fetchTasks()`, `createTask()`, `selectTask()`, `toggleSection()`
- Refs: camelCase — `showNewTask`, `newTaskTitle`, `collapsedSections`, `filterStatus`
- Computed: camelCase — `selectedTask`, `filteredTasks`, `grouped`, `visibleSections`
- Constants: camelCase for arrays/objects, UPPER_SNAKE for mock data — `sectionOrder`, `sectionLabels`, `MOCK_TASKS`

## Code Style

**Formatting:**
- No ESLint or Prettier config detected — formatting is manual/editor-based
- Go: standard `gofmt` formatting (implicit via Go toolchain)
- TypeScript: 2-space indentation, single quotes in Vue/TS files
- Vue templates: 2-space indentation

**Linting:**
- No `.eslintrc`, `.prettierrc`, or `biome.json` detected
- TypeScript strict mode enabled in `frontend/tsconfig.json` but with `noImplicitAny: false` and `noUnusedParameters: false` (relaxed)
- Go: no golint or golangci-lint config detected

**Line Length:**
- No enforced limit; template lines in Vue can be long (100+ chars is common for Tailwind classes)

## Import Organization

**Go Import Order (3 groups, blank-line separated):**
1. Standard library (`fmt`, `log`, `sort`, `strings`, `net/http`, etc.)
2. Third-party packages (`github.com/spf13/viper`, `github.com/wailsapp/wails/v3`, etc.)
3. Internal packages (`dev.azure.com/xbox/xb-tasks/domain`, `dev.azure.com/xbox/xb-tasks/internal/db`)

Example from `internal/auth/auth.go`:
```go
import (
    "context"
    "crypto/rand"
    "fmt"
    "net/http"

    "github.com/pkg/browser"
    "github.com/wailsapp/wails/v3/pkg/application"

    "dev.azure.com/xbox/xb-tasks/internal/db"
    "dev.azure.com/xbox/xb-tasks/domain"
)
```

**Vue/TS Import Order (no enforced rule, but observed pattern):**
1. Vue core (`vue`, `vue-router`)
2. Store imports (`@/stores/tasks`)
3. UI component imports (`@/components/ui/button`)
4. Custom component imports (`@/components/TaskDetail.vue`)
5. Utility imports (`@/lib/utils`)
6. Icon imports (`lucide-vue-next`)

**Path Aliases:**
- `@/` maps to `frontend/src/` — configured in `frontend/vite.config.ts` and `frontend/tsconfig.json`
- Wails bindings imported via relative path: `../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice`

## Vue Component Patterns

**Script Setup (Composition API exclusively):**
- All components use `<script setup lang="ts">` — no Options API
- No `defineComponent()` usage anywhere

**Props Pattern:**
```vue
<script setup lang="ts">
const props = defineProps<{
  task: Task
  selected?: boolean
}>()
```

**Emits Pattern:**
```vue
const emit = defineEmits<{
  select: [id: number]
  toggleStatus: [id: number]
}>()
```
Or for simple emits:
```vue
const emit = defineEmits<{ close: [] }>()
```

**Store Usage:**
- Access store at top of `<script setup>`: `const taskStore = useTaskStore()`
- Call store actions in lifecycle: `onMounted(() => { taskStore.fetchTasks() })`
- Bind store state directly in templates: `taskStore.selectedTask`, `taskStore.loading`

**Slot-based Layout Pattern:**
- `PageHeader.vue` uses named slots `#left` and `#right` for consistent page headers
- `AppShell.vue` uses default slot for page content

**Computed → Switch Pattern for Status/Priority Mapping:**
- Common pattern: `computed()` with `switch` on status/priority to return icon, color, label config
- Used in `StatusBadge.vue`, `PriorityBadge.vue`, `TaskRow.vue`, `TasksView.vue`

Example from `frontend/src/components/ui/StatusBadge.vue`:
```typescript
const config = computed(() => {
  switch (props.status) {
    case 'todo': return { icon: Circle, variant: 'outline' as const, classes: '...', label: 'To Do' }
    case 'in_progress': return { icon: CircleDot, variant: 'default' as const, classes: '...', label: 'In Progress' }
    // ...
  }
})
```

## shadcn-vue / UI Component Patterns

**Component Location:**
- All shadcn-vue primitives live in `frontend/src/components/ui/{component-name}/`
- Each component directory has an `index.ts` barrel file re-exporting the component(s) and variant definitions

**Barrel File Pattern:**
```typescript
// frontend/src/components/ui/button/index.ts
export { default as Button } from "./Button.vue"
export const buttonVariants = cva("...", { variants: { ... } })
export type ButtonVariants = VariantProps<typeof buttonVariants>
```

**CVA (class-variance-authority) for Variants:**
- Used in `button/index.ts` and `badge/index.ts`
- Variants defined with `cva()`, typed with `VariantProps`

**cn() Utility for Tailwind Merging:**
- Defined in `frontend/src/lib/utils.ts`
- Combines `clsx` + `tailwind-merge` for conditional class merging
- Used throughout all components: `cn('base-classes', condition && 'conditional-class')`

**Custom UI Components (not shadcn primitives):**
- `StatusBadge.vue`, `PriorityBadge.vue`, `TagChip.vue`, `AdoBadge.vue` — app-specific badge components in `frontend/src/components/ui/`
- These compose the shadcn `Badge` component with domain-specific logic

**Icon Usage:**
- Icons from `lucide-vue-next` imported individually
- Rendered with `<component :is="iconComponent" :size="N" :stroke-width="N" />`
- Common sizes: 10, 12, 13, 14, 16, 18, 20
- Common stroke widths: 1.75, 2, 2.5

## Pinia Store Patterns

**Setup Store Syntax (Composition API style):**
```typescript
export const useTaskStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  // ...computed, actions...
  return { tasks, loading, /* exposed members */ }
})
```

**Interface-First:**
- Each store defines its own TypeScript interface at the top of the file: `export interface Task { ... }`

**Mock-First Development Pattern:**
- All stores have `let useMock = true` flag and `MOCK_*` arrays for offline development
- Store actions check `if (useMock)` and fall back to mock data
- On Wails binding failure (`catch`), set `useMock = true` and load mocks
- This enables running the frontend in a browser without the Go backend

**Async Dynamic Import for Wails Bindings:**
```typescript
const { List } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
tasks.value = (await List(status)) as Task[]
```

**Loading State Pattern:**
- `loading` ref set `true` before async call, `false` in `finally`

## Error Handling

**Go Backend:**
- Use `fmt.Errorf("context: %w", err)` for wrapping errors — consistent throughout all service methods
- Validation returns zero-value struct + error: `return domain.Task{}, fmt.Errorf("title is required")`
- Row scan errors: `log.Printf()` + `continue` (skip bad row, don't fail entire list)
- Fatal errors only at startup in `main.go`: `log.Fatalf("failed to init config: %v", err)`

Example from `internal/app/tasks.go`:
```go
if err != nil {
    return domain.Task{}, fmt.Errorf("create task: %w", err)
}
```

**Vue Frontend:**
- `try/catch` with empty `catch` blocks — errors silently swallowed in most cases
- Auth store has `error` ref for user-facing error messages
- No global error handler or toast notification system
- Wails binding failures trigger mock fallback (graceful degradation)

## Logging

**Go Backend:**
- Standard library `log` package — `log.Printf()` for warnings, `log.Fatalf()` for fatal startup errors
- No structured logging framework (no slog, zerolog, etc.)
- Config has `log.level` setting but it's not wired to any logger

**Vue Frontend:**
- No logging framework
- No `console.log` statements in committed code
- No error tracking service

## Comments

**Go:**
- Godoc-style comments on exported types and methods: `// TaskService is bound to the frontend via Wails.`
- Section comments in SQL schema: table creation blocks
- Inline comments for non-obvious logic: `// Verify parent exists`, `// Check for circular dependency before inserting`

**Vue/TS:**
- Section separator comments with Unicode box drawing: `// ── Store & emits ──`
- JSDoc on `PageHeader.vue` with `/** ... */` block
- No JSDoc/TSDoc on functions or store methods

## Tailwind CSS Patterns

**Version:** Tailwind v4 with `@tailwindcss/vite` plugin and `tw-animate-css`

**Theme:**
- CSS custom properties defined in `frontend/src/style.css` for both `:root` (light) and `.dark` (dark)
- Uses shadcn-vue's oklch color system (`--primary`, `--foreground`, `--muted-foreground`, etc.)
- Custom surface/border/text variables layered on top: `--surface-primary`, `--text-primary`, `--border-default`

**Common Class Patterns:**
- Sizes: `text-[10px]`, `text-[11px]`, `text-[13px]`, `text-xs`, `text-sm` — fine-grained sizes for information density
- Opacity: `text-muted-foreground/40`, `bg-primary/[0.06]`, `border-border/50` — slash notation for opacity
- Responsive hiding: `hidden lg:inline`
- Tabular numbers: `tabular-nums` class for numeric alignment
- Truncation: `truncate`, `line-clamp-2`

**Desktop App Specific:**
- `titlebar-drag` / `titlebar-no-drag` CSS classes for macOS titlebar drag regions
- `-webkit-app-region: drag` / `no-drag` in `style.css`

## Module Design

**Go Services:**
- One service per file in `internal/app/`: `tasks.go`, `projects.go`, `deps.go`
- Each service struct holds a `*db.DB` reference, injected via constructor
- Services registered with Wails in `main.go` via `application.NewService()`
- No interfaces defined for services (concrete types only)

**Vue Exports:**
- Components exported as default from `.vue` files
- UI components re-exported from barrel `index.ts` files
- Stores export the `defineStore` result and its interface
- Composables export named functions

**Wails Bindings (auto-generated):**
- Located in `frontend/bindings/dev.azure.com/xbox/xb-tasks/internal/`
- Mirror the Go package structure
- Each service gets its own `.ts` file (e.g., `taskservice.ts`, `projectservice.ts`)
- Imported dynamically in stores — never statically

---

*Convention analysis: 2025-06-26*
