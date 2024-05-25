import { createRouter, createWebHistory } from 'vue-router'

const Home = () => import('../views/Home.vue')
const Login = () => import('../views/Login.vue')
const SignUp = () => import('../views/SignUp.vue')

const routes = [
  {
    path: '/',
    component: Home,
    meta: {
      title: 'Calendar Reminder'
    }
  },
  {
    path: '/login',
    component: Login,
    meta: {
      title: 'Calendar Reminder | Login'
    }
  },
  {
    path: '/sign-up',
    component: SignUp,
    meta: {
      title: 'Calendar Reminder | Sign Up'
    }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  window.document.title = to.meta.title
  next()
})

export default router