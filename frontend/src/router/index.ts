import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      redirect: '/tasks',
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
      path: '/playground/tasks',
      name: 'playground-tasks',
      component: () => import('../views/PlaygroundTasksView.vue'),
    },
    {
      path: '/playground/dashboard',
      name: 'playground-dashboard',
      component: () => import('../views/PlaygroundDashboardView.vue'),
    },
    {
      path: '/playground/detail',
      name: 'playground-detail',
      component: () => import('../views/PlaygroundDetailView.vue'),
    },
    {
      path: '/playground/ado',
      name: 'playground-ado',
      component: () => import('../views/PlaygroundAdoView.vue'),
    },
    {
      path: '/playground/chain',
      name: 'playground-chain',
      component: () => import('../views/PlaygroundChainView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { hideShell: true },
    },
  ],
})

export default router
