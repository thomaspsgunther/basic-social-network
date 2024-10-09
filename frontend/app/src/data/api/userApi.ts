import User from '../models/User';
import { axiosInstance } from './axiosInstance';

const userApi = {
  create: async (userData: Omit<User, 'id'>) => {
    const response = await axiosInstance.post('/users', userData);
    return response;
  },
  get: async (idList: string) => {
    const response = await axiosInstance.get(`/users/${idList}`);
    return response;
  },
  update: async (user: User) => {
    const response = await axiosInstance.put(`/users/${user.id}`, user);
    return response;
  },
  remove: async (id: string) => {
    const response = await axiosInstance.delete(`/users/${id}`);
    return response;
  },
  search: async (searchTerm: string) => {
    const response = await axiosInstance.delete(`/users/search/${searchTerm}`);
    return response;
  },
};

export default userApi;
