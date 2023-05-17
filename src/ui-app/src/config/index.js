// set to true or false
const local = true;

var API_BASE_URL = '';
var LIKE_API_PORT = '';
var COUNTER_API_PORT = '';

if (local) {
    API_BASE_URL = 'http://localhost';
    LIKE_API_PORT = ':9080';
    COUNTER_API_PORT = ':9090';
} else {
    API_BASE_URL = 'REPLACE_ME'; // configure
}

export const LIKE_API_BASE_URL = `${API_BASE_URL}${LIKE_API_PORT}/api/like`;
export const COUNTER_API_BASE_URL = `${API_BASE_URL}${COUNTER_API_PORT}/api/counter`;