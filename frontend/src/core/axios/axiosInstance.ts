import axios from 'axios';

const apiHost = process.env.EXPO_PUBLIC_API_HOST;
const apiPort = process.env.EXPO_PUBLIC_API_PORT;

const axiosInstance = axios.create({
  baseURL: `http://${apiHost}:${apiPort}/api/v1`,
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
