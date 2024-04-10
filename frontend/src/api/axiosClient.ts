import { auth } from '../config/firebase';
import axios from 'axios';
import { notification } from 'antd';
const axiosApiInstance = axios.create();

type ApiError = {
    message?: string;
    type?: string;
};

// Request interceptor for API calls
// Use the below interceptor to make transformations to an axios request,
// before the request is sent to the API
axiosApiInstance.interceptors.request.use(
    async (config) => {
        // set token
        const access_token = await auth.currentUser?.getIdToken();
        config.headers = {
            Accept: '*/*',
            'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8',
            Authorization: `Bearer ${access_token}`
        };

        // provide full url
        config.url = process.env.REACT_APP_API! + config.url;
        if (config.data) {
            config.data = JSON.stringify(config.data);
        } else {
            config.data = '{}';
        }

        return config;
    },
    (error) => {
        notification.error({ message: 'Unexpected request error', description: error.message });
        Promise.reject(error);
    }
);

//Use the below interceptor to handle axios (request) errors
axiosApiInstance.interceptors.response.use(
    (response) => {
        return response;
    },
    async function (error) {
        // Error 😨
        if (error.response) {
            /*
             * The request was made and the server responded with a
             * status code that falls out of the range of 2xx
             */
            let apiError = error.response.data as ApiError;
            if (apiError.type == "warning") {
                notification.warning({ message: 'Warning!', description: apiError.message });
            } else {
                notification.error({ message: 'Error ' + error.response.status, description: apiError.message });
            }
        } else if (error.request) {
            /*
             * The request was made but no response was received, `error.request`
             * is an instance of XMLHttpRequest in the browser and an instance
             * of http.ClientRequest in Node.js
             */
            //e.g. the backend was stopped
            // console.log(error.request.data); // doesnt work
            notification.error({ message: 'Connection error' }); // alt. Network Error
        } else {
            // Something happened in setting up the request and triggered an Error
            notification.error({ message: 'Unexpected response error', description: error.message });
        }
        return Promise.reject(error);
    }
);

export default axiosApiInstance;
