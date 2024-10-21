import { Post } from '@/src/features/shared/data/models/Post';
import { User } from '@/src/features/shared/data/models/User';

import { IPostRepository } from '../repositories/PostRepository';

interface IPostUsecase {
  createPost(post: Omit<Post, 'id'>): Promise<Post>;
  listPosts(limit: number, cursor?: string): Promise<Post[]>;
  getPost(id: string): Promise<Post>;
  updatePost(post: Post): Promise<boolean>;
  deletePost(id: string): Promise<boolean>;
  likePost(userId: string, postId: string): Promise<boolean>;
  unlikePost(userId: string, postId: string): Promise<boolean>;
  checkIfUserLikedPost(userId: string, postId: string): Promise<boolean>;
  getLikes(postId: string): Promise<User[]>;
}

export class PostUsecaseImpl implements IPostUsecase {
  private repository: IPostRepository;

  constructor(repository: IPostRepository) {
    this.repository = repository;
  }

  async createPost(post: Omit<Post, 'id'>): Promise<Post> {
    if (!post.image) {
      throw new Error('post image is required');
    }

    const createdPost: Post = await this.repository.createPost(post);

    return createdPost;
  }

  async listPosts(limit: number, cursor?: string): Promise<Post[]> {
    const posts: Post[] = await this.repository.listPosts(limit, cursor);

    return posts;
  }

  async getPost(id: string): Promise<Post> {
    const post: Post = await this.repository.getPost(id);

    return post;
  }

  async updatePost(post: Post): Promise<boolean> {
    if (!post.image) {
      throw new Error('post image is required');
    }

    await this.repository.updatePost(post);

    return true;
  }

  async deletePost(id: string): Promise<boolean> {
    await this.repository.deletePost(id);

    return true;
  }

  async likePost(userId: string, postId: string): Promise<boolean> {
    const didLike: boolean = await this.repository.likePost(userId, postId);

    return didLike;
  }

  async unlikePost(userId: string, postId: string): Promise<boolean> {
    const didUnlike: boolean = await this.repository.unlikePost(userId, postId);

    return didUnlike;
  }

  async checkIfUserLikedPost(userId: string, postId: string): Promise<boolean> {
    const liked: boolean = await this.repository.checkIfUserLikedPost(
      userId,
      postId,
    );

    return liked;
  }

  async getLikes(id: string): Promise<User[]> {
    const users: User[] = await this.repository.getLikes(id);

    return users;
  }
}
