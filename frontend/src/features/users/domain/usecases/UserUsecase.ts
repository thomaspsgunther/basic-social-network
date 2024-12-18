import { Post } from '@/src/features/shared/data/models/Post';
import { User } from '@/src/features/shared/data/models/User';

import { IUserRepository } from '../repositories/UserRepository';

interface IUserUsecase {
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

export class UserUsecaseImpl implements IUserUsecase {
  private repository: IUserRepository;

  constructor(repository: IUserRepository) {
    this.repository = repository;
  }

  async getUserById(id: string): Promise<User> {
    const user: User = await this.repository.getUserById(id);

    return user;
  }

  async getUsersBySearch(searchTerm: string): Promise<User[]> {
    const users: User[] = await this.repository.getUsersBySearch(searchTerm);

    return users;
  }

  async listUserPosts(
    id: string,
    limit: number,
    cursor?: string,
  ): Promise<Post[]> {
    const posts: Post[] = await this.repository.listUserPosts(
      id,
      limit,
      cursor,
    );

    return posts;
  }

  async updateUser(user: User): Promise<boolean> {
    if (!user.username) {
      throw new Error('username is required');
    }

    const didUpdate: boolean = await this.repository.updateUser(user);

    return didUpdate;
  }

  async deleteUser(id: string): Promise<boolean> {
    const didDelete: boolean = await this.repository.deleteUser(id);

    return didDelete;
  }

  async followUser(followerId: string, followedId: string): Promise<boolean> {
    const didFollow: boolean = await this.repository.followUser(
      followerId,
      followedId,
    );

    return didFollow;
  }

  async unfollowUser(followerId: string, followedId: string): Promise<boolean> {
    const didUnfollow: boolean = await this.repository.unfollowUser(
      followerId,
      followedId,
    );

    return didUnfollow;
  }

  async userFollowsUser(
    followerId: string,
    followedId: string,
  ): Promise<boolean> {
    const follows: boolean = await this.repository.userFollowsUser(
      followerId,
      followedId,
    );

    return follows;
  }

  async getUserFollowers(id: string): Promise<User[]> {
    const users: User[] = await this.repository.getUserFollowers(id);

    return users;
  }

  async getUserFollowed(id: string): Promise<User[]> {
    const users: User[] = await this.repository.getUserFollowed(id);

    return users;
  }
}
