import { Post } from '@/src/features/shared/data/models/Post';
import { User } from '@/src/features/shared/data/models/User';

export interface IPostRepository {
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
