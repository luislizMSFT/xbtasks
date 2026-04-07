<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, shallowRef, triggerRef } from 'vue'
import {
  forceSimulation,
  forceLink,
  forceManyBody,
  forceCenter,
  forceCollide,
  forceX,
  forceY,
  type Simulation,
  type SimulationNodeDatum,
  type SimulationLinkDatum,
} from 'd3-force'
import { select } from 'd3-selection'
import { zoom, zoomIdentity, type ZoomBehavior } from 'd3-zoom'
import { drag, type DragBehavior, type SubjectPosition } from 'd3-drag'
import { adoTypeHex, adoTypeIcon } from '@/lib/styles'

// ─── Types ──────────────────────────────────────────────────────────────────

import type { ForceGraphNode, ForceGraphEdge, ForceGraphEdgeKind } from '@/types'

export type GraphNode = ForceGraphNode
export type GraphEdge = ForceGraphEdge

interface SimNode extends SimulationNodeDatum {
  id: number
  title: string
  status: string
  priority: string
  projectId: number | null
  source: 'local' | 'ado'
  type?: string
}

interface SimEdge extends SimulationLinkDatum<SimNode> {
  source: SimNode
  target: SimNode
  kind: ForceGraphEdgeKind
}

// ─── Props / Emits ──────────────────────────────────────────────────────────

const props = defineProps<{
  nodes: GraphNode[]
  edges: GraphEdge[]
  selectedId: number | null
}>()

const emit = defineEmits<{
  select: [id: number]
}>()

// ─── Constants ──────────────────────────────────────────────────────────────

const STATUS_COLORS: Record<string, string> = {
  done: '#10b981',
  in_progress: '#3b82f6',
  in_review: '#8b5cf6',
  blocked: '#ef4444',
  todo: '#71717a',
  cancelled: '#a1a1aa',
}

const PRIORITY_RADIUS: Record<string, number> = {
  P0: 16,
  P1: 12,
  P2: 9,
  P3: 7,
}

const HULL_COLORS = [
  '#3b82f6', '#8b5cf6', '#f59e0b', '#10b981',
  '#ec4899', '#06b6d4', '#f97316', '#6366f1',
]

// ─── Reactive state ─────────────────────────────────────────────────────────

const containerRef = ref<HTMLDivElement | null>(null)
const svgRef = ref<SVGSVGElement | null>(null)
const width = ref(800)
const height = ref(600)

const simulatedNodes = shallowRef<SimNode[]>([])
const simulatedEdges = shallowRef<SimEdge[]>([])

const transformStr = ref('translate(0,0) scale(1)')
const currentZoom = ref(1)
const showLabels = computed(() => currentZoom.value > 0.6)

let simulation: Simulation<SimNode, SimEdge> | null = null
let zoomBehavior: ZoomBehavior<SVGSVGElement, unknown> | null = null
let dragBehavior: DragBehavior<SVGElement, SimNode, SimNode | SubjectPosition> | null = null
let resizeObserver: ResizeObserver | null = null
let animFrame = 0
let isDragging = false
let dragStartX = 0
let dragStartY = 0
const DRAG_THRESHOLD = 3

// ─── Helpers ────────────────────────────────────────────────────────────────

function nodeRadius(priority: string): number {
  return PRIORITY_RADIUS[priority] ?? 8
}

function statusColor(status: string): string {
  return STATUS_COLORS[status] ?? '#71717a'
}

function truncate(s: string, max = 20): string {
  return s.length > max ? s.slice(0, max) + '…' : s
}

function isBlockedEdge(edge: SimEdge): boolean {
  return edge.target?.status === 'blocked'
}

function isParentEdge(edge: SimEdge): boolean {
  return edge.kind === 'parent'
}

function isAdoNode(node: SimNode): boolean {
  return node.source === 'ado'
}

// ─── Convex hull computation ────────────────────────────────────────────────

function cross(O: [number, number], A: [number, number], B: [number, number]): number {
  return (A[0] - O[0]) * (B[1] - O[1]) - (A[1] - O[1]) * (B[0] - O[0])
}

function convexHull(points: [number, number][]): [number, number][] {
  if (points.length < 3) return points
  const sorted = [...points].sort((a, b) => a[0] - b[0] || a[1] - b[1])
  const lower: [number, number][] = []
  for (const p of sorted) {
    while (lower.length >= 2 && cross(lower[lower.length - 2], lower[lower.length - 1], p) <= 0) {
      lower.pop()
    }
    lower.push(p)
  }
  const upper: [number, number][] = []
  for (const p of sorted.reverse()) {
    while (upper.length >= 2 && cross(upper[upper.length - 2], upper[upper.length - 1], p) <= 0) {
      upper.pop()
    }
    upper.push(p)
  }
  upper.pop()
  lower.pop()
  return lower.concat(upper)
}

const projectHulls = computed(() => {
  const groups = new Map<number, SimNode[]>()
  for (const n of simulatedNodes.value) {
    if (n.projectId != null && n.x != null && n.y != null) {
      let arr = groups.get(n.projectId)
      if (!arr) {
        arr = []
        groups.set(n.projectId, arr)
      }
      arr.push(n)
    }
  }

  const hulls: { projectId: number; path: string; color: string }[] = []
  let colorIdx = 0
  for (const [pid, nodes] of groups) {
    if (nodes.length < 3) continue
    const pts: [number, number][] = nodes.map(n => [n.x!, n.y!])
    const hull = convexHull(pts)
    if (hull.length < 3) continue

    // Expand hull outward by padding
    const pad = 24
    const cx = hull.reduce((s, p) => s + p[0], 0) / hull.length
    const cy = hull.reduce((s, p) => s + p[1], 0) / hull.length
    const expanded = hull.map(([x, y]) => {
      const dx = x - cx
      const dy = y - cy
      const dist = Math.sqrt(dx * dx + dy * dy) || 1
      return [x + (dx / dist) * pad, y + (dy / dist) * pad] as [number, number]
    })

    const d = expanded.map((p, i) => `${i === 0 ? 'M' : 'L'}${p[0]},${p[1]}`).join(' ') + ' Z'
    hulls.push({
      projectId: pid,
      path: d,
      color: HULL_COLORS[colorIdx % HULL_COLORS.length],
    })
    colorIdx++
  }
  return hulls
})

// ─── Project clustering force ───────────────────────────────────────────────

function forceCluster(strength: number) {
  let nodes: SimNode[] = []

  function force(alpha: number) {
    const centroids = new Map<number, { x: number; y: number; count: number }>()
    for (const n of nodes) {
      if (n.projectId == null) continue
      let c = centroids.get(n.projectId)
      if (!c) {
        c = { x: 0, y: 0, count: 0 }
        centroids.set(n.projectId, c)
      }
      c.x += n.x ?? 0
      c.y += n.y ?? 0
      c.count++
    }
    for (const c of centroids.values()) {
      c.x /= c.count
      c.y /= c.count
    }
    for (const n of nodes) {
      if (n.projectId == null) continue
      const c = centroids.get(n.projectId)!
      n.vx! += (c.x - n.x!) * strength * alpha
      n.vy! += (c.y - n.y!) * strength * alpha
    }
  }

  force.initialize = (_nodes: SimNode[]) => {
    nodes = _nodes
  }

  return force
}

// ─── Simulation setup ───────────────────────────────────────────────────────

function buildSimulation() {
  if (simulation) simulation.stop()

  const simNodes: SimNode[] = props.nodes.map(n => ({
    id: n.id,
    title: n.title,
    status: n.status,
    priority: n.priority,
    projectId: n.projectId,
    source: n.source,
    type: n.type,
    x: n.x ?? width.value / 2 + (Math.random() - 0.5) * 200,
    y: n.y ?? height.value / 2 + (Math.random() - 0.5) * 200,
  }))

  const nodeMap = new Map(simNodes.map(n => [n.id, n]))

  const simEdges: SimEdge[] = props.edges
    .filter(e => nodeMap.has(e.source) && nodeMap.has(e.target))
    .map(e => ({
      source: nodeMap.get(e.source)!,
      target: nodeMap.get(e.target)!,
      kind: e.kind,
    }))

  simulation = forceSimulation<SimNode>(simNodes)
    .force(
      'link',
      forceLink<SimNode, SimEdge>(simEdges)
        .id(d => d.id)
        .distance(120),
    )
    .force('charge', forceManyBody<SimNode>().strength(-300))
    .force('center', forceCenter(width.value / 2, height.value / 2))
    .force(
      'collide',
      forceCollide<SimNode>(d => nodeRadius(d.priority) + 8),
    )
    .force('x', forceX<SimNode>(width.value / 2).strength(0.05))
    .force('y', forceY<SimNode>(height.value / 2).strength(0.05))
    .force('cluster', forceCluster(0.15) as any)
    .on('tick', () => {
      cancelAnimationFrame(animFrame)
      animFrame = requestAnimationFrame(() => {
        simulatedNodes.value = simNodes
        simulatedEdges.value = simEdges
        triggerRef(simulatedNodes)
        triggerRef(simulatedEdges)
      })
    })

  simulatedNodes.value = simNodes
  simulatedEdges.value = simEdges
}

// ─── Zoom ───────────────────────────────────────────────────────────────────

function setupZoom() {
  if (!svgRef.value) return
  const svg = select(svgRef.value)

  zoomBehavior = zoom<SVGSVGElement, unknown>()
    .scaleExtent([0.2, 3])
    .on('zoom', (event) => {
      const t = event.transform
      transformStr.value = `translate(${t.x},${t.y}) scale(${t.k})`
      currentZoom.value = t.k
    })

  svg.call(zoomBehavior).on('dblclick.zoom', null)
}

// ─── Drag ───────────────────────────────────────────────────────────────────

function setupDrag() {
  dragBehavior = drag<SVGElement, SimNode>()
    .on('start', (event, d) => {
      isDragging = false
      dragStartX = event.x
      dragStartY = event.y
      if (!event.active) simulation?.alphaTarget(0.3).restart()
      d.fx = d.x
      d.fy = d.y
    })
    .on('drag', (event, d) => {
      const dx = Math.abs(event.x - dragStartX)
      const dy = Math.abs(event.y - dragStartY)
      if (dx > DRAG_THRESHOLD || dy > DRAG_THRESHOLD) {
        isDragging = true
      }
      d.fx = event.x
      d.fy = event.y
    })
    .on('end', (event, d) => {
      if (!event.active) simulation?.alphaTarget(0)
      if (isDragging) {
        // Keep position fixed after real drag
        d.fx = event.x
        d.fy = event.y
      } else {
        // Was a click, not a drag — emit select
        d.fx = null
        d.fy = null
        emit('select', d.id)
      }
      isDragging = false
    })
}

function applyDrag() {
  if (!svgRef.value || !dragBehavior) return
  const nodeElements = select(svgRef.value).selectAll<SVGElement, SimNode>('.node-circle')
  nodeElements.call(dragBehavior as any)
}

// ─── Node click handler ────────────────────────────────────────────────────

function onNodeClick(event: MouseEvent, nodeId: number) {
  event.stopPropagation()
  emit('select', nodeId)
}

function onBackgroundClick(event: MouseEvent) {
  // Only deselect if the click is directly on the SVG or root <g>
  const target = event.target as Element
  if (target.tagName === 'svg' || target.classList.contains('graph-bg')) {
    emit('select', -1)
  }
}

// ─── Resize observer ───────────────────────────────────────────────────────

function observeSize() {
  if (!containerRef.value) return
  resizeObserver = new ResizeObserver(entries => {
    for (const entry of entries) {
      width.value = entry.contentRect.width
      height.value = entry.contentRect.height
    }
    // Recenter forces on resize
    if (simulation) {
      simulation.force('center', forceCenter(width.value / 2, height.value / 2))
      simulation.force('x', forceX<SimNode>(width.value / 2).strength(0.05))
      simulation.force('y', forceY<SimNode>(height.value / 2).strength(0.05))
      simulation.alpha(0.3).restart()
    }
  })
  resizeObserver.observe(containerRef.value)
}

// ─── Lifecycle ──────────────────────────────────────────────────────────────

onMounted(() => {
  observeSize()
  buildSimulation()
  setupZoom()
  setupDrag()
  // Need to wait a tick for Vue to render circles before d3 can select them
  requestAnimationFrame(() => applyDrag())
})

watch(
  () => [props.nodes, props.edges],
  () => {
    buildSimulation()
    setupZoom()
    setupDrag()
    requestAnimationFrame(() => applyDrag())
  },
  { deep: true },
)

// Re-apply drag when simulation updates (new DOM circles)
watch(simulatedNodes, () => {
  requestAnimationFrame(() => applyDrag())
})

onUnmounted(() => {
  simulation?.stop()
  simulation = null
  resizeObserver?.disconnect()
  cancelAnimationFrame(animFrame)
})
</script>

<template>
  <div ref="containerRef" class="w-full h-full relative bg-background overflow-hidden">
    <svg
      ref="svgRef"
      :width="width"
      :height="height"
      class="absolute inset-0"
      @click="onBackgroundClick"
    >
      <defs>
        <!-- Default arrowhead -->
        <marker
          id="arrowhead"
          viewBox="0 0 10 10"
          refX="10"
          refY="5"
          markerWidth="6"
          markerHeight="6"
          orient="auto-start-reverse"
        >
          <path d="M 0 0 L 10 5 L 0 10 z" fill="#3f3f46" />
        </marker>
        <!-- Blocked arrowhead -->
        <marker
          id="arrowhead-blocked"
          viewBox="0 0 10 10"
          refX="10"
          refY="5"
          markerWidth="6"
          markerHeight="6"
          orient="auto-start-reverse"
        >
          <path d="M 0 0 L 10 5 L 0 10 z" fill="#ef4444" opacity="0.5" />
        </marker>
        <!-- Parent hierarchy arrowhead -->
        <marker
          id="arrowhead-parent"
          viewBox="0 0 10 10"
          refX="10"
          refY="5"
          markerWidth="5"
          markerHeight="5"
          orient="auto-start-reverse"
        >
          <path d="M 0 0 L 10 5 L 0 10 z" fill="#6366f1" opacity="0.4" />
        </marker>
        <!-- Glow filter for hover -->
        <filter id="glow">
          <feGaussianBlur stdDeviation="3" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
        <!-- Red shadow filter for blocked edges -->
        <filter id="blocked-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feGaussianBlur in="SourceGraphic" stdDeviation="3" result="blur" />
          <feFlood flood-color="#ef4444" flood-opacity="0.3" result="color" />
          <feComposite in="color" in2="blur" operator="in" result="shadow" />
          <feMerge>
            <feMergeNode in="shadow" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
        <!-- ADO node dashed stroke pattern -->
        <filter id="ado-glow" x="-30%" y="-30%" width="160%" height="160%">
          <feGaussianBlur stdDeviation="2" result="blur" />
          <feFlood flood-color="#6366f1" flood-opacity="0.15" result="color" />
          <feComposite in="color" in2="blur" operator="in" result="glow" />
          <feMerge>
            <feMergeNode in="glow" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>

      <g :transform="transformStr">
        <!-- Background rect for click detection -->
        <rect
          class="graph-bg"
          :width="width * 5"
          :height="height * 5"
          :x="-width * 2.5"
          :y="-height * 2.5"
          fill="transparent"
        />

        <!-- Project cluster hulls -->
        <path
          v-for="hull in projectHulls"
          :key="'hull-' + hull.projectId"
          :d="hull.path"
          :fill="hull.color"
          fill-opacity="0.05"
          :stroke="hull.color"
          stroke-opacity="0.1"
          stroke-width="1"
        />

        <!-- Edges -->
        <!-- Blocked edges: dashed + red shadow -->
        <line
          v-for="(edge, i) in simulatedEdges.filter(e => isBlockedEdge(e) && !isParentEdge(e))"
          :key="'edge-blocked-' + i"
          :x1="edge.source.x"
          :y1="edge.source.y"
          :x2="edge.target.x"
          :y2="edge.target.y"
          stroke="#ef4444"
          stroke-opacity="0.5"
          stroke-width="1.5"
          stroke-dasharray="6 4"
          filter="url(#blocked-shadow)"
          marker-end="url(#arrowhead-blocked)"
        />
        <!-- Normal dependency edges: animated flowing -->
        <line
          v-for="(edge, i) in simulatedEdges.filter(e => !isBlockedEdge(e) && !isParentEdge(e))"
          :key="'edge-flow-' + i"
          class="edge-flowing"
          :x1="edge.source.x"
          :y1="edge.source.y"
          :x2="edge.target.x"
          :y2="edge.target.y"
          stroke="#3f3f46"
          stroke-opacity="0.6"
          stroke-width="1"
          stroke-dasharray="8 4"
          marker-end="url(#arrowhead)"
        />
        <!-- Parent hierarchy edges: dotted, lighter -->
        <line
          v-for="(edge, i) in simulatedEdges.filter(e => isParentEdge(e))"
          :key="'edge-parent-' + i"
          :x1="edge.source.x"
          :y1="edge.source.y"
          :x2="edge.target.x"
          :y2="edge.target.y"
          stroke="#6366f1"
          stroke-opacity="0.3"
          stroke-width="1"
          stroke-dasharray="2 3"
          marker-end="url(#arrowhead-parent)"
        />

        <!-- Nodes -->
        <g
          v-for="node in simulatedNodes"
          :key="'node-' + node.id"
          class="cursor-pointer"
        >
          <!-- Local node: solid fill circle -->
          <circle
            v-if="!isAdoNode(node)"
            class="node-circle"
            :cx="node.x"
            :cy="node.y"
            :r="nodeRadius(node.priority)"
            :fill="statusColor(node.status)"
            :stroke="selectedId === node.id ? '#ffffff' : 'transparent'"
            :stroke-width="selectedId === node.id ? 2 : 0"
            :transform="
              selectedId === node.id
                ? `translate(${node.x},${node.y}) scale(1.15) translate(${-(node.x ?? 0)},${-(node.y ?? 0)})`
                : undefined
            "
            :data-node-id="node.id"
          />
          <!-- ADO-only node: type-colored border + icon -->
          <g v-else :filter="selectedId === node.id ? 'url(#ado-glow)' : undefined">
            <rect
              class="node-circle"
              :x="(node.x ?? 0) - nodeRadius(node.priority)"
              :y="(node.y ?? 0) - nodeRadius(node.priority)"
              :width="nodeRadius(node.priority) * 2"
              :height="nodeRadius(node.priority) * 2"
              rx="3"
              :fill="adoTypeHex(node.type)"
              fill-opacity="0.15"
              :stroke="selectedId === node.id ? '#ffffff' : adoTypeHex(node.type)"
              :stroke-width="selectedId === node.id ? 2 : 1.5"
              stroke-dasharray="3 2"
              :data-node-id="node.id"
              :transform="
                selectedId === node.id
                  ? `translate(${node.x},${node.y}) scale(1.15) translate(${-(node.x ?? 0)},${-(node.y ?? 0)})`
                  : undefined
              "
            />
            <!-- ADO type icon via foreignObject -->
            <foreignObject
              v-if="nodeRadius(node.priority) >= 9"
              :x="(node.x ?? 0) - 6"
              :y="(node.y ?? 0) - 6"
              width="12"
              height="12"
              class="pointer-events-none overflow-visible"
            >
              <component
                :is="adoTypeIcon(node.type)"
                :size="11"
                :style="{ color: adoTypeHex(node.type) }"
              />
            </foreignObject>
          </g>
          <!-- Label -->
          <text
            v-if="showLabels"
            :x="(node.x ?? 0) + nodeRadius(node.priority) + 4"
            :y="(node.y ?? 0) + 3"
            class="select-none pointer-events-none"
            :fill="isAdoNode(node) ? '#6366f1' : 'currentColor'"
            :font-size="isAdoNode(node) ? 9 : 10"
            :font-style="isAdoNode(node) ? 'italic' : 'normal'"
          >
            {{ truncate(node.title) }}
          </text>
        </g>
      </g>
    </svg>
  </div>
</template>

<style scoped>
.node-circle:hover {
  filter: url(#glow);
}

/* Animated flowing effect for normal dependency edges */
.edge-flowing {
  animation: edge-flow 1.5s linear infinite;
}

@keyframes edge-flow {
  from {
    stroke-dashoffset: 24;
  }
  to {
    stroke-dashoffset: 0;
  }
}
</style>
