import { axiosInstance } from '@/src/core/axios/axiosInstance';
import { User } from '@/src/features/shared/data/models/User';

export const userApi = {
  get: async (idList: string) => {
    const response = await axiosInstance.get(`/users/${idList}`);

    return response;
  },
  search: async (searchTerm: string) => {
    const response = await axiosInstance.delete(`/users/search/${searchTerm}`);

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
  follow: async (followerId: string, followedId: string) => {
    const response = await axiosInstance.post(
      `/users/follow/${followerId}_${followedId}`,
    );

    return response;
  },
  unfollow: async (followerId: string, followedId: string) => {
    const response = await axiosInstance.delete(
      `/users/follow/${followerId}_${followedId}`,
    );

    return response;
  },
  userFollowsUser: async (followerId: string, followedId: string) => {
    const response = await axiosInstance.get(
      `/users/checkfollow/${followerId}_${followedId}`,
    );

    return response;
  },
  getFollowers: async (id: string) => {
    const response = await axiosInstance.get(`/users/followers/${id}`);

    return response;
  },
  getFollowed: async (id: string) => {
    const response = await axiosInstance.get(`/users/followed/${id}`);

    return response;
  },
};
