import User from '../models/User';
import { axiosInstance } from './axiosInstance';

const loginApi = {
  login: async (userData: Omit<User, 'id'>) => {
    const response = await axiosInstance.post('/login', userData);
    return response;
  },
  refreshToken: async (token: string) => {
    const response = await axiosInstance.post('/login/refreshtoken', token);
    return response;
  },
};

export default loginApi;
