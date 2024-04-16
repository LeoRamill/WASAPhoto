<script>
export default {	

	props:['modal_id','likes'],

	methods: {
		async goToProfileLiker(identifier, nickname){
			this.$searched_state.username = nickname;
			this.$searched_state.headers.Authorization = identifier;
			this.$router.push({ path: `/users/${nickname}/profile` })
		},
	}


}
</script>

<template>
	<div class="modal fade modal-dialog modal-dialog-centered modal-dialog-scrollable my-modal-disp-none" :id="modal_id" tabindex="-1" aria-hidden="true">
	<div class="modal-dialog">
		<div class="modal-content">
		<div class="modal-header">
			<h1 class="modal-title fs-5" :id="modal_id">Likes</h1>
			<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
		</div>
		<div class="modal-body">
			<UserIdent v-for="(like,index) in likes" 
			:key="index" 
			:identifier="like['like-id']['code-like']['identifier']"
			:nickname="like['from-user']"

			/>
		</div>
		</div>
	</div>
	</div>
</template>

<style>
.my-modal-disp-none{
	display: none;
}

</style>