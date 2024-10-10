import { User } from '../../../shared/data/models/User';

export interface IUserRepository {
  getUsersById(idList: string): Promise<User[]>;
  getUsersBySearch(searchTerm: string): Promise<User[]>;
  updateUser(user: User): Promise<boolean>;
  deleteUser(id: string): Promise<boolean>;
  //   follow(followerId: string, followedId: string): Promise<boolean>;
  //   unfollow(followerId: string, followedId: string): Promise<boolean>;
  //   userFollowsUser(followerId: string, followedId: string): Promise<boolean>;
  //   getFollowers(id: string): Promise<User[]>;
  //   getFollowed(id: string): Promise<User[]>;
}
