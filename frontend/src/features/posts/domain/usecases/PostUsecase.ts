import { Post } from '@/src/features/shared/data/models/Post';
import { IPostRepository } from '../repositories/PostRepository';

interface IPostUsecase {
    createPost(post: Omit<Post, 'id'>): Promise<Post>;
    listPosts(limit: number, cursor?: string): Promise<Post[]>;
    likePost(userId: string, postId: string): Promise<boolean>;
    unlikePost(userId: string, postId: string): Promise<boolean>;
    checkIfUserLikedPost(userId: string, postId: string): Promise<boolean>;
    getLikes(postId: string): Promise<string[]>;
  }

export class PostUsecaseImpl implements IPostUsecase {
  private repository: IPostRepository;

  constructor(repository: IPostRepository) {
    this.repository = repository;
  }

  async createPost(post: Omit<Post, 'id'>): Promise<Post> {
    const createdPost: Post = await this.repository.createPost(post);
    return createdPost;
  }

  async listPosts(limit: number, cursor?: string): Promise<Post[]> {
    const posts: Post[] = await this.repository.listPosts(limit, cursor);
    return posts;
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
    const liked: boolean = await this.repository.checkIfUserLikedPost(userId, postId);
    return liked;
  }

  async getLikes(postId: string): Promise<string[]> {
    const userIds: string[] = await this.repository.getLikes(postId);
    return userIds;
  }
}