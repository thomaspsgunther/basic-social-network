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
    const response = await axiosInstance.delete(`/posts/${id}`);

    return response;
  },
  like: async (userId: string, postId: string) => {
    const response = await axiosInstance.post(
      `/posts/${postId}/likes/${userId}`,
    );

    return response;
  },
  getLikes: async (id: string) => {
    const response = await axiosInstance.get(`/posts/${id}/likes`);

    return response;
  },
  unlike: async (userId: string, postId: string) => {
    const response = await axiosInstance.delete(
      `/posts/${postId}/likes/${userId}`,
    );

    return response;
  },
  checkLiked: async (userId: string, postId: string) => {
    const response = await axiosInstance.get(
      `/posts/${postId}/likes/check/${userId}`,
    );

    return response;
  },
};
