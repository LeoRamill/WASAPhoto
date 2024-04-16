<script>
import CommentWriter from './CommentWriter.vue'

export default {
    data() {
        return {
            errormsg: "",
            message: "",

            photoURL: "",
            liked: false,
            allComments: [],
            allLikes: [],
            cntlike: 0,
            have_i_liked_this: false,
			comment_open: false
        };
    },
    props: ['owner','owner_id', 'likes', 'comments', "upload_date", "photo_id", "isOwner"],
    methods: {
        loadPhoto() {
            // Get photo : "/users/:id/photos/:photo_id"
            this.photoURL = __API_URL__ + "/users/" + this.owner + "/profile/photos/" + this.photo_id;
        },
        async deletePhoto() {
            try {
                // Delete photo: /users/:id/photos/:photo_id
                await this.$axios.delete("/users/" + this.owner + "/profile/photos/" + this.photo_id, {
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": this.$user_state.headers.Authorization
                    }
                });
                // location.reload()
                this.$emit("removePhoto", this.photo_id);
            }
            catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Request error, please Login before doing some action." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 401) {
                    this.errormsg = "You are not allowed to do this action because you are not the stream's owner." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later." + e.toString();
					this.message = this.errormsg;
                }else if (e.response && e.response.status === 404) {
                    this.errormsg = "Photo-id or usname-id not found " + e.toString();
					this.message = this.errormsg;
                }
            }
        },
        photoOwnerClick: function () {
            this.$searched_state.username = this.owner
            this.$searched_state.headers.Authorization = this.owner_id
            this.$router.replace("/users/" + this.owner + "/profile");

        },
        async Like() {
            try{
                const req_body = {
                    "like-id": {
                        "code-like": {
                            "identifier": this.$user_state.headers.Authorization
                        }
                    },
                    "photo-id": {
                        "code-image": {
                            "identifier": this.photo_id
                        }
                    },
                    "from-user": this.$user_state.username
                };
                let response = await this.$axios.put("/users/" + this.$user_state.username + "/homepage/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization, req_body, {
                    headers: {
                        "Authorization": this.$user_state.headers.Authorization
                    }
                });
                // Only for consistency, the component does this internally.
                this.liked = true;
                this.cntlike++;
                this.allLikes.push(response.data);
            }catch(e){
				if (e.response && e.response.status === 400) {
					this.errormsg = "Request error, please Login before doing some action or ask to add a like to a valid photo/like/user." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 401) {
                    this.errormsg = "You are not allowed to do this action because you are not the stream's owner." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later." + e.toString();
					this.message = this.errormsg;
                }else if (e.response && e.response.status === 403) {
                    this.errormsg = "No put like because you are banned " + e.toString();
					this.message = this.errormsg;
                }
            }

        },
        async Unlike() {
            try{
                await this.$axios.delete("/users/" + this.$user_state.username + "/homepage/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization, {
                headers: {
                    "Authorization": this.$user_state.headers.Authorization
                    }
                });
                this.liked = false;
                this.cntlike--;
                this.allLikes.pop();
            }catch(e){
				if (e.response && e.response.status === 400) {
					this.errormsg = "Request error, please Login before doing some action or ask to add a like to a valid photo/like/user." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 401) {
                    this.errormsg = "You are not allowed to do this action because you are not the stream's owner." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later." + e.toString();
					this.message = this.errormsg;
                }else if (e.response && e.response.status === 403) {
                    this.errormsg = "No put like because you are banned " + e.toString();
					this.message = this.errormsg;
                }else if (e.response && e.response.status === 404) {
                    this.errormsg = "No like found, so you can't unlike the post " + e.toString();
					this.message = this.errormsg;}
            }
        },

        async toggleLike() {
            this.liked = !this.liked;
            if (!this.liked) {
                this.Unlike();
            }
            else {
                this.Like();
            }
        },
        removeCommentFromList(value) {
            this.allComments = this.allComments.filter(obj => obj["comment-id"]["code-comment"]["identifier"] !== value);
        },
        addCommentToList(comment) {
            this.allComments.push(comment);
        },
    },
	async toggleCommentList() {
            this.comment_open = !this.comment_open;

        },
    async mounted() {
        await this.loadPhoto();
        if (this.likes != null) {
            this.allLikes = this.likes;
            this.cntlike = this.allLikes.length;
        }
        if (this.likes != null) {
            this.liked = this.allLikes.some(obj => obj['like-id']['code-like']['identifier'] === this.$user_state.headers.Authorization);
        }
        if (this.comments != null) {
            this.allComments = this.comments;
        }
        console.log(this.allComments);
    },
    components: { CommentWriter }
}
</script>

<template>
	<div class="container-fluid mt-3 mb-5 ">

        <Like 
		:modal_id="'like_modal'+photo_id" 
		:likes="allLikes" />

        <CommentWriter 
		:modal_id="'comment_modal'+photo_id" 
		:comments_list="allComments" 
		:photo_owner="owner" 
		:photo_id="photo_id"

		@eliminateComment="removeCommentFromList"
		@addComment="addCommentToList"

		/>

        <div class="d-flex flex-row justify-content-center">

            <div class="card my-card">
                <div class="d-flex justify-content-end">

					<button class="my-trnsp-btn m-0 p-1 me-auto" @click="photoOwnerClick">
                            	<i> @{{owner}}</i>
							</button>

                    <button v-if="isOwner" class="my-trnsp-btn my-dlt-btn me-2" @click="deletePhoto">
                        <div class="nav-link">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#trash-2" />
                            </svg>
                        </div>
					</button>

                </div>
                <div class="d-flex justify-content-center photo-background-color">
                    <img :src="photoURL" class="card-img-top img-fluid">
                </div>

                <div class="card-body">

                    <div class="container">

                        <div class="d-flex flex-row justify-content-end align-items-center mb-2">

							<div class="col-sm-2 my-trnsp-btn m-0 p-1 me-auto">
                    			<div @click="toggleLike()" class="icon-container">
                        			<svg width="24" height="24" viewBox="0 0 24 24">

                            			<path v-if="liked"
                                			d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                                			fill="red"></path>
                            			<path v-else
                                			d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                                			fill="none" stroke="red"></path>
                        			</svg>
									<div>
										<i data-bs-toggle="modal" :data-bs-target="'#like_modal'+photo_id" class="my-comment-color ">
									<p> Like: {{this.cntlike}}</p>
                                		</i>
									</div>
                        			
                    			</div>

							
							</div>
                            <button class=" my-trnsp-btn m-0 p-1  d-flex justify-content-center align-items-center" 
							data-bs-toggle="modal" :data-bs-target="'#comment_modal'+photo_id">

                                <i class="my-comment-color fa-regular fa-comment me-1" @click="commentClick"></i>
                                <i class="my-comment-color-2"> Comments: {{allComments != null ? allComments.length : 0}}</i>

                            </button>
						
                        </div>

                        <div class="d-flex flex-row justify-content-start align-items-center ">
                            <p> Uploaded on {{upload_date}}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
.photo-background-color{
	background-color: grey;
}

.my-card{
	width: 27rem;
	border-color: black;
	border-width: thin;
}

.my-heart-color{
	color: grey;
}
.my-heart-color:hover{
	color: red;
}

.my-comment-color {
	color: grey;
}
.my-comment-color:hover{
	color: black;
}

.my-comment-color-2{
	color:grey
}

.my-dlt-btn{
	font-size: 19px;
}
.my-dlt-btn:hover{
	font-size: 19px;
	color: var(--color-red-danger);
}
</style>