import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: 'http://10.0.2.2:8080/api/v1',
  timeout: 10000,
});

const setAuthToken = (token: string | null) => {
  if (token) {
    axiosInstance.defaults.headers['Authorization'] = `Bearer ${token}`;
  } else {
    delete axiosInstance.defaults.headers['Authorization'];
  }
};

export { axiosInstance, setAuthToken };
