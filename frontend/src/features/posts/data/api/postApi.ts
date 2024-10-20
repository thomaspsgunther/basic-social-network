import { axiosInstance } from '@/src/core/axios/axiosInstance';
import { Post } from '@/src/features/shared/data/models/Post';

export const postApi = {
  create: async (post: Omit<Post, 'id'>) => {
    const response = await axiosInstance.post('/posts', post);

    return response;
  },
  list: async (limit: number, cursor?: string) => {
    if (cursor) {
      const response = await axiosInstance.get(
        `/posts?limit=${limit}&cursor=${cursor}`,
      );

      return response;
    } else {
      const response = await axiosInstance.get(`/posts?limit=${limit}`);

      return response;
    }
  },
  get: async (id: string) => {
    const response = await axiosInstance.get(`/posts/${id}`);

    return response;
  },
  update: async (post: Post) => {
    const response = await axiosInstance.put(`/posts/${post.id}`, post);

    return response;
  },
  remove: async (id: string) => {
    const response = await axiosInstance.delete(`/posts${id}`);

    return response;
  },
  like: async (userId: string, postId: string) => {
    const response = await axiosInstance.post(
      `/posts/likes/${userId}_${postId}`,
    );

    return response;
  },
  unlike: async (userId: string, postId: string) => {
    const response = await axiosInstance.delete(
      `/posts/likes/${userId}_${postId}`,
    );

    return response;
  },
  checkLiked: async (userId: string, postId: string) => {
    const response = await axiosInstance.get(
      `/posts/check/${userId}_${postId}`,
    );

    return response;
  },
  getLikes: async (postId: string) => {
    const response = await axiosInstance.get(`/posts/likes/${postId}`);

    return response;
  },
};
