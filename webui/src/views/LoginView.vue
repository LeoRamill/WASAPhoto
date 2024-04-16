
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

			id:"",
        }
    },

    methods: {

		async initialize(){
			this.$user_state.current_view = this.$views.LOGIN;
			this.$user_state.username = null;
		},

        async DoLogin() {

            // Re-initializing variables to their default value.
			this.errormsg = "";
			this.loading = true;

			// get username from the input
            let username = document.getElementById("login-form").value;

            // check username regex

            if (!username.match("^[a-zA-Z0-9_.]{6,12}$")) {
                alert("Invalid username");
                return;
            }

			try{
				let response = await this.$axios.post("/session", {
                "username-string": username
            });

			this.$user_state.username = username;
            this.$user_state.headers.Authorization = response.data;
            /*
            localStorage.setItem('BearerToken', this.$user_state.headers.Authorization),
            localStorage.setItem('Username', this.$user_state.username),
            localStorage.setItem('usernameProfileToView', this.$user_state.username),
            */
			console.log(this.$user_state.username);
			
            this.$route.params.username = this.$user_state.username
            localStorage.setItem('BearerToken', this.$user_state.headers.Authorization),
            localStorage.setItem('Username', this.$user_state.username),
            localStorage.setItem('usernameProfileToView', this.$user_state.username),
            // Re-addressing the page to the personal profile page of a user.
            this.$router.push({ path: `/users/${this.$user_state.username}/homepage` })
			

			} catch(e){
                if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us. This may be due to an incorrect insertion. Be sure to resepct the rules!" + e.toString();
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An Internal Error occurred. We are sorry for the inconvenient. Please try again later." + + e.toString();
				} else {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
        },

		mounted(){
			this.initialize()
		}

    }
}
</script>

<template>
    <div class="container text-center pt-3 pb-2">
        <h2>Login</h2>
    </div>


    <div class="h-75 d-flex align-items-center justify-content-center">
        <form class="border border-dark p-5 rounded shadow-lg">
            <!-- Username input -->
            <div class="form-outline mb-4">
                <input type="text" id="login-form" class="form-control" pattern="^[a-zA-Z0-9_.]{6,12}$" />
                <label class="form-label" for="login-form">Username</label>
            </div>

            <!-- Submit button -->
            <button type="button" class="btn btn-primary btn-block mb-4" @click="DoLogin()">Sign in</button>

        </form>
        <div>
            <LoadingSpinner v-if="loading"></LoadingSpinner>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
    </div>

</template>

<style>
</style>