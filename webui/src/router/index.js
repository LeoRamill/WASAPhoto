import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import StreamView from '../views/StreamView.vue'
import SearchView from '../views/SearchView.vue'
import ProfileView from '../views/ProfileView.vue'
import BanView from '../views/BanView.vue'
import UpdateUsname from '../views/UpdateUsname.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			redirect: '/session'
		},
		{
			path: '/session',
			component: LoginView
		},
		{
			path: '/users/:username/homepage',
			component: StreamView
		},
		{
			path: '/users/:username/',
			component: SearchView
		},
		{
			path: '/users/:username/profile',
			component: ProfileView
		},
		{
			path: '/users/:username/profile/banned',
			component: BanView
		},
		{
			path: '/users/:username/settings',
			component: UpdateUsname
		}

	]
})

export default router
