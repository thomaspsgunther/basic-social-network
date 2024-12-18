import { Post } from '@/src/features/shared/data/models/Post';
import { User } from '@/src/features/shared/data/models/User';

export interface IUserRepository {
  getUserById(id: string): Promise<User>;
  getUsersBySearch(searchTerm: string): Promise<User[]>;
  listUserPosts(id: string, limit: number, cursor?: string): Promise<Post[]>;
  updateUser(user: User): Promise<boolean>;
  deleteUser(id: string): Promise<boolean>;
  followUser(followerId: string, followedId: string): Promise<boolean>;
  unfollowUser(followerId: string, followedId: string): Promise<boolean>;
  userFollowsUser(followerId: string, followedId: string): Promise<boolean>;
  getUserFollowers(id: string): Promise<User[]>;
  getUserFollowed(id: string): Promise<User[]>;
}
