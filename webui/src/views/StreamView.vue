<script>
import ErrorMsg from '../components/ErrorMsg.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'

export default {
	data: function () {
		return {
			errormsg: "",
			loading: false,
			photos: [],
			us: this.$user_state.username,
			sameUser : false,
			message: "",
			// Retrieving from the Cache the Username and the Bearer Authenticaiton Token.
            username: localStorage.getItem('Username'),
            BearerToken: localStorage.getItem('BearerToken'),
		}
	},

	methods: {

		async initialize() {
			this.$user_state.current_view = this.$views.STREAM;
			this.$user_state.username = this.username;
			this.$user_state.headers.Authorization = this.BearerToken

		},


		async loadStream() {
			
			// Re-initializing variables to their default value.
			this.errormsg = "";
			this.loading = true;

			/*
			this.$user_state.username = this.username
			this.$user_state.headers.Authorization = this.BearerToken
			*/
			this.$route.params.username = this.$user_state.username

			this.$user_state.current_view = this.$views.STREAM;
			console.log(this.$user_state.username);
			console.log(this.username)
			console.log(this.BearerToken)
			
			try {
				let response = await this.$axios.get("/users/" + this.$user_state.username + "/homepage", {
					headers: {
						"Content-Type": "application/json",
						"Authorization": this.$user_state.headers.Authorization
					}
				}
				)

				if (response.data != null){
					this.photos = response.data
				}

				if (response.status===204){
					this.message = "There's no content yet, follow somebody or no one posted yet"
				}
				
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Request error, please Login before doing some action or ask to get the list of photo of the stream of a valid user." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 401) {
                    this.errormsg = "You are not allowed to do this action because you are not the stream's owner." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 204) {
                    this.errormsg = "In the Internal DB there is not anymore the content you have asked." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later." + e.toString();
					this.message = this.errormsg;
                }else if (e.response && e.response.status === 404) {
                    this.errormsg = "No Photo found" + e.toString();
					this.message = this.errormsg;
				}else {
                    this.errormsg = "Please Login before with an Authorized profile to view this page. " + e.toString();
					this.message = this.errormsg;
                }
			}
			this.loading = false;
		}
	},

	async mounted() {
		this.initialize()
		this.loadStream()
	}

}
</script>

<template>
	<div class="container-fluid">

		<div class="container-fluid mt-3">
			<div class="row ">
				<div class="col-12 d-flex justify-content-center">
					<h1> {{ this.$user_state.username }}'s Stream</h1>
				</div>
			</div>
		</div>

		<div class="row">
			<Photo v-for="(photo,index) in photos" 
                    :key="index" 
                    :owner="photo['nickname']" 
					:owner_id="photo['user-id']['identifier']"
                    :photo_id="photo['photo-id']['identifier']" 
                    :comments="photo['comment-collection']" 
                    :likes="photo['like-collection']" 
                    :upload_date="photo['date-time']" 
                    :isOwner="sameUser" 
                    />
					{{ photos['like-collection'] }}
		</div>

	</div>

    <div v-if="this.message!='' " class="h-25">
			<h2 class="d-flex justify-content-center mt-5" style="color: black;">{{this.message }}</h2>
    </div>

	<div>
        <LoadingSpinner v-if="loading"></LoadingSpinner>
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
</template>

<style>
</style>