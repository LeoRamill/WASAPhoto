
<script setup>

import {  RouterView } from 'vue-router'



</script>

<script>
export default {
	data: function(){
			return{
			// Retrieving from the Cache the Username and the Bearer Authenticaiton Token.
			usernameLogged: localStorage.getItem('Username'),
			username: localStorage.getItem('Username') == localStorage.getItem('usernameProfileToView') ? localStorage.getItem('Username') : localStorage.getItem('usernameProfileToView'),
            BearerToken: localStorage.getItem('BearerToken'),
			}
		},
	methods: {
		
		async goToLogin(){
			localStorage.clear();
			this.$user_state.username = null;
			this.$user_state.headers.Authorization = null;
            this.$searched_state.username = null;
			this.$searched_state.headers.Authorization = null;
			localStorage.clear();
			localStorage.setItem('usernameProfileToView', localStorage.getItem('Username'));
			this.$router.push({path: `/session`})
		},
		async goToMyStream() {
			if (this.$user_state.username  == null) {
				return
			}
			localStorage.setItem('usernameProfileToView', localStorage.getItem('Username'));
			this.$router.push({ path: "/users/" + this.$user_state.username  + "/homepage" })
		},

		async goToSearch() {
			if (this.$user_state.username == null) {
				return
			}
			localStorage.setItem('usernameProfileToView', localStorage.getItem('Username'));
			this.$router.push({ path: `/users/${this.$user_state.username}/` })
		},

		async goToMyProfile() {

			if (this.$user_state.username  == null) {
				return
			}
			localStorage.setItem('usernameProfileToView', localStorage.getItem('Username'));
			this.$router.push({ path: `/users/${this.$user_state.username}/profile` })
		},

		async goToBan() {
			if (this.$user_state.username  == null) {
				return
			}
			localStorage.setItem('usernameProfileToView', localStorage.getItem('Username'));
			this.$router.push({ path: `/users/${this.$user_state.username}/profile/banned` })
		},

		async goToUpdateProfile() {
			if (this.$user_state.username  == null) {
				return
			}
			localStorage.setItem('usernameProfileToView', localStorage.getItem('Username'));
			this.$router.push({ path: `/users/${this.$user_state.username}/settings` })
		},

	},

}

</script>


<template>
    <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
        <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-7">WASA Photo</a>
    </header>

    <div class="container-fluid">
        <div class="row">
            <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
                <div class="position-sticky pt-3 sidebar-sticky">
                    <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
                        <span>Menu</span>
                    </h6>
                    <ul class="nav flex-column">
                        <li class="nav-item" @click="goToLogin" style="margin-left: 10px; margin-top: 10px; margin-right: 10px;">
                            <!-- <RouterLink to="/session/" class="nav-link"> -->
                                <svg class="feather" style="color:#4a77d4; margin-right: 10px; "><use href="/feather-sprite-v4.29.0.svg#log-in"/></svg>
                                <b>Logout</b>
                            <!-- </RouterLink> -->
                        </li>
                    </ul>

                    <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-5 mb-1 text-muted text-uppercase">
                        <span>PERSONAL PROFILE</span>
                    </h6>
                    <ul class="nav flex-column">
                        <li class="nav-item" @click="goToMyProfile(usname)" style="margin-left: 10px; margin-top: 10px; margin-right: 10px;">
                            <!-- <RouterLink :to="'/users/'+username" class="nav-link" > -->
                                <svg class="feather" style="color:#4a77d4; margin-right: 10px;"><use href="/feather-sprite-v4.29.0.svg#instagram"/></svg>
                                <b>My Profile </b>
                            <!-- </RouterLink> -->
                        </li>
                        <li class="nav-item" @click="goToMyStream" style="margin-left: 10px; margin-top: 10px; margin-right: 10px;">
                            <!-- <RouterLink :to="'/users/'+username+'/myStream/'" class="nav-link" > -->
                                <svg class="feather" style="color:#4a77d4;margin-right: 10px; "><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                                <b>My Stream</b>
                            <!-- </RouterLink> -->
                        </li>
                    </ul>

                    <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-5 mb-1 text-muted text-uppercase">
                        <span>ACTIONS</span>
                    </h6>
                    <ul class="nav flex-column">
                        <li class="nav-item" @click="goToSearch" style="margin-left: 10px; margin-top: 10px; margin-right: 10px;">
                            <!-- <RouterLink to="/search/" class="nav-link"> -->
                                <svg class="feather" style="color:#4a77d4; margin-right: 10px; "><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
                                <b>Search</b>
                            <!-- </RouterLink> -->
                        </li>

						<li class="nav-item" @click="goToUpdateProfile" style="margin-left: 10px; margin-top: 10px; margin-right: 10px;">
							<!-- <RouterLink :to="'/users/'+username+'/update/'" class="nav-link" > -->
								<svg class="feather" style="color:#4a77d4;  10px;margin-right: 10px;  "><use href="/feather-sprite-v4.29.0.svg#edit"/></svg>
								<b>Update Profile</b>
							<!-- </RouterLink> -->
						</li>

						<li class="nav-item" @click="goToBan" style="margin-left: 10px; margin-top: 10px; margin-right: 10px;">
							<!-- <RouterLink :to="'/users/'+username+'/ban/'" class="nav-link" > -->
								<svg class="feather" style="color:#4a77d4;  margin-right: 10px; "><use href="/feather-sprite-v4.29.0.svg#lock"/></svg>
								<b>Ban</b>
							<!-- </RouterLink> -->
						</li>
                    </ul>

                
                </div>
            </nav>

            <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                <RouterView />
            </main>
        </div>
    </div>

</template>
