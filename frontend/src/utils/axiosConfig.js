import axios from 'axios';

// Configure axios defaults
axios.defaults.baseURL = 'http://localhost:8080';

// REQUEST INTERCEPTOR REMOVED
// axios.interceptors.request.use(
//   config => {
//     if (config.url && config.url.startsWith('/v1/')) {
//       config.url = '/api' + config.url;
//     }
//     return config;
//   },
//   error => {
//     return Promise.reject(error);
//   }
// );

export default axios; 