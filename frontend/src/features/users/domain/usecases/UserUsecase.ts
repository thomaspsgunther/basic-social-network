import { User } from '@/src/features/shared/data/models/User';

import { IUserRepository } from '../repositories/UserRepository';

interface IUserUsecase {
  getUsersById(idList: string): Promise<User[]>;
  getUsersBySearch(searchTerm: string): Promise<User[]>;
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

  async getUsersById(idList: string): Promise<User[]> {
    const users: User[] = await this.repository.getUsersById(idList);

    return users;
  }

  async getUsersBySearch(searchTerm: string): Promise<User[]> {
    const users: User[] = await this.repository.getUsersBySearch(searchTerm);

    return users;
  }

  async updateUser(user: User): Promise<boolean> {
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
