import axios from 'axios';

import { axiosInstance } from '@/src/core/axios/axiosInstance';
import { User } from '@/src/features/shared/data/models/User';

export const loginApi = {
  register: async (userData: Omit<User, 'id'>) => {
    try {
      const response = await axiosInstance.post('/login/register', userData);

      return response;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response) {
          const backendMessage = error.response.data;
          if (backendMessage) {
            throw new Error(backendMessage);
          } else {
            throw new Error(error.message);
          }
        } else if (error.request) {
          throw new Error(error.request);
        } else {
          throw new Error(error.message);
        }
      } else {
        throw new Error('unexpected error');
      }
    }
  },
  login: async (userData: Omit<User, 'id'>) => {
    try {
      const response = await axiosInstance.post('/login', userData);

      return response;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response) {
          const backendMessage = error.response.data;
          if (backendMessage) {
            throw new Error(backendMessage);
          } else {
            throw new Error(error.message);
          }
        } else if (error.request) {
          throw new Error(error.request);
        } else {
          throw new Error(error.message);
        }
      } else {
        throw new Error('unexpected error');
      }
    }
  },
  refreshToken: async (token: string) => {
    const response = await axiosInstance.post('/login/refreshtoken', token);

    return response;
  },
};
