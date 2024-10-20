import { Post } from '@/src/features/shared/data/models/Post';
import { User } from '@/src/features/shared/data/models/User';

export interface IUserRepository {
  getUsersById(idList: string): Promise<User[]>;
  getUsersBySearch(searchTerm: string): Promise<User[]>;
  getUserPosts(id: string, limit: number, cursor?: string): Promise<Post[]>;
  updateUser(user: User): Promise<boolean>;
  deleteUser(id: string): Promise<boolean>;
  followUser(followerId: string, followedId: string): Promise<boolean>;
  unfollowUser(followerId: string, followedId: string): Promise<boolean>;
  userFollowsUser(followerId: string, followedId: string): Promise<boolean>;
  getUserFollowers(id: string): Promise<User[]>;
  getUserFollowed(id: string): Promise<User[]>;
}
