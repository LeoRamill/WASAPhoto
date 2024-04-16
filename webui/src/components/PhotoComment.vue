<script>
export default {
    data(){
        return {
            user: "",
        }
    },
	props: ['content','author','photo_owner','comment_id','photo_id','nickname'],

    methods:{
        async deleteComment(){
            try{
                // Delete comment: "/users/:id/photos/:photo_id/comments/:comment_id"
                await this.$axios.delete("/users/"+this.photo_owner+"/homepage/"+this.photo_id+"/comments/"+this.comment_id, {
							headers: {
								"Authorization": this.$user_state.headers.Authorization
							}
						})

                this.$emit('eliminateComment',this.comment_id)

            }catch (e){
				if (e.response && e.response.status === 400) {
					this.errormsg = "Request error, please Login before doing some action or ask to add a comment to a valid photo/comment/user." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 401) {
                    this.errormsg = "You are not allowed to do this action because you are not the stream's owner." + e.toString();
					this.message = this.errormsg
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later." + e.toString();
					this.message = this.errormsg;
                }else if (e.response && e.response.status === 403) {
                    this.errormsg = "No put comment because you are banned " + e.toString();
					this.message = this.errormsg;
                }else if (e.response && e.response.status === 404) {
                    this.errormsg = "No comment found, so you can't uncomment the post " + e.toString();
					this.message = this.errormsg;}
            }
        },
    },

    mounted(){
        this.user = this.$user_state.username
        console.log(this.user === this.author)
        console.log(this.user === this.photo_owner)
    }

}
</script>

<template>
	<div class="container-fluid">

        <hr>
        <div class="row">
            <div class="col-10">
                <h5> @{{nickname}}</h5>
            </div>

            <div class="col-2">
                <button v-if="user === nickname || user === photo_owner" class="my-trnsp-btn my-dlt-btn me-2" @click="deleteComment">
                    <div class="nav-link">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#trash-2" />
                            </svg>
                        </div>
                </button>
            </div>

        </div>

        <div class="row">
            <div class="col-12">
                {{content}}
            </div>

        </div>
        <hr>
    </div>
</template>

<style>
.my-btn-comm{
    border: none;
}
.my-btn-comm:hover{
    border: none;
    color: var(--color-red-danger);
    transform: scale(1.1);
}

</style>