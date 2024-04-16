<script>
import ErrorMsg from '../components/ErrorMsg.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'

export default{
    // Including some components that will be used in the page.
    components: {
        LoadingSpinner,
        ErrorMsg,
    },

    data: function(){
        return{
            message:"",
            errormsg: "",
            loading: false,

            username: localStorage.getItem('Username'),
            BearerToken: localStorage.getItem('BearerToken'),
            
            user_searchList: [],
        }
    },

    methods:{
        async initialize() {
			this.$user_state.current_view = this.$views.SEARCH;
            console.log(this.username)
            console.log(this.BearerToken)
            this.$user_state.username = this.username
			this.$user_state.headers.Authorization = this.BearerToken
		},
        
		async goToProfile(usname, id) {
            this.$searched_state.username = usname;
            this.$searched_state.headers.Authorization = id
			this.$router.push({ path: `/users/${usname}/profile` });
		},
        async searchedUsers(){
            // Re-initializing variables to their default value.
            this.$user_state.current_view = this.$views.SEARCH;
            console.log(this.$user_state.username)
			this.errormsg = "";
			this.loading = true;

            /*

            this.$user_state.username = this.username
			this.$user_state.headers.Authorization = this.BearerToken

            */

            let search = document.querySelector("input").value;
            search = search.trim();

            if (search.length>0){
                const seacher_id = this.$user_state.headers.Authorization;
                
                if (seacher_id==null){
                    return
                }
            }

            try{
                let response = await this.$axios.get("/users/" + this.$user_state.username + "/", {
                    params: {
                        "searched-id": search
                    },

                    headers:{
                        "Authorization": this.$user_state.headers.Authorization,
                        "user-id": this.$user_state.username
                    }
                })

                if (response.status== 200){
                    // list of user
                    this.user_searchList = response.data;

                } else if (response.status==204){
                    this.message = "No users found"
                }else{
                    this.user_searchList = [];
                }


            } catch(e){
                // If an error is encountered, display it!
				if (e.response && e.response.status === 400) {
					this.errormsg = "The body was not parsable JSON or username invalid" + e.toString();
                } else if (e.response && e.response.status === 403) {
                    this.errormsg = "An Unauthorized Action has been blocked. You are not allowed to do this action because you are not the profile's owner." + e.toString();
                } else if (e.response && e.response.status === 404) {
                    this.errormsg = "The searcher-id is not found" + e.toString();
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later." + e.toString();
                } else {
                    this.errormsg = "Please Login before with an Authorized profile to view this page. " + e.toString();
                }
            }
            this.loading = false;
        },

        async mounted(){
            await this.initialize();
            await this.searchedUsers();
        }
    }
}
</script>

<template>
	<div class="container-fluid">
		<div class="row boxrel1">
			<div class="col d-flex justify-content-center mb-2">
				<h1>{{ this.$user_state.username }}'s Searching</h1>
			</div>
		</div>
            <div class="row mt-2 boxrel2">
			<div class="col d-flex justify-content-center">
				<div class="input-group mb-3 w-50">
					<input
                        class="form-control" 
                        id="SearchBox" 
                        type="text" 
                        placeholder="Search Username..." 
                        aria-label="Search"/>
					<div class="input-group-append">
						<button class="btn btn-primary btn-block mb-4" 
						@click="searchedUsers">
						Search</button>
					</div>
                    
				</div>
                
			</div>
            <div class="row">
			<UserIdent v-for="(user,index) in user_searchList" 
			:key="index" 
			:identifier="user['user-id']['identifier']"
			:nickname="user['username']"

            @clickedUser="goToProfile(user['username'], user['user-id']['identifier'])"
			/>
        </div>
            
        </div>

    <div v-if="this.message!='' " class="h-25">
			<h2 class="d-flex justify-content-center mt-5" style="color: black;">{{this.message }}</h2>
    </div>

    </div>
	<!-- Let's report the Error Message(if any), and the Loading Spinner if needed. -->
    <div class="row">
            <LoadingSpinner v-if="loading"></LoadingSpinner>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>

</template>

<style>


.boxrel1 {
 position:relative; 
 top:35px;
}

.boxrel2 {
 position:relative; 
 top:50px;
}

.boxrel3 {
 position:relative; 
 left:100px;
}

</style>