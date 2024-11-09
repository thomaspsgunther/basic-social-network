import { axiosInstance } from '@/src/core/axios/axiosInstance';

import { Comment } from '../models/Comment';

export const commentApi = {
  create: async (comment: Omit<Comment, 'id'>) => {
    const response = await axiosInstance.post('/comments', comment);

    return response;
  },
  getFromPost: async (postId: string) => {
    const response = await axiosInstance.get(`/comments/post/${postId}`);

    return response;
  },
  update: async (comment: Comment) => {
    const response = await axiosInstance.put(
      `/comments/${comment.id}`,
      comment,
    );

    return response;
  },
  remove: async (id: string) => {
    const response = await axiosInstance.delete(`/comments/${id}`);

    return response;
  },
};
