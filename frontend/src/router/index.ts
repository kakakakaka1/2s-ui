// Composables
import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import Data from '@/store/modules/data'

const routes = [
  {
    path: '/login',
    name: 'pages.login',
    component: Login,
  },
  {
    path: '/',
    component: () => import('@/layouts/default/Default.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '/',
        name: 'pages.home',
        component: () => import('@/views/Home.vue'),
      },
      {
        path: '/inbounds',
        name: 'pages.inbounds',
        component: () => import('@/views/Inbounds.vue'),
      },
      {
        path: '/clients',
        name: 'pages.clients',
        component: () => import('@/views/Clients.vue'),
      },  
      {
        path: '/outbounds',
        name: 'pages.outbounds',
        component: () => import('@/views/Outbounds.vue'),
      },
      {
        path: '/services',
        name: 'pages.services',
        component: () => import('@/views/Services.vue'),
      },
      {
        path: '/endpoints',
        name: 'pages.endpoints',
        component: () => import('@/views/Endpoints.vue'),
      },
      {
        path: '/rules',
        name: 'pages.rules',
        component: () => import('@/views/Rules.vue'),
      },
      {
        path: '/tls',
        name: 'pages.tls',
        component: () => import('@/views/Tls.vue'),
      },
      {
        path: '/basics',
        name: 'pages.basics',
        component: () => import('@/views/Basics.vue'),
      },
      {
        path: '/dns',
        name: 'pages.dns',
        component: () => import('@/views/Dns.vue'),
      },
      {
        path: '/nodes',
        name: 'pages.nodes',
        component: () => import('@/views/Nodes.vue'),
      },
      {
        path: '/admins',
        name: 'pages.admins',
        component: () => import('@/views/Admins.vue'),
      },
      {
        path: '/settings',
        name: 'pages.settings',
        component: () => import('@/views/Settings.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory((window as any).BASE_URL),
  routes,
})

const DEFAULT_TITLE = '2S-UI'
let intervalId:any

// Chunk names carry a per-build id and build.sh wipes web/html before copying, so
// after an in-place upgrade an already-open tab asks for chunks the server no longer
// has. Reload once to pick up the new index.html; the flag keeps a genuinely
// unreachable server from turning that into a reload loop.
const RELOADED_KEY = '2sui-chunk-reload'

router.onError((err) => {
  if (!/dynamically imported module|Importing a module script failed|Failed to fetch/i.test(String(err))) return
  if (sessionStorage.getItem(RELOADED_KEY)) return
  sessionStorage.setItem(RELOADED_KEY, '1')
  window.location.reload()
})

router.afterEach(() => sessionStorage.removeItem(RELOADED_KEY))

// Navigation guard to check authentication state
router.beforeEach((to) => {
  // Load default data
  if (to.path !== '/login') {
    loadDataInterval()
  } else {
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = undefined
    }
  }
})

const loadDataInterval = () => {
  if (intervalId) return
  Data().loadData()
  intervalId = setInterval(() => {
    Data().loadData()
  }, 10000)
}

export default router
