<script setup lang="ts">
import { ref } from 'vue'
import type { Component } from 'vue'
import { cn } from '@/lib/utils'
import {
  LayoutDashboard, CheckSquare, FolderKanban, Settings,
  Sun, Moon, Search, Plus, Activity, Bell, Clock,
  ChevronRight, ChevronLeft, X,
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

// --- Types ---

interface MockTask {
  title: string
  priority: string
  status: string
  adoId: string
  time: string
}

interface NavItem {
  icon: Component
  label: string
  active: boolean
}

// --- Mock data ---

const mockTasks: MockTask[] = [
  { title: 'Fix auth redirect loop', priority: 'P0', status: 'in_progress', adoId: '#48291', time: '30m' },
  { title: 'Build sidebar navigation', priority: 'P1', status: 'in_progress', adoId: '', time: '1h' },
  { title: 'Review PR #234 \u2014 schema changes', priority: 'P1', status: 'in_review', adoId: '#48350', time: '2h' },
  { title: 'Add rate limiting to gateway', priority: 'P1', status: 'todo', adoId: '#48400', time: '3h' },
  { title: 'Migrate auth to MSAL v3', priority: 'P0', status: 'blocked', adoId: '#48100', time: '1d' },
  { title: 'Write E2E tests for task CRUD', priority: 'P2', status: 'done', adoId: '', time: '2d' },
]

const navItems: NavItem[] = [
  { icon: LayoutDashboard, label: 'Dashboard', active: false },
  { icon: CheckSquare, label: 'Tasks', active: true },
  { icon: FolderKanban, label: 'Projects', active: false },
  { icon: Settings, label: 'Settings', active: false },
]

const activityEvents = [
  { text: 'Task #48291 moved to In Progress', time: '2m ago' },
  { text: 'PR #234 review requested', time: '15m ago' },
  { text: 'New comment on #48350', time: '1h ago' },
  { text: 'Task #48400 created', time: '3h ago' },
  { text: 'Sprint planning completed', time: '1d ago' },
]

// --- State ---

const showActivityA = ref(false)
const showActivityB = ref(false)
const showActivityC = ref(false)
const sidebarBExpanded = ref(true)
const darkA = ref(false)
const darkB = ref(false)
const darkC = ref(false)

// --- Helpers ---

function priorityColor(p: string): string {
  switch (p) {
    case 'P0': return 'bg-red-500'
    case 'P1': return 'bg-orange-500'
    case 'P2': return 'bg-amber-400'
    default: return 'bg-zinc-300 dark:bg-zinc-600'
  }
}

function statusColor(s: string): string {
  switch (s) {
    case 'in_progress': return 'bg-blue-500'
    case 'in_review': return 'bg-violet-500'
    case 'todo': return 'bg-zinc-400'
    case 'blocked': return 'bg-red-500'
    case 'done': return 'bg-emerald-500'
    default: return 'bg-zinc-300'
  }
}

function statusLabel(s: string): string {
  switch (s) {
    case 'in_progress': return 'In Progress'
    case 'in_review': return 'In Review'
    case 'todo': return 'To Do'
    case 'blocked': return 'Blocked'
    case 'done': return 'Done'
    default: return s
  }
}
</script>

<template>
  <TooltipProvider :delay-duration="200">
    <ScrollArea class="flex-1 h-full">
      <div class="px-6 py-6 space-y-4">
        <!-- Header -->
        <div>
          <h1 class="text-lg font-semibold text-foreground">Shell / App Chrome Playground</h1>
          <p class="text-sm text-muted-foreground mt-1">
            Compare shell layout variants side by side. Each tab renders a self-contained mock layout.
          </p>
        </div>

        <!-- Variant tabs -->
        <Tabs default-value="compact">
          <TabsList>
            <TabsTrigger value="compact">A: Compact Rail</TabsTrigger>
            <TabsTrigger value="expanded">B: Expanded Nav</TabsTrigger>
            <TabsTrigger value="floating">C: Floating Panels</TabsTrigger>
          </TabsList>

          <!-- ===== VARIANT A: Compact Rail ===== -->
          <TabsContent value="compact">
            <div class="border rounded-xl overflow-hidden bg-background relative" style="height: 580px">
              <!-- macOS traffic lights -->
              <div class="absolute top-2.5 left-3 flex items-center gap-1.5 z-10">
                <div class="w-3 h-3 rounded-full bg-[#FF5F57]" />
                <div class="w-3 h-3 rounded-full bg-[#FFBD2E]" />
                <div class="w-3 h-3 rounded-full bg-[#28C840]" />
              </div>

              <div class="flex h-full">
                <!-- Sidebar: 48px icon rail -->
                <div class="w-12 shrink-0 border-r border-border/50 flex flex-col items-center bg-muted/30">
                  <div class="pt-[48px]" />
                  <div class="flex flex-col items-center gap-1 px-1">
                    <Tooltip v-for="item in navItems" :key="item.label">
                      <TooltipTrigger as-child>
                        <Button
                          :variant="item.active ? 'secondary' : 'ghost'"
                          size="icon-sm"
                          :class="cn('h-8 w-8', item.active && 'bg-primary/10 text-primary')"
                        >
                          <component :is="item.icon" :size="18" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent side="right">
                        <p>{{ item.label }}</p>
                      </TooltipContent>
                    </Tooltip>
                  </div>
                  <div class="mt-auto pb-3">
                    <Button variant="ghost" size="icon-sm" class="h-8 w-8" @click="darkA = !darkA">
                      <Sun v-if="darkA" :size="16" />
                      <Moon v-else :size="16" />
                    </Button>
                  </div>
                </div>

                <!-- Main area -->
                <div class="flex-1 flex flex-col min-w-0">
                  <!-- Top bar: 36px -->
                  <div class="h-9 shrink-0 border-b border-border/50 flex items-center px-3 gap-3">
                    <span class="text-xs font-medium text-muted-foreground">Tasks</span>
                    <ChevronRight :size="12" class="text-muted-foreground/50" />
                    <span class="text-xs font-medium text-foreground">All</span>

                    <!-- Stats dots -->
                    <div class="flex items-center gap-2.5 ml-auto mr-auto">
                      <Tooltip>
                        <TooltipTrigger>
                          <span class="flex items-center gap-1">
                            <span class="w-2 h-2 rounded-full bg-blue-500 inline-block" />
                            <span class="text-[10px] text-muted-foreground font-medium">2</span>
                          </span>
                        </TooltipTrigger>
                        <TooltipContent><p>In Progress</p></TooltipContent>
                      </Tooltip>
                      <Tooltip>
                        <TooltipTrigger>
                          <span class="flex items-center gap-1">
                            <span class="w-2 h-2 rounded-full bg-violet-500 inline-block" />
                            <span class="text-[10px] text-muted-foreground font-medium">1</span>
                          </span>
                        </TooltipTrigger>
                        <TooltipContent><p>In Review</p></TooltipContent>
                      </Tooltip>
                      <Tooltip>
                        <TooltipTrigger>
                          <span class="flex items-center gap-1">
                            <span class="w-2 h-2 rounded-full bg-red-500 inline-block" />
                            <span class="text-[10px] text-muted-foreground font-medium">1</span>
                          </span>
                        </TooltipTrigger>
                        <TooltipContent><p>Blocked</p></TooltipContent>
                      </Tooltip>
                      <Tooltip>
                        <TooltipTrigger>
                          <span class="flex items-center gap-1">
                            <span class="w-2 h-2 rounded-full bg-emerald-500 inline-block" />
                            <span class="text-[10px] text-muted-foreground font-medium">1</span>
                          </span>
                        </TooltipTrigger>
                        <TooltipContent><p>Done</p></TooltipContent>
                      </Tooltip>
                    </div>

                    <!-- Actions -->
                    <div class="flex items-center gap-1">
                      <Button variant="ghost" size="icon-sm" class="h-7 w-7">
                        <Search :size="14" />
                      </Button>
                      <Button variant="ghost" size="icon-sm" class="h-7 w-7">
                        <Plus :size="14" />
                      </Button>
                      <Button
                        :variant="showActivityA ? 'secondary' : 'ghost'"
                        size="icon-sm"
                        class="h-7 w-7"
                        @click="showActivityA = !showActivityA"
                      >
                        <Activity :size="14" />
                      </Button>
                    </div>
                  </div>

                  <!-- Content -->
                  <div class="flex-1 flex min-h-0">
                    <ScrollArea class="flex-1">
                      <div class="py-1">
                        <div
                          v-for="task in mockTasks"
                          :key="task.title"
                          class="flex items-center px-3 py-2 border-b border-border/30 hover:bg-muted/50 cursor-default"
                        >
                          <div :class="cn('w-0.5 h-5 rounded-full mr-3 shrink-0', priorityColor(task.priority))" />
                          <div :class="cn('w-2 h-2 rounded-full mr-2.5 shrink-0', statusColor(task.status))" />
                          <span class="text-sm text-foreground truncate flex-1">
                            {{ task.title }}
                            <span v-if="task.adoId" class="text-muted-foreground ml-1.5 text-xs">{{ task.adoId }}</span>
                          </span>
                          <span class="text-xs text-muted-foreground shrink-0 ml-3 tabular-nums">{{ task.time }}</span>
                        </div>
                      </div>
                    </ScrollArea>

                    <!-- Activity sidebar: 220px -->
                    <div
                      v-if="showActivityA"
                      class="w-[220px] shrink-0 border-l border-border/50 flex flex-col bg-muted/20"
                    >
                      <div class="px-3 py-2 flex items-center justify-between border-b border-border/30">
                        <span class="text-xs font-semibold text-foreground">Activity</span>
                        <Button variant="ghost" size="icon-sm" class="h-6 w-6" @click="showActivityA = false">
                          <X :size="12" />
                        </Button>
                      </div>
                      <ScrollArea class="flex-1">
                        <div class="py-1">
                          <div
                            v-for="(event, i) in activityEvents"
                            :key="i"
                            class="px-3 py-2 border-b border-border/20"
                          >
                            <p class="text-xs text-foreground leading-snug">{{ event.text }}</p>
                            <p class="text-[10px] text-muted-foreground mt-0.5">{{ event.time }}</p>
                          </div>
                        </div>
                      </ScrollArea>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <p class="text-xs text-muted-foreground mt-2 px-1">
              Key idea: Maximum content space, minimal chrome. 48px rail clears macOS traffic lights with pt-[48px].
            </p>
          </TabsContent>

          <!-- ===== VARIANT B: Expanded Nav ===== -->
          <TabsContent value="expanded">
            <div class="border rounded-xl overflow-hidden bg-background relative" style="height: 580px">
              <!-- macOS traffic lights -->
              <div class="absolute top-2.5 left-3 flex items-center gap-1.5 z-10">
                <div class="w-3 h-3 rounded-full bg-[#FF5F57]" />
                <div class="w-3 h-3 rounded-full bg-[#FFBD2E]" />
                <div class="w-3 h-3 rounded-full bg-[#28C840]" />
              </div>

              <div class="flex h-full">
                <!-- Sidebar: collapsible 200px / 48px -->
                <div
                  :class="cn(
                    'shrink-0 border-r border-border/50 flex flex-col bg-muted/30 transition-all duration-200 overflow-hidden',
                    sidebarBExpanded ? 'w-[200px]' : 'w-12'
                  )"
                >
                  <div class="pt-[48px]" />

                  <div :class="cn('flex flex-col gap-0.5', sidebarBExpanded ? 'px-2' : 'px-1 items-center')">
                    <template v-for="item in navItems" :key="item.label">
                      <!-- Collapsed: icon-only with tooltip -->
                      <Tooltip v-if="!sidebarBExpanded">
                        <TooltipTrigger as-child>
                          <Button
                            :variant="item.active ? 'secondary' : 'ghost'"
                            size="icon-sm"
                            :class="cn('h-8 w-8', item.active && 'bg-primary/10 text-primary')"
                          >
                            <component :is="item.icon" :size="18" />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent side="right">
                          <p>{{ item.label }}</p>
                        </TooltipContent>
                      </Tooltip>
                      <!-- Expanded: icon + label -->
                      <Button
                        v-else
                        :variant="item.active ? 'secondary' : 'ghost'"
                        size="sm"
                        :class="cn(
                          'justify-start gap-2.5 h-8',
                          item.active && 'bg-primary/10 text-primary'
                        )"
                      >
                        <component :is="item.icon" :size="16" />
                        <span class="text-sm truncate">{{ item.label }}</span>
                      </Button>
                    </template>
                  </div>

                  <!-- Bottom: theme + collapse -->
                  <div :class="cn('mt-auto pb-2 flex flex-col gap-1', sidebarBExpanded ? 'px-2' : 'px-1 items-center')">
                    <Separator class="mb-1" />
                    <Button variant="ghost" size="icon-sm" class="h-8 w-8" @click="darkB = !darkB">
                      <Sun v-if="darkB" :size="16" />
                      <Moon v-else :size="16" />
                    </Button>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <Button variant="ghost" size="icon-sm" class="h-8 w-8" @click="sidebarBExpanded = !sidebarBExpanded">
                          <ChevronLeft v-if="sidebarBExpanded" :size="16" />
                          <ChevronRight v-else :size="16" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent side="right">
                        <p>{{ sidebarBExpanded ? 'Collapse' : 'Expand' }}</p>
                      </TooltipContent>
                    </Tooltip>
                  </div>
                </div>

                <!-- Main area -->
                <div class="flex-1 flex flex-col min-w-0">
                  <!-- Top bar: 44px -->
                  <div class="h-11 shrink-0 border-b border-border/50 flex items-center px-3 gap-3">
                    <span class="text-sm font-semibold text-foreground whitespace-nowrap">Team ADO Tool</span>

                    <!-- Search input -->
                    <div class="flex-1 max-w-md mx-auto">
                      <div class="flex items-center h-7 rounded-md border border-border/60 bg-muted/30 px-2.5 gap-2">
                        <Search :size="13" class="text-muted-foreground shrink-0" />
                        <span class="text-xs text-muted-foreground flex-1">Search...</span>
                        <Badge variant="outline" class="text-[10px] px-1.5 py-0 h-4 font-normal rounded-sm">
                          <span class="font-mono">&#8984;K</span>
                        </Badge>
                      </div>
                    </div>

                    <!-- Actions -->
                    <div class="flex items-center gap-1.5">
                      <Button variant="default" size="sm" class="h-7 text-xs gap-1.5 px-2.5">
                        <Plus :size="13" />
                        New Task
                      </Button>
                      <Button variant="ghost" size="icon-sm" class="h-7 w-7">
                        <Bell :size="14" />
                      </Button>
                      <Button
                        :variant="showActivityB ? 'secondary' : 'ghost'"
                        size="icon-sm"
                        class="h-7 w-7"
                        @click="showActivityB = !showActivityB"
                      >
                        <Activity :size="14" />
                      </Button>
                      <div class="w-6 h-6 rounded-full bg-primary/20 flex items-center justify-center ml-0.5">
                        <span class="text-[10px] font-semibold text-primary">JD</span>
                      </div>
                    </div>
                  </div>

                  <!-- Stats strip: 28px stacked progress -->
                  <div class="h-7 shrink-0 border-b border-border/30 flex items-center px-3 gap-0.5">
                    <Tooltip>
                      <TooltipTrigger class="flex-none">
                        <div class="h-2 rounded-l-full bg-blue-500" style="width: 80px" />
                      </TooltipTrigger>
                      <TooltipContent><p>In Progress: 2</p></TooltipContent>
                    </Tooltip>
                    <Tooltip>
                      <TooltipTrigger class="flex-none">
                        <div class="h-2 bg-violet-500" style="width: 40px" />
                      </TooltipTrigger>
                      <TooltipContent><p>In Review: 1</p></TooltipContent>
                    </Tooltip>
                    <Tooltip>
                      <TooltipTrigger class="flex-none">
                        <div class="h-2 bg-zinc-400" style="width: 40px" />
                      </TooltipTrigger>
                      <TooltipContent><p>To Do: 1</p></TooltipContent>
                    </Tooltip>
                    <Tooltip>
                      <TooltipTrigger class="flex-none">
                        <div class="h-2 bg-red-500" style="width: 40px" />
                      </TooltipTrigger>
                      <TooltipContent><p>Blocked: 1</p></TooltipContent>
                    </Tooltip>
                    <Tooltip>
                      <TooltipTrigger class="flex-1">
                        <div class="h-2 rounded-r-full bg-emerald-500 w-full" />
                      </TooltipTrigger>
                      <TooltipContent><p>Done: 1</p></TooltipContent>
                    </Tooltip>
                  </div>

                  <!-- Content: two-column dashboard -->
                  <div class="flex-1 flex min-h-0">
                    <ScrollArea class="flex-1">
                      <div class="p-3 grid grid-cols-2 gap-3">
                        <!-- Tasks column -->
                        <div>
                          <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2 px-1">Tasks</h3>
                          <div class="border rounded-lg border-border/50 overflow-hidden">
                            <div
                              v-for="(task, i) in mockTasks"
                              :key="task.title"
                              :class="cn(
                                'flex items-center px-3 py-2 hover:bg-muted/50',
                                i < mockTasks.length - 1 && 'border-b border-border/30'
                              )"
                            >
                              <div :class="cn('w-0.5 h-5 rounded-full mr-2.5 shrink-0', priorityColor(task.priority))" />
                              <div :class="cn('w-1.5 h-1.5 rounded-full mr-2 shrink-0', statusColor(task.status))" />
                              <span class="text-xs text-foreground truncate flex-1">{{ task.title }}</span>
                              <span class="text-[10px] text-muted-foreground shrink-0 ml-2 tabular-nums">{{ task.time }}</span>
                            </div>
                          </div>
                        </div>

                        <!-- PRs column -->
                        <div>
                          <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2 px-1">Pull Requests</h3>
                          <div class="border rounded-lg border-border/50 overflow-hidden">
                            <div class="flex items-center px-3 py-2 border-b border-border/30 hover:bg-muted/50">
                              <div class="w-1.5 h-1.5 rounded-full bg-emerald-500 mr-2 shrink-0" />
                              <span class="text-xs text-foreground truncate flex-1">feat: add task filtering</span>
                              <Badge variant="outline" class="text-[10px] px-1.5 py-0 h-4 ml-2 rounded-sm">#234</Badge>
                            </div>
                            <div class="flex items-center px-3 py-2 border-b border-border/30 hover:bg-muted/50">
                              <div class="w-1.5 h-1.5 rounded-full bg-violet-500 mr-2 shrink-0" />
                              <span class="text-xs text-foreground truncate flex-1">fix: auth token refresh</span>
                              <Badge variant="outline" class="text-[10px] px-1.5 py-0 h-4 ml-2 rounded-sm">#231</Badge>
                            </div>
                            <div class="flex items-center px-3 py-2 hover:bg-muted/50">
                              <div class="w-1.5 h-1.5 rounded-full bg-emerald-500 mr-2 shrink-0" />
                              <span class="text-xs text-foreground truncate flex-1">refactor: gateway middleware</span>
                              <Badge variant="outline" class="text-[10px] px-1.5 py-0 h-4 ml-2 rounded-sm">#228</Badge>
                            </div>
                          </div>
                        </div>
                      </div>
                    </ScrollArea>

                    <!-- Activity sidebar: 280px, grouped -->
                    <div
                      v-if="showActivityB"
                      class="w-[280px] shrink-0 border-l border-border/50 flex flex-col bg-muted/20"
                    >
                      <div class="px-3 py-2 flex items-center justify-between border-b border-border/30">
                        <span class="text-xs font-semibold text-foreground">Activity</span>
                        <Button variant="ghost" size="icon-sm" class="h-6 w-6" @click="showActivityB = false">
                          <X :size="12" />
                        </Button>
                      </div>
                      <ScrollArea class="flex-1">
                        <div class="text-[10px] font-semibold text-muted-foreground uppercase tracking-wider px-3 pt-3 pb-1">Today</div>
                        <div
                          v-for="(event, i) in activityEvents.slice(0, 3)"
                          :key="'today-' + i"
                          class="px-3 py-2 border-b border-border/20"
                        >
                          <p class="text-xs text-foreground leading-snug">{{ event.text }}</p>
                          <p class="text-[10px] text-muted-foreground mt-0.5">{{ event.time }}</p>
                        </div>
                        <Separator />
                        <div class="text-[10px] font-semibold text-muted-foreground uppercase tracking-wider px-3 pt-3 pb-1">Yesterday</div>
                        <div
                          v-for="(event, i) in activityEvents.slice(3)"
                          :key="'yesterday-' + i"
                          class="px-3 py-2 border-b border-border/20"
                        >
                          <p class="text-xs text-foreground leading-snug">{{ event.text }}</p>
                          <p class="text-[10px] text-muted-foreground mt-0.5">{{ event.time }}</p>
                        </div>
                      </ScrollArea>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <p class="text-xs text-muted-foreground mt-2 px-1">
              Key idea: Discoverable navigation with icon+label sidebar. Integrated search with shortcut hint. Stats as stacked progress segments.
            </p>
          </TabsContent>

          <!-- ===== VARIANT C: Floating Panels ===== -->
          <TabsContent value="floating">
            <div class="border rounded-xl overflow-hidden bg-background relative" style="height: 580px">
              <!-- macOS traffic lights -->
              <div class="absolute top-2.5 left-3 flex items-center gap-1.5 z-10">
                <div class="w-3 h-3 rounded-full bg-[#FF5F57]" />
                <div class="w-3 h-3 rounded-full bg-[#FFBD2E]" />
                <div class="w-3 h-3 rounded-full bg-[#28C840]" />
              </div>

              <div class="flex h-full">
                <!-- Sidebar: 56px icon rail, rounded items, subtle bg -->
                <div class="w-14 shrink-0 border-r border-border/50 flex flex-col items-center bg-muted/20">
                  <div class="pt-[48px]" />
                  <div class="flex flex-col items-center gap-1.5 px-1.5">
                    <Tooltip v-for="item in navItems" :key="item.label">
                      <TooltipTrigger as-child>
                        <Button
                          :variant="item.active ? 'secondary' : 'ghost'"
                          size="icon-sm"
                          :class="cn(
                            'h-9 w-9 rounded-lg',
                            item.active && 'bg-primary/10 text-primary shadow-sm'
                          )"
                        >
                          <component :is="item.icon" :size="18" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent side="right">
                        <p>{{ item.label }}</p>
                      </TooltipContent>
                    </Tooltip>
                  </div>
                  <div class="mt-auto pb-3">
                    <Button variant="ghost" size="icon-sm" class="h-9 w-9 rounded-lg" @click="darkC = !darkC">
                      <Sun v-if="darkC" :size="16" />
                      <Moon v-else :size="16" />
                    </Button>
                  </div>
                </div>

                <!-- Main area -->
                <div class="flex-1 flex flex-col min-w-0 relative">
                  <!-- Top bar: 40px -->
                  <div class="h-10 shrink-0 border-b border-border/50 flex items-center px-3 gap-2">
                    <!-- Clickable breadcrumb segments -->
                    <Button variant="ghost" size="sm" class="h-6 px-1.5 text-xs text-muted-foreground">Dashboard</Button>
                    <ChevronRight :size="12" class="text-muted-foreground/50 shrink-0" />
                    <Button variant="ghost" size="sm" class="h-6 px-1.5 text-xs font-medium text-foreground">Tasks</Button>

                    <!-- Stats dots with tooltips -->
                    <div class="flex items-center gap-2 ml-auto mr-auto">
                      <Tooltip>
                        <TooltipTrigger>
                          <span class="flex items-center gap-1">
                            <span class="w-2 h-2 rounded-full bg-blue-500 inline-block" />
                            <span class="text-[10px] text-muted-foreground">2</span>
                          </span>
                        </TooltipTrigger>
                        <TooltipContent><p>In Progress: 2</p></TooltipContent>
                      </Tooltip>
                      <Tooltip>
                        <TooltipTrigger>
                          <span class="flex items-center gap-1">
                            <span class="w-2 h-2 rounded-full bg-violet-500 inline-block" />
                            <span class="text-[10px] text-muted-foreground">1</span>
                          </span>
                        </TooltipTrigger>
                        <TooltipContent><p>In Review: 1</p></TooltipContent>
                      </Tooltip>
                      <Tooltip>
                        <TooltipTrigger>
                          <span class="flex items-center gap-1">
                            <span class="w-2 h-2 rounded-full bg-red-500 inline-block" />
                            <span class="text-[10px] text-muted-foreground">1</span>
                          </span>
                        </TooltipTrigger>
                        <TooltipContent><p>Blocked: 1</p></TooltipContent>
                      </Tooltip>
                      <Tooltip>
                        <TooltipTrigger>
                          <span class="flex items-center gap-1">
                            <span class="w-2 h-2 rounded-full bg-emerald-500 inline-block" />
                            <span class="text-[10px] text-muted-foreground">1</span>
                          </span>
                        </TooltipTrigger>
                        <TooltipContent><p>Done: 1</p></TooltipContent>
                      </Tooltip>
                    </div>

                    <!-- Actions -->
                    <div class="flex items-center gap-1">
                      <Button variant="ghost" size="icon-sm" class="h-7 w-7">
                        <Search :size="14" />
                      </Button>
                      <Button variant="ghost" size="icon-sm" class="h-7 w-7">
                        <Plus :size="14" />
                      </Button>
                      <Button
                        :variant="showActivityC ? 'secondary' : 'ghost'"
                        size="icon-sm"
                        class="h-7 w-7"
                        @click="showActivityC = !showActivityC"
                      >
                        <Clock :size="14" />
                      </Button>
                    </div>
                  </div>

                  <!-- Content with floating inset -->
                  <div class="flex-1 p-1 min-h-0">
                    <ScrollArea class="h-full rounded-lg bg-muted/20">
                      <div class="py-1">
                        <div
                          v-for="task in mockTasks"
                          :key="task.title"
                          class="flex items-center px-3 py-2.5 mx-1 rounded-md hover:bg-muted/60 cursor-default"
                        >
                          <div :class="cn('w-0.5 h-5 rounded-full mr-3 shrink-0', priorityColor(task.priority))" />
                          <div :class="cn('w-2 h-2 rounded-full mr-2.5 shrink-0', statusColor(task.status))" />
                          <span class="text-sm text-foreground truncate flex-1">
                            {{ task.title }}
                            <span v-if="task.adoId" class="text-muted-foreground ml-1.5 text-xs">{{ task.adoId }}</span>
                          </span>
                          <Badge variant="outline" class="text-[10px] px-1.5 py-0 h-4 ml-2 shrink-0 rounded-sm">
                            {{ statusLabel(task.status) }}
                          </Badge>
                          <span class="text-xs text-muted-foreground shrink-0 ml-2 tabular-nums">{{ task.time }}</span>
                        </div>
                      </div>
                    </ScrollArea>
                  </div>

                  <!-- Floating activity dropdown panel -->
                  <div
                    v-if="showActivityC"
                    class="absolute top-11 right-2 w-80 bg-background border border-border rounded-xl shadow-lg z-20 flex flex-col overflow-hidden"
                    style="height: 400px"
                  >
                    <!-- Arrow indicator -->
                    <div class="absolute -top-1.5 right-4 w-3 h-3 bg-background border-l border-t border-border rotate-45 z-10" />
                    <div class="px-3 py-2.5 flex items-center justify-between border-b border-border/50 relative z-20 bg-background">
                      <span class="text-xs font-semibold text-foreground">Recent Activity</span>
                      <Button variant="ghost" size="icon-sm" class="h-6 w-6" @click="showActivityC = false">
                        <X :size="12" />
                      </Button>
                    </div>
                    <ScrollArea class="flex-1">
                      <div class="py-1">
                        <div
                          v-for="(event, i) in activityEvents"
                          :key="i"
                          class="px-3 py-2.5 hover:bg-muted/50"
                        >
                          <p class="text-xs text-foreground leading-snug">{{ event.text }}</p>
                          <p class="text-[10px] text-muted-foreground mt-0.5">{{ event.time }}</p>
                        </div>
                      </div>
                    </ScrollArea>
                  </div>
                </div>
              </div>
            </div>
            <p class="text-xs text-muted-foreground mt-2 px-1">
              Key idea: Clean, modern floating feel. Content has inset margin for breathing room. Activity is a dropdown panel, not a fixed sidebar.
            </p>
          </TabsContent>
        </Tabs>
      </div>
    </ScrollArea>
  </TooltipProvider>
</template>
