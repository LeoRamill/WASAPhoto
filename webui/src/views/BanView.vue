<script>

import ErrorMsg from '../components/ErrorMsg.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'

export default {

    // Including some components that will be used in the page.
    components: {
        LoadingSpinner,
        ErrorMsg,
    },

	data: function () {
		return {
			message:"",
			errormsg: "",



			loading: false,
			bans: [],
			us: this.$user_state.username,
			sameUser : false,

            username: localStorage.getItem('Username'),
            BearerToken: localStorage.getItem('BearerToken'),

		}
	},

	methods: {
		async initialize() {
			this.$user_state.current_view = this.$views.BAN;
			this.$user_state.username = this.username
			this.$user_state.headers.Authorization = this.BearerToken
		},

		async loadBan() {
			
			// Re-initializing variables to their default value.
			this.errormsg = "";
			this.loading = true;
			/*
			this.$user_state.username = this.username
			this.$user_state.headers.Authorization = this.BearerToken
			*/
			this.$route.params.username = this.$user_state.username
			this.$user_state.current_view = this.$views.BAN;
			console.log(this.$user_state.username);
			
			try {
				let response = await this.$axios.get("/users/" + this.$user_state.username + "/profile/banned", {
					headers: {
						"Content-Type": "application/json",
						"Authorization": this.$user_state.headers.Authorization
					}
				}
				)

				if (response.data != null){
					this.bans = response.data
					console.log(this.bans)
				}

				if (response.status===204){
					this.message = "You haven't banned anyone"
				}
				
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Request error, please Login before doing some action or ask to add a comment to a valid photo/comment/user." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 401) {
                    this.errormsg = "You are not allowed to do this action because you are not the stream's owner." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later." + e.toString();
					this.message = this.errormsg;
                }
			}
			this.loading = false;
		},

		async goToProfileBan(identifier, nickname){
			this.$searched_state.username = nickname;
			this.$searched_state.headers.Authorization = identifier;
			this.$router.push({ path: `/users/${nickname}/profile` })
		},
	},


	async mounted() {
		this.initialize()
		this.loadBan()
	}

}
</script>



<template>
	<div class="container-fluid">
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="container-fluid mt-3">
			<div class="row ">
				<div class="col-12 d-flex justify-content-center">
					<h1> {{ this.$user_state.username }}'s Bans</h1>
				</div>
			</div>
		</div>

        <div>
            <LoadingSpinner v-if="loading"></LoadingSpinner>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>		

		<div class="row">
			<UserIdent v-for="(ban,index) in bans" 
			:key="index" 
			:identifier="ban['user-id']['identifier']"
			:nickname="ban['username']"

            @clickedUser="goToProfileBan(ban['user-id']['identifier'],ban['username'])"
			/>
		</div>
		
		<div v-if="bans.length === 0" class="row ">
			<h2 class="d-flex justify-content-center mt-5" style="color: black;">There's no content yet, follow somebody!</h2>
		</div>
	</div>
	<div v-if="this.message!='' " class="h-25">
			<h2 class="d-flex justify-content-center mt-5" style="color: black;">{{this.message }}</h2>
    </div>
</template>

<!-- Declaration of the style(scoped) to use. -->
<style>
</style>