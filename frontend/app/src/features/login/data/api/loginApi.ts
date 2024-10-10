import { axiosInstance } from '../../../shared/data/api/axiosInstance';
import User from '../../../shared/data/models/User';

const loginApi = {
  register: async (userData: Omit<User, 'id'>) => {
    const response = await axiosInstance.post('/login/register', userData);
    return response;
  },
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
