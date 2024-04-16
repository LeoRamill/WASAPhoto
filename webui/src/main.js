import {createApp,
     reactive} from 'vue'

// import router, app and axios
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';

// import all the components in main.js
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Photo from './components/PhotoPost.vue'
import Like from './components/Like.vue'
import UserIdent from './components/UserIdent.vue'
import Follower from './components/Follower.vue'
import CommentWriter from './components/CommentWriter.vue'
import PhotoComment from './components/PhotoComment.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)

// define costant views
const views = {
     LOGIN: "login",
     STREAM: "register",
     PROFILE: "profile",
     BAN: "ban",
     SEARCH: "search",
     UPDATE: "update"

 }
 

// define the state of login user
 var state = {
     headers: {
         Authorization: null
     },
     username: null,
     current_view: null
 
 }
 app.config.globalProperties.$user_state = reactive(state);
// define the state of searched user --> use in case of searching
 var state_searched = {
    headers: {
        Authorization: null
    },
    username: null,
    current_view: null

}
 app.config.globalProperties.$searched_state = reactive(state_searched);
 
 app.config.globalProperties.$views = views;
 app.config.globalProperties.$axios = axios;


 


app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("Photo", Photo);
app.component("Like", Like);
app.component("UserIdent", UserIdent);
app.component("Follower", Follower);
app.component("CommentWriter", CommentWriter)
app.component("PhotoComment", PhotoComment)

app.use(router)
app.mount('#app')
