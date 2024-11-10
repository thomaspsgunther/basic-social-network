import { axiosInstance } from '@/src/core/axios/axiosInstance';
import { User } from '@/src/features/shared/data/models/User';

export const userApi = {
  get: async (idList: string) => {
    const response = await axiosInstance.get(`/users/${idList}`);

    return response;
  },
  search: async (searchTerm: string) => {
    const response = await axiosInstance.get(`/users/search/${searchTerm}`);

    return response;
  },
  listPosts: async (id: string, limit: number, cursor?: string) => {
    if (cursor) {
      const response = await axiosInstance.get(
        `/users/${id}/posts?limit=${limit}&cursor=${cursor}`,
      );

      return response;
    } else {
      const response = await axiosInstance.get(
        `/users/${id}/posts?limit=${limit}`,
      );

      return response;
    }
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
    const response = await axiosInstance.get(`/users/${id}/followers`);

    return response;
  },
  getFollowed: async (id: string) => {
    const response = await axiosInstance.get(`/users/${id}/followed`);

    return response;
  },
};
