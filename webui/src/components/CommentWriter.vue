<script>
export default {	
	data(){
		return{
            errormsg: "",
            message: "",
			
			commentValue:"",
			commentid: ""
		}
	},
	props:['modal_id','comments_list','photo_owner','photo_id'],

	methods: {


		async restore(){
			if(document.value != ""){
				document.value = "";
			} 
		},

		async addComment(){
			try{
					
				const req_body = {
						"photo-id": {
							"code-image": {
							"identifier": this.photo_id,
							}
						},
						"text": this.commentValue,
						"user-id": {
							"code-user": {
							"identifier": this.$user_state.headers.Authorization
							}
						},
						"from-user": this.$user_state.username
					}

				let response = await this.$axios.post("/users/" + this.photo_owner + "/homepage/" + this.photo_id+"/comments", 
						req_body, {
							headers: {
								"Authorization": this.$user_state.headers.Authorization
							}
						});
				
				// this.comments_list.push(response.data)
				this.commentid = response.data["comment-id"]["code-comment"]["identifier"]
				console.log(this.commentid)

				this.$emit('addComment',response.data)

				console.log(this.commentValue)
				this.commentValue = ""
				
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

		eliminateCommentToParent(value){
			this.$emit('eliminateComment',value)
		},

		addCommentToParent(newCommentJSON){
			this.$emit('addComment',newCommentJSON)
		},
	},
}
</script>

<template>


    <div class="modal fade my-modal-disp-none" :id="modal_id" tabindex="-1" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-dialog modal-dialog-scrollable ">
            <div class="modal-content">

                <div class="modal-header">
                    <h1 class="modal-title fs-5" :id="modal_id">Comments</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>

                <div class="modal-body">
                    <PhotoComment v-for="(comm,index) in comments_list" 
					:key="index" 
					:author="comm['user-id']['code-user']['identifier']" 
					:nickname="comm['from-user']"
					:comment_id="comm['comment-id']['code-comment']['identifier']"
					:photo_id="comm['photo-id']['code-image']['identifier']"
					:content="comm['text']"
					:photo_owner="photo_owner"
					

					@eliminateComment="eliminateCommentToParent"
					/>



                </div>
                <div class="modal-footer d-flex justify-content-center w-100">
                    <div class="row w-100 ">
                        <div class="col-10">
                            <div class="mb-3 me-auto boxrel1 ">

                                <textarea class="form-control" id="exampleFormControlTextarea1" 
								placeholder="Add a comment..." rows="1" maxLength="30" v-model="commentValue"></textarea>
                            </div>
                        </div>

						<div class="col-2 d-flex align-items-center">
                            <button type="button" class="btn btn-primary" 
							@click.prevent="addComment">
							Send
							</button>
                    </div>
				</div>
            </div>
        </div>
    </div>
</div>

</template>

<style> 
.my-modal-disp-none{
	display: none;
}
.boxrel1 {
 position:relative; 
 top:7px;
}
</style>