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
            errormsg: "",
            loading: false,
			nickname:"",

			username: localStorage.getItem('Username'),
            BearerToken: localStorage.getItem('BearerToken'),

		}
	},

	methods:{

		async initialize() {
			this.$user_state.current_view = this.$views.UPDATE;
			this.$user_state.username = this.username
			this.$user_state.headers.Authorization = this.BearerToken
		},

        async changeName(){
            // Re-initializing variables to their default value.
			this.errormsg = "";
			this.loading = true;
			this.$user_state.current_view = this.$views.UPDATE;
			/*
			this.$user_state.username = this.username
			this.$user_state.headers.Authorization = this.BearerToken
			*/
            try{
                const req_body = {
                    "username-string": this.nickname
                }
                const res = await this.$axios.put("/users/"+this.$user_state.username+"/profile", req_body, {
                                headers: {
                                    "Authorization": this.$user_state.headers.Authorization
                                }})

                if (res == undefined) {
                console.log("Error: undefined response");
                return
                }

				if (res.status==201){
					this.$user_state.username = res.data;
					console.log(res.data)
					console.log(this.nickname)
					localStorage.setItem('BearerToken', this.$user_state.headers.Authorization),
					localStorage.setItem('Username', this.$user_state.username),
					localStorage.setItem('usernameProfileToView', this.$user_state.username)
					this.$router.push({ path: `/users/${this.$user_state.username}/profile`})
					this.nickname=""
				}

            }catch(e){
				 if (e.response && e.response.status === 401) {
                    this.errormsg = "You are not allowed to do this action because you are not the stream's owner." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later." + e.toString();
					this.message = this.errormsg;
                }
            }
			this.loading = false;    
		},

		async mounted() {
		this.initialize()
	}

	}



}
</script>

<template>
	<div class="container-fluid">
		<div class="row">
			<div class="col d-flex justify-content-center mb-2">
				<h1>{{ this.$route.params.username }}'s Settings</h1>
			</div>
		</div>

		<div class="row mt-2">
			<div class="col d-flex justify-content-center">
				<div class="input-group mb-3 w-25">
					<input
						id="updateName"
						type="text"
						class="form-control w-25"
						placeholder="Your new nickname..."
						maxlength="16"
						minlength="3"
						v-model="this.nickname"
					/>
					<div class="input-group-append">
						<button class="btn btn-outline-secondary" 
						@click="changeName"
						:disabled="this.nickname === null || nickname.length >16 || nickname.length <3 || nickname.trim().length===0">
						Modify</button>
					</div>
				</div>
			</div>
		</div>

		<div class="row" >
			<div v-if="nickname.trim().length>0" class="col d-flex justify-content-center">
				Preview: {{nickname}} @{{ this.$route.params.username }}
			</div>
		</div>

        <div>
            <LoadingSpinner v-if="loading"></LoadingSpinner>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
	</div>
	
</template>

<style>
</style>