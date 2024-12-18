import { Post } from '@/src/features/shared/data/models/Post';
import { User } from '@/src/features/shared/data/models/User';

import { IPostRepository } from '../../domain/repositories/PostRepository';
import { postApi } from '../api/postApi';

export class PostRepositoryImpl implements IPostRepository {
  async createPost(post: Omit<Post, 'id'>): Promise<Post> {
    const response = await postApi.create(post);
    const createdPost: Post = response.data;

    return createdPost;
  }

  async listPosts(limit: number, cursor?: string): Promise<Post[]> {
    const response = await postApi.list(limit, cursor);
    const posts: Post[] = response.data ? response.data : [];

    return posts;
  }

  async getPost(id: string): Promise<Post> {
    const response = await postApi.get(id);
    const post: Post = response.data;

    return post;
  }

  async updatePost(post: Post): Promise<boolean> {
    await postApi.update(post);

    return true;
  }

  async deletePost(id: string): Promise<boolean> {
    await postApi.remove(id);

    return true;
  }

  async likePost(userId: string, postId: string): Promise<boolean> {
    await postApi.like(userId, postId);

    return true;
  }

  async unlikePost(userId: string, postId: string): Promise<boolean> {
    await postApi.unlike(userId, postId);

    return true;
  }

  async checkIfUserLikedPost(userId: string, postId: string): Promise<boolean> {
    const response = await postApi.checkLiked(userId, postId);
    const liked: boolean = response.data.liked;

    return liked;
  }

  async getLikes(id: string): Promise<User[]> {
    const response = await postApi.getLikes(id);
    const users: User[] = response.data;

    return users;
  }
}
