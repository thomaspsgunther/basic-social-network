import { IUserRepository } from '../../domain/repositories/UserRepository';
import { User } from '../../../shared/data/models/User';
import { userApi } from '../api/userApi';

export class UserRepositoryImpl implements IUserRepository {
  async getUsersById(idList: string): Promise<User[]> {
    const response = await userApi.get(idList);
    const users: User[] = response.data;

    return users;
  }

  async getUsersBySearch(searchTerm: string): Promise<User[]> {
    const response = await userApi.search(searchTerm);
    const users: User[] = response.data;

    return users;
  }

  async updateUser(user: User): Promise<boolean> {
    try {
      await userApi.update(user);

      return true;
    } catch (_error) {
      return false;
    }
  }

  async deleteUser(id: string): Promise<boolean> {
    try {
      await userApi.remove(id);

      return true;
    } catch (_error) {
      return false;
    }
  }

  async followUser(followerId: string, followedId: string): Promise<boolean> {
    try {
      await userApi.follow(followerId, followedId);

      return true;
    } catch (_error) {
      return false;
    }
  }

  async unfollowUser(followerId: string, followedId: string): Promise<boolean> {
    try {
      await userApi.unfollow(followerId, followedId);

      return true;
    } catch (_error) {
      return false;
    }
  }

  async userFollowsUser(
    followerId: string,
    followedId: string,
  ): Promise<boolean> {
    const response = await userApi.userFollowsUser(followerId, followedId);
    const follows: boolean = response.data;

    return follows;
  }

  async getUserFollowers(id: string): Promise<User[]> {
    const response = await userApi.getFollowers(id);
    const users: User[] = response.data;

    return users;
  }

  async getUserFollowed(id: string): Promise<User[]> {
    const response = await userApi.getFollowed(id);
    const users: User[] = response.data;

    return users;
  }
}
