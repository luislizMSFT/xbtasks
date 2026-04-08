import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: () => import('../views/TasksView.vue'),
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
    },
    {
      path: '/projects',
      name: 'projects',
      component: () => import('../views/ProjectsView.vue'),
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/SettingsView.vue'),
    },
    {
      path: '/ado',
      name: 'ado',
      component: () => import('../views/AdoView.vue'),
    },
    {
      path: '/dependencies',
      name: 'dependencies',
      component: () => import('../views/DependencyGraphView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { hideShell: true },
    },
    // ── Playground pages (isolated testing) ──
    {
      path: '/playground/task-detail',
      name: 'playground-task-detail',
      component: () => import('../views/playground/PlaygroundTaskDetail.vue'),
      meta: { hideShell: true },
    },
    {
      path: '/playground/dashboard-header',
      name: 'playground-dashboard-header',
      component: () => import('../views/playground/PlaygroundDashboardHeader.vue'),
      meta: { hideShell: true },
    },
    {
      path: '/playground/task-tree',
      name: 'playground-task-tree',
      component: () => import('../views/playground/PlaygroundTaskTree.vue'),
      meta: { hideShell: true },
    },
    {
      path: '/playground/task-styling',
      name: 'playground-task-styling',
      component: () => import('../views/playground/PlaygroundTaskStyling.vue'),
      meta: { hideShell: true },
    },
    {
      path: '/playground/lifecycle',
      name: 'playground-lifecycle',
      component: () => import('../views/playground/PlaygroundLifecycle.vue'),
      meta: { hideShell: true },
    },
    {
      path: '/playground/dashboard',
      name: 'playground-dashboard',
      component: () => import('../views/playground/PlaygroundDashboard.vue'),
      meta: { hideShell: true },
    },
    {
      path: '/playground/integrated',
      name: 'playground-integrated',
      component: () => import('../views/playground/PlaygroundIntegrated.vue'),
      meta: { hideShell: true },
    },
  ],
})

export default router
