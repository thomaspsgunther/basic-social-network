import { axiosInstance } from '@/src/core/axios/axiosInstance';
import { Post } from '@/src/features/shared/data/models/Post';

export const postApi = {
    create: async (post: Omit<Post, 'id'>) => {
        const response = await axiosInstance.post(`/posts`, post);

        return response;
    },

      list: async (limit: number = 10, cursor: string = '') => {
        const response = await axiosInstance.get(`/posts`, {
            params: { limit, cursor }
        });
        return response;
    },

    like: async (userId: string, postId: string) => {
        const response = await axiosInstance.post(`/posts/likes/${userId}_${postId}`);
        return response;
    },

    unlike: async (userId: string, postId: string) => {
        const response = await axiosInstance.delete(`/posts/likes/${userId}_${postId}`);
        return response;
    },

    checkLiked: async (userId: string, postId: string) => {
        const response = await axiosInstance.get(`/posts/check/${userId}_${postId}`);
        return response;
    },

    getLikes: async (postId: string) => {
        const response = await axiosInstance.get(`/posts/likes/${postId}`);
        return response;
    },
    
};
