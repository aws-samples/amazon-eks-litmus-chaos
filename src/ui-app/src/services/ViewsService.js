import axios from "axios";

const COUNTER_API_BASE_URL = "http://localhost:9090/api/counter"; // local

class ViewsService {

    getViews(){
        // console.log(COUNTER_API_BASE_URL);
        return axios.get(COUNTER_API_BASE_URL + '/count');
    }

    addView(){
        return axios.post(COUNTER_API_BASE_URL + '/count');
    }
}

export default new ViewsService()