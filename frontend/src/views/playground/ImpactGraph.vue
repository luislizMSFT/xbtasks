<script setup lang="ts">
/**
 * ImpactGraph — Structured SVG-based DAG showing how work items
 * relate to each other (dependencies, parents, PR links, triggers).
 *
 * Uses a layered horizontal layout computed via topological sort.
 * Supports pan/zoom and highlights upstream/downstream chains on selection.
 */
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { ZoomIn, ZoomOut, Maximize } from 'lucide-vue-next'

// ── Types ──

interface ImpactNode {
  id: number
  title: string
  status: string
  priority: string
  type?: string
  adoId?: string
}

interface ImpactEdge {
  source: number
  target: number
  kind: 'dependency' | 'parent' | 'pr_link' | 'triggers'
}

const props = defineProps<{
  nodes: ImpactNode[]
  edges: ImpactEdge[]
  selectedId: number | null
}>()

const emit = defineEmits<{
  select: [id: number]
}>()

// ── Constants ──

const NODE_WIDTH = 240
const NODE_HEIGHT = 54
const LAYER_SPACING = 300
const VERTICAL_SPACING = 74
const PADDING = 60

const STATUS_FILL: Record<string, string> = {
  done: '#10b981',
  in_progress: '#3b82f6',
  in_review: '#8b5cf6',
  blocked: '#ef4444',
  todo: '#71717a',
  cancelled: '#a1a1aa',
}

// ADO type → color. Used when node.type is an ADO type.
const TYPE_FILL: Record<string, string> = {
  epic: '#f97316',
  feature: '#22c55e',
  'user story': '#a855f7',
  task: '#f59e0b',
  bug: '#ef4444',
  pr: '#3b82f6',
  personal: '#71717a',
}

const EDGE_COLORS: Record<string, string> = {
  dependency: '#9ca3af',
  parent: '#6366f1',
  pr_link: '#22c55e',
  triggers: '#f59e0b',
}

const EDGE_DASH: Record<string, string> = {
  dependency: '',
  parent: '6 4',
  pr_link: '2 3',
  triggers: '',
}

// ── Layout computation ──

interface PositionedNode extends ImpactNode {
  x: number
  y: number
  layer: number
}

const layout = computed(() => {
  const nodeList = props.nodes
  const edgeList = props.edges

  if (nodeList.length === 0) return { nodes: [] as PositionedNode[], edges: edgeList }

  // Build adjacency: source -> targets (outgoing) and target -> sources (incoming)
  const incoming = new Map<number, Set<number>>()
  const outgoing = new Map<number, Set<number>>()
  for (const n of nodeList) {
    incoming.set(n.id, new Set())
    outgoing.set(n.id, new Set())
  }
  const nodeIds = new Set(nodeList.map(n => n.id))
  for (const e of edgeList) {
    if (!nodeIds.has(e.source) || !nodeIds.has(e.target)) continue
    incoming.get(e.target)!.add(e.source)
    outgoing.get(e.source)!.add(e.target)
  }

  // Topological layering (Kahn's algorithm variant with cycle handling)
  const layers = new Map<number, number>()
  const inDegree = new Map<number, number>()
  for (const n of nodeList) {
    inDegree.set(n.id, incoming.get(n.id)!.size)
  }

  // Initial queue: nodes with no incoming edges
  let queue: number[] = []
  for (const n of nodeList) {
    if (inDegree.get(n.id) === 0) {
      queue.push(n.id)
      layers.set(n.id, 0)
    }
  }

  const visited = new Set<number>()
  while (queue.length > 0) {
    const nextQueue: number[] = []
    for (const id of queue) {
      visited.add(id)
      const layer = layers.get(id)!
      for (const target of outgoing.get(id) ?? []) {
        const candidateLayer = layer + 1
        if (!layers.has(target) || layers.get(target)! < candidateLayer) {
          layers.set(target, candidateLayer)
        }
        inDegree.set(target, inDegree.get(target)! - 1)
        if (inDegree.get(target) === 0 && !visited.has(target)) {
          nextQueue.push(target)
        }
      }
    }
    queue = nextQueue
  }

  // Handle cycles: assign unvisited nodes to max layer of their deps + 1
  for (const n of nodeList) {
    if (!visited.has(n.id)) {
      let maxDepLayer = -1
      for (const src of incoming.get(n.id) ?? []) {
        if (layers.has(src)) {
          maxDepLayer = Math.max(maxDepLayer, layers.get(src)!)
        }
      }
      layers.set(n.id, maxDepLayer + 1)
    }
  }

  // Group nodes by layer
  const layerGroups = new Map<number, ImpactNode[]>()
  for (const n of nodeList) {
    const l = layers.get(n.id) ?? 0
    if (!layerGroups.has(l)) layerGroups.set(l, [])
    layerGroups.get(l)!.push(n)
  }

  // Sort layers
  const sortedLayers = Array.from(layerGroups.keys()).sort((a, b) => a - b)

  // Barycenter ordering: reduce edge crossings by sorting nodes within each
  // layer by the average y-position of their connected nodes in previous layers
  const tempY = new Map<number, number>()
  for (const layerIdx of sortedLayers) {
    const group = layerGroups.get(layerIdx)!
    if (layerIdx > 0) {
      // Sort by barycenter of incoming neighbors
      group.sort((a, b) => {
        const aNeighbors = (incoming.get(a.id) ?? new Set())
        const bNeighbors = (incoming.get(b.id) ?? new Set())
        const aAvg = aNeighbors.size > 0
          ? [...aNeighbors].reduce((s, id) => s + (tempY.get(id) ?? 0), 0) / aNeighbors.size
          : 0
        const bAvg = bNeighbors.size > 0
          ? [...bNeighbors].reduce((s, id) => s + (tempY.get(id) ?? 0), 0) / bNeighbors.size
          : 0
        return aAvg - bAvg
      })
    }
    const totalHeight = group.length * NODE_HEIGHT + (group.length - 1) * (VERTICAL_SPACING - NODE_HEIGHT)
    const startY = -totalHeight / 2
    for (let i = 0; i < group.length; i++) {
      tempY.set(group[i].id, startY + i * VERTICAL_SPACING + NODE_HEIGHT / 2)
    }
  }

  // Position nodes using computed ordering
  const positioned: PositionedNode[] = []
  for (const layerIdx of sortedLayers) {
    const group = layerGroups.get(layerIdx)!
    const totalHeight = group.length * NODE_HEIGHT + (group.length - 1) * (VERTICAL_SPACING - NODE_HEIGHT)
    const startY = -totalHeight / 2

    for (let i = 0; i < group.length; i++) {
      const n = group[i]
      positioned.push({
        ...n,
        x: PADDING + layerIdx * LAYER_SPACING,
        y: startY + i * VERTICAL_SPACING,
        layer: layerIdx,
      })
    }
  }

  return { nodes: positioned, edges: edgeList }
})

// ── Selection highlight chains ──

const highlightedIds = computed<Set<number>>(() => {
  if (props.selectedId == null) return new Set()

  const sel = props.selectedId
  const result = new Set<number>([sel])

  const nodeIds = new Set(props.nodes.map(n => n.id))
  if (!nodeIds.has(sel)) return result

  // Build adjacency for traversal
  const incoming = new Map<number, number[]>()
  const outgoing = new Map<number, number[]>()
  for (const n of props.nodes) {
    incoming.set(n.id, [])
    outgoing.set(n.id, [])
  }
  for (const e of props.edges) {
    if (!nodeIds.has(e.source) || !nodeIds.has(e.target)) continue
    outgoing.get(e.source)!.push(e.target)
    incoming.get(e.target)!.push(e.source)
  }

  // Trace upstream (what it depends on)
  const upQueue = [...(incoming.get(sel) ?? [])]
  while (upQueue.length > 0) {
    const id = upQueue.pop()!
    if (result.has(id)) continue
    result.add(id)
    for (const src of incoming.get(id) ?? []) {
      if (!result.has(src)) upQueue.push(src)
    }
  }

  // Trace downstream (what depends on it)
  const downQueue = [...(outgoing.get(sel) ?? [])]
  while (downQueue.length > 0) {
    const id = downQueue.pop()!
    if (result.has(id)) continue
    result.add(id)
    for (const tgt of outgoing.get(id) ?? []) {
      if (!result.has(tgt)) downQueue.push(tgt)
    }
  }

  return result
})

const highlightedEdges = computed<Set<string>>(() => {
  if (props.selectedId == null) return new Set()
  const ids = highlightedIds.value
  const result = new Set<string>()
  for (const e of props.edges) {
    if (ids.has(e.source) && ids.has(e.target)) {
      result.add(`${e.source}-${e.target}`)
    }
  }
  return result
})

function isNodeHighlighted(id: number): boolean {
  if (props.selectedId == null) return true
  return highlightedIds.value.has(id)
}

function isEdgeHighlighted(e: ImpactEdge): boolean {
  if (props.selectedId == null) return true
  return highlightedEdges.value.has(`${e.source}-${e.target}`)
}

// ── Node position lookup ──

const nodePositions = computed(() => {
  const map = new Map<number, { x: number; y: number }>()
  for (const n of layout.value.nodes) {
    map.set(n.id, { x: n.x, y: n.y })
  }
  return map
})

// ── Edge path computation ──

function edgePath(e: ImpactEdge): string {
  const src = nodePositions.value.get(e.source)
  const tgt = nodePositions.value.get(e.target)
  if (!src || !tgt) return ''

  const x1 = src.x + NODE_WIDTH
  const y1 = src.y + NODE_HEIGHT / 2
  const x2 = tgt.x
  const y2 = tgt.y + NODE_HEIGHT / 2

  const dx = Math.abs(x2 - x1) * 0.5
  return `M ${x1} ${y1} C ${x1 + dx} ${y1}, ${x2 - dx} ${y2}, ${x2} ${y2}`
}

// ── SVG viewBox ──

const viewBox = computed(() => {
  const nodes = layout.value.nodes
  if (nodes.length === 0) return '0 0 600 400'

  let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity
  for (const n of nodes) {
    minX = Math.min(minX, n.x)
    minY = Math.min(minY, n.y)
    maxX = Math.max(maxX, n.x + NODE_WIDTH)
    maxY = Math.max(maxY, n.y + NODE_HEIGHT)
  }

  const pad = PADDING
  return `${minX - pad} ${minY - pad} ${maxX - minX + pad * 2} ${maxY - minY + pad * 2}`
})

// ── Pan and zoom ──

const containerRef = ref<HTMLDivElement | null>(null)
const svgRef = ref<SVGSVGElement | null>(null)
const containerWidth = ref(800)
const containerHeight = ref(600)

const panX = ref(0)
const panY = ref(0)
const scale = ref(1)
const isPanning = ref(false)
const panStart = ref({ x: 0, y: 0 })

let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  if (containerRef.value) {
    const rect = containerRef.value.getBoundingClientRect()
    containerWidth.value = rect.width
    containerHeight.value = rect.height

    resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        containerWidth.value = entry.contentRect.width
        containerHeight.value = entry.contentRect.height
      }
    })
    resizeObserver.observe(containerRef.value)
  }
})

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
})

function onWheel(event: WheelEvent) {
  event.preventDefault()
  const delta = event.deltaY > 0 ? -0.1 : 0.1
  scale.value = Math.min(2.0, Math.max(0.3, scale.value + delta))
}

function onPointerDown(event: PointerEvent) {
  // Only pan on background clicks (not on nodes)
  if ((event.target as SVGElement).closest('.impact-node')) return
  isPanning.value = true
  panStart.value = { x: event.clientX - panX.value, y: event.clientY - panY.value }
  ;(event.currentTarget as Element)?.setPointerCapture(event.pointerId)
}

function onPointerMove(event: PointerEvent) {
  if (!isPanning.value) return
  panX.value = event.clientX - panStart.value.x
  panY.value = event.clientY - panStart.value.y
}

function onPointerUp(event: PointerEvent) {
  isPanning.value = false
  ;(event.currentTarget as Element)?.releasePointerCapture(event.pointerId)
}

function resetView() {
  panX.value = 0
  panY.value = 0
  scale.value = 1
}

function zoomIn() {
  scale.value = Math.min(2.0, scale.value + 0.15)
}

function zoomOut() {
  scale.value = Math.max(0.3, scale.value - 0.15)
}

// ── Helpers ──

function truncate(text: string, max: number): string {
  return text.length > max ? text.slice(0, max - 1) + '\u2026' : text
}

// Auto-fit view on mount
watch(
  () => [containerWidth.value, containerHeight.value, layout.value] as const,
  () => {
    if (scale.value !== 1 || panX.value !== 0 || panY.value !== 0) return // user has interacted
    autoFit()
  },
  { flush: 'post' },
)

function autoFit() {
  const nodes = layout.value.nodes
  if (nodes.length === 0) return

  let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity
  for (const n of nodes) {
    minX = Math.min(minX, n.x)
    minY = Math.min(minY, n.y)
    maxX = Math.max(maxX, n.x + NODE_WIDTH)
    maxY = Math.max(maxY, n.y + NODE_HEIGHT)
  }

  const graphW = maxX - minX + PADDING * 2
  const graphH = maxY - minY + PADDING * 2
  const scaleX = containerWidth.value / graphW
  const scaleY = containerHeight.value / graphH
  const fitScale = Math.min(scaleX, scaleY, 1.0) * 0.9 // 90% to leave margin
  scale.value = Math.max(0.3, fitScale)
}

function statusFill(status: string): string {
  return STATUS_FILL[status] ?? STATUS_FILL.todo
}

function nodeFill(node: ImpactNode): string {
  if (node.type) {
    const typeFill = TYPE_FILL[node.type]
    if (typeFill) return typeFill
  }
  return statusFill(node.status)
}
</script>

<template>
  <div class="relative w-full h-full flex flex-col bg-background text-foreground">
    <!-- Toolbar -->
    <div
      class="absolute top-2 right-2 z-10 flex items-center gap-1 rounded-md border border-border bg-background/80 backdrop-blur px-1.5 py-1"
    >
      <button
        class="p-1 rounded hover:bg-muted transition-colors text-muted-foreground hover:text-foreground"
        title="Zoom in"
        @click="zoomIn"
      >
        <ZoomIn :size="14" />
      </button>
      <button
        class="p-1 rounded hover:bg-muted transition-colors text-muted-foreground hover:text-foreground"
        title="Zoom out"
        @click="zoomOut"
      >
        <ZoomOut :size="14" />
      </button>
      <button
        class="p-1 rounded hover:bg-muted transition-colors text-muted-foreground hover:text-foreground"
        title="Reset view"
        @click="resetView"
      >
        <Maximize :size="14" />
      </button>
      <span class="text-[10px] text-muted-foreground tabular-nums px-1">
        {{ Math.round(scale * 100) }}%
      </span>
    </div>

    <!-- Edge legend -->
    <div
      class="absolute bottom-2 left-2 z-10 flex items-center gap-3 rounded-md border border-border bg-background/80 backdrop-blur px-2 py-1 text-[10px] text-muted-foreground"
    >
      <span class="flex items-center gap-1">
        <svg width="20" height="6"><line x1="0" y1="3" x2="20" y2="3" stroke="#9ca3af" stroke-width="1.5" /></svg>
        dependency
      </span>
      <span class="flex items-center gap-1">
        <svg width="20" height="6"><line x1="0" y1="3" x2="20" y2="3" stroke="#6366f1" stroke-width="1.5" stroke-dasharray="6 4" /></svg>
        parent
      </span>
      <span class="flex items-center gap-1">
        <svg width="20" height="6"><line x1="0" y1="3" x2="20" y2="3" stroke="#22c55e" stroke-width="1.5" stroke-dasharray="2 3" /></svg>
        pr_link
      </span>
      <span class="flex items-center gap-1">
        <svg width="20" height="6"><line x1="0" y1="3" x2="20" y2="3" stroke="#f59e0b" stroke-width="1.5" /></svg>
        triggers
      </span>
    </div>

    <!-- SVG canvas -->
    <div
      ref="containerRef"
      class="flex-1 min-h-0 overflow-hidden cursor-grab"
      :class="{ 'cursor-grabbing': isPanning }"
      @wheel.prevent="onWheel"
      @pointerdown="onPointerDown"
      @pointermove="onPointerMove"
      @pointerup="onPointerUp"
    >
      <svg
        ref="svgRef"
        :width="containerWidth"
        :height="containerHeight"
        :viewBox="viewBox"
        class="select-none"
      >
        <g :transform="`translate(${panX / scale}, ${panY / scale}) scale(${scale})`">
          <!-- Arrowhead markers -->
          <defs>
            <marker
              v-for="kind in (['dependency', 'parent', 'pr_link', 'triggers'] as const)"
              :key="kind"
              :id="`arrow-${kind}`"
              markerWidth="8"
              markerHeight="6"
              refX="8"
              refY="3"
              orient="auto"
            >
              <polygon points="0 0, 8 3, 0 6" :fill="EDGE_COLORS[kind]" />
            </marker>
          </defs>

          <!-- Edges (rendered first, behind nodes) -->
          <path
            v-for="(edge, idx) in layout.edges"
            :key="`e-${idx}`"
            :d="edgePath(edge)"
            fill="none"
            :stroke="EDGE_COLORS[edge.kind] ?? '#9ca3af'"
            :stroke-width="isEdgeHighlighted(edge) ? 2 : 1.2"
            :stroke-dasharray="EDGE_DASH[edge.kind] ?? ''"
            :opacity="isEdgeHighlighted(edge) ? 1 : 0.15"
            :marker-end="`url(#arrow-${edge.kind})`"
            class="transition-opacity duration-200"
          />

          <!-- Nodes -->
          <g
            v-for="node in layout.nodes"
            :key="node.id"
            class="impact-node cursor-pointer"
            :transform="`translate(${node.x}, ${node.y})`"
            :opacity="isNodeHighlighted(node.id) ? 1 : 0.3"
            @click.stop="emit('select', node.id)"
          >
            <!-- Background rect -->
            <rect
              :width="NODE_WIDTH"
              :height="NODE_HEIGHT"
              rx="6"
              ry="6"
              :fill="nodeFill(node)"
              fill-opacity="0.1"
              :stroke="nodeFill(node)"
              :stroke-width="selectedId === node.id ? 2.5 : 1"
              :stroke-opacity="selectedId === node.id ? 1 : 0.4"
              class="transition-all duration-150"
            />
            <!-- Left color bar -->
            <rect
              x="0"
              y="0"
              width="3"
              :height="NODE_HEIGHT"
              rx="1.5"
              :fill="nodeFill(node)"
            />
            <!-- Selected scale effect via a slightly larger outer rect -->
            <rect
              v-if="selectedId === node.id"
              x="-3"
              y="-3"
              :width="NODE_WIDTH + 6"
              :height="NODE_HEIGHT + 6"
              rx="8"
              ry="8"
              fill="none"
              :stroke="nodeFill(node)"
              stroke-width="1.5"
              stroke-opacity="0.5"
            />
            <!-- Title -->
            <text
              :x="10"
              :y="20"
              font-size="12"
              font-weight="500"
              fill="#e5e5e5"
              class="dark:fill-neutral-200"
            >
              {{ truncate(node.title, 30) }}
            </text>
            <!-- Priority + ADO ID -->
            <text
              :x="10"
              :y="38"
              font-size="10"
              fill="#a3a3a3"
            >
              <tspan font-weight="600" :fill="nodeFill(node)">{{ node.priority }}</tspan>
              <tspan v-if="node.type" dx="6" :fill="nodeFill(node)" font-weight="500" class="uppercase">{{ node.type }}</tspan>
              <tspan v-if="node.adoId" dx="6" fill="#60a5fa">#{{ node.adoId }}</tspan>
            </text>
            <!-- Hover rect (invisible, for hover effect) -->
            <rect
              :width="NODE_WIDTH"
              :height="NODE_HEIGHT"
              rx="6"
              ry="6"
              fill="white"
              fill-opacity="0"
              class="hover:fill-opacity-[0.05] transition-all duration-100"
            />
          </g>
        </g>
      </svg>
    </div>
  </div>
</template>
