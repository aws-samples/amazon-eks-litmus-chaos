import axios from "axios";
import { COUNTER_API_BASE_URL } from "../config";

class ViewsService {

    getViews(){
        return axios.get(COUNTER_API_BASE_URL + '/count');
    }

    addView(){
        return axios.post(COUNTER_API_BASE_URL + '/count');
    }
}

export default new ViewsService()