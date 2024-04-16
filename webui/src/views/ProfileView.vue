<script>
import ErrorMsg from '../components/ErrorMsg.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'

export default {

	data: function() {

		return {

			errormsg: "",
			loading: false,
            message:"",

			userExists: false,
			banStatus: false,

            nickname: "",
            // username: "",


			followStatus: false,
			currentIsBanned: false,

			followerCnt: 0,
			followingCnt:0,
			postCnt:0,

            uplphotoid: "",
			photos: [],
            following: [],
            followers: [],
            bans:[],

            BearerToken: localStorage.getItem('BearerToken'),
			usernameLogged: localStorage.getItem('Username'),
			username: localStorage.getItem('Username') == localStorage.getItem('usernameProfileToView') ? localStorage.getItem('Username') : localStorage.getItem('usernameProfileToView'),
		}
	},

    watch:{
        currentPath(newid,oldid){
            if (newid !== oldid){
                this.loadInfo()
            }
        },
    },

	computed:{

        currentPath(){
            return this.$route.params.username
        },
        

		sameUser(){
        
			return this.$route.params.username == this.$user_state.username
		},
	},

	methods: {

		async initialize() {
			this.$user_state.current_view = this.$views.PROFILE;
			this.$user_state.username = this.username;
			this.$user_state.headers.Authorization = this.BearerToken
		},

        async uploadFile(){
            let fileInput = document.getElementById('fileUploader')

            const file = fileInput.files[0];
            const reader = new FileReader();


            reader.readAsArrayBuffer(file);

            reader.onload = async () => {
                // Post photo: /users/:id/photos
                let response = await this.$axios.post("/users/"+this.$route.params.username+"/profile/photos", reader.result, {
                    headers: {
                    "Content-Type": file.type,
                    "Authorization": this.$user_state.headers.Authorization
                    },
                })
                this.photos.push(response.data)
                this.uplphotoid = response.data['photo-id']['identifier']
                this.postCnt += 1
                console.log(this.uplphotoid)

            };

        },


		async followClick(){
			this.username = this.$route.params.username;
            try{
                if (this.followStatus){ 
					const res = await this.$axios.delete("/users/" + this.$user_state.username + "/profile/following/" + this.username, {
						headers: {
							"Authorization": this.$user_state.headers.Authorization
						}
					});
					if (res.statusText != "No Content") {

						alert("Error: " + res.statusText);
						return
					}
                    this.followerCnt -=1
                }else{
					const req_body = {
						"code-user": {
							"identifier": this.$searched_state.headers.Authorization
						}
					}

					const res = await this.$axios.put("/users/" + this.$user_state.username + "/profile/following/" + this.username, 
					req_body, {
						headers: {
							"Authorization": this.$user_state.headers.Authorization
						}
					});

					if (res.statusText != "No Content") {

						alert("Error: " + res.statusText);
						return
					}
					this.followerCnt +=1
                }
                this.followStatus = !this.followStatus
            }catch (e){
                this.errormsg = e.toString();
            }
            
		},

		async banClick(){
			this.username = this.$route.params.username;
            try{
                if (this.banStatus){
					const res = await this.$axios.delete("/users/" + this.$user_state.username + "/profile/banned/" + this.username, {
						headers: {
							"Authorization": this.$user_state.headers.Authorization
						}
					});

					if (res.statusText != "No Content") {

						alert("Error: " + res.statusText);
						return
					}
                    this.loadInfo()
                }else{
					const req_body = {
						"code-user": {
							"identifier": this.$searched_state.headers.Authorization
						}
					}
					const res = await this.$axios.put("/users/" + this.$user_state.username + "/profile/banned/" + this.username, req_body,
					 {
						headers: {
							"Authorization": this.$user_state.headers.Authorization
						}
					});

					if (res.statusText != "No Content") {

						alert("Error: " + res.statusText);
						return
					}
                    this.followStatus = false
                }
                this.banStatus = !this.banStatus
            }catch(e){
                this.errormsg = e.toString();
            }
		},

        async getBans(){
            try {
				let response = await this.$axios.get("/users/" + this.$user_state.username + "/profile/banned", {
					headers: {
						"Content-Type": "application/json",
						"Authorization": this.$user_state.headers.Authorization
					}
				}
				)
                
                if(response.status==204){
                    this.bans = []
                }

				if (response.data != null){
					this.bans = response.data
				}


				
			} catch (e) {
				this.errormsg = e.toString()
			}
        },

		async loadInfo(){
            /*
                this.$user_state.username = this.username
                this.$user_state.headers.Authorization = this.BearerToken   
            */
                this.errormsg = "";
                this.loading = true;
                if (this.$route.params.username === undefined){
                    return
                }
                try{
                    console.log(this.$searched_state.username) 
                    console.log(this.$searched_state.headers.Authorization)
                    // Get user profile: /users/:id
                    let response = await this.$axios.get("/users/"+this.$route.params.username+"/profile",{
                    headers : {
						"Content-Type": "application/json"
					}});
                    console.log(response.data)
                    

                    this.banStatus = false
                    this.userExists = true
                    this.currentIsBanned = false


                    if (response.status === 206){
                        this.banStatus = true
                        return
                    }
                    
                    if (response.status === 204){
                        this.userExists = false
                    }

                    
                    // await this.getBans()
                    this.banStatus = this.bans.some(obj => obj['username'] === this.$route.params.username)
                    console.log(response.status)

                    if (this.banStatus === true){
                        this.bans.pop()
                    }
                    
                    this.nickname = response.data.nickname
                    this.followerCnt = response.data.followers != null ? response.data.followers.length : 0
                    this.followingCnt = response.data.following != null? response.data.following.length : 0
                    this.postCnt = response.data.posts != null ? response.data.posts.length : 0
                    this.followStatus = response.data.followers != null ? response.data.followers.find(obj => obj["user-id"]["identifier"] === this.$user_state.headers.Authorization) : false
                    this.photos = response.data.posts != null ? response.data.posts : []
                    this.followers = response.data.followers != null ? response.data.followers : []
                    this.following = response.data.following != null ? response.data.following : []

                    console.log(response.data)
                    console.log(this.userExists)
                    console.log(this.currentIsBanned)

                }catch(e){
                    this.currentIsBanned = true
                    if (e.response && e.response.status === 400) {
                        this.errormsg = "Request error, please Login before doing some action " + e.toString();
                        this.message = this.errormsg
                    } else if (e.response && e.response.status === 401) {
                        this.errormsg = "You are not allowed to do this action." + e.toString();
                        this.message = this.errormsg
                    } else if (e.response && e.response.status === 500) {
                        this.errormsg = "An internal error occurred. We will be notified. Please try again later." + e.toString();
                        this.message = this.errormsg;
                    }else if (e.response && e.response.status === 404) {
                        this.errormsg = "Doesn't Exist the Profile " + e.toString();
                        this.message = this.errormsg;}
                }
                this.loading = false;
            
		},

        goToChangeName(){
            this.$router.push('/users/'+this.$route.params.username+'/settings')
        },

        removePhotoFromList(photo_id){
			this.photos = this.photos.filter(obj => obj['photo-id']['identifier'] !== photo_id)
		},
	},

	async mounted(){
        this.initialize()
        this.getBans()
		this.loadInfo()

	}

}
</script>

<template>
    <Follower
		:modal_id="'foll_modal'+nickname" 
		:followers="followers" />

    <div class="container-fluid " v-if="!currentIsBanned && userExists">
        <div class="row boxrel1 ">
            <div class="col-15 d-flex justify-content-center">
                <div class="card w-50  container-fluid">

                    <div class="row">
                        <div class="col">
                            <div class="card-body d-flex justify-content-between align-items-center">
                                <h5 class="card-title p-0 me-auto mt-auto">@{{nickname }}</h5>

                                <button v-if="!sameUser && !banStatus" @click="followClick" class="btn btn-success ms-2">
                                    {{followStatus ? "Unfollow" : "Follow"}}
                                </button>

                                <button v-if="!sameUser" @click="banClick" class="btn btn-danger ms-2">
                                    {{banStatus ? "Unban" : "Ban"}}
                                </button>

                                <button v-if="sameUser" @click="goToChangeName" class="btn btn-primary ms-2">
                                    Change Name 
                                </button>
                                

                            </div>
                        </div>
                    </div>

                    <div v-if="!banStatus" class="row mt-1 mb-1">
                        <div class="col-4 d-flex justify-content-start">
                            <h6 class="ms-3 p-0 ">Posts: {{postCnt}}</h6>
                        </div>
                    
                        <div class="col-4 d-flex justify-content-end">
                            <h6 class=" p-0 me-3">Follower: {{followerCnt}}</h6>
                        </div>
                    
                        <div class="col-4 d-flex justify-content-end">
                            <h6 class=" p-0 me-3">Following: {{followingCnt}}</h6>
                        </div>
                    </div>
                </div>
            </div>
        </div>


        <div class="row boxrel2">

            <div class="container-fluid mt-3">

                <div class="row  ">
                    <div class="col-12 d-flex justify-content-center">
                        <h2> Personal Gallery</h2>
                        <input autocomplete="fileUploader" id="fileUploader" name="fileUploader" type="file" class="profile-file-upload" @change="uploadFile" accept=".jpg, .png">
                        <label v-if="sameUser" class="btn my-btn-add-photo ms-2 d-flex align-items-center boxrel3" for="fileUploader"> Add Image </label>
                    </div>
                </div>

                <div class="row ">
                    <div class="col-3"></div>
                    <div class="col-6">
                        <hr class="border border-dark">
                    </div>
                    <div class="col-3"></div>
                </div>
            </div>
        </div>

        <div>
            <LoadingSpinner v-if="loading"></LoadingSpinner>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>

        <div class="row">
            <div class="col boxrel1">

                <div v-if="!banStatus && postCnt>0">
                    <Photo v-for="(photo,index) in photos" 
                    :key="index" 
                    :owner="photo['nickname']"
                    :owner_id="photo['user-id']['identifier']"
                    :photo_id="photo['photo-id']['identifier']" 
                    :comments="photo['comment-collection']" 
                    :likes="photo['like-collection'] != null ? photo['like-collection'] : []" 
                    :upload_date="photo['date-time']" 
                    :isOwner="sameUser" 
                    
                    @removePhoto="removePhotoFromList"
                    />

                </div>
                
                <div v-else class="mt-5 ">
                    <h2 class="d-flex justify-content-center" style="color: black;">No posts yet</h2>
                </div>

            </div>
        </div>


    </div>

    
</template>

<style>
.profile-file-upload{
    display: none;
}

.my-nav-icon-gear{
    color: grey;
}
.my-nav-icon-gear:hover{
    transform: scale(1.3);
}

.my-btn-add-photo{
    background-color: 	limegreen	;
    border-color: grey;
}
.my-btn-add-photo:hover{
    color: white;
    background-color: 	limegreen	;
    border-color: white
}

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
 left:30px;
}

</style>

