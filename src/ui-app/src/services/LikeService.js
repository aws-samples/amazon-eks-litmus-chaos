import axios from "axios";

const LIKE_API_BASE_URL = "http://localhost:9080/api/like"; // local

class LikeService {

    getLikes(){
        // console.log(LIKE_API_BASE_URL);
        return axios.get(LIKE_API_BASE_URL);
    }
    
    addLike(id){
        return axios.post(LIKE_API_BASE_URL, id);
    }
}

export default new LikeService()