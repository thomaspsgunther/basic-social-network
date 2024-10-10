import User from '../../../shared/data/models/User';
import userApi from '../api/userApi';

class UserRepository {
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
    } catch (error) {
      console.error('User update failed:', error);
      return false;
    }
  }

  async deleteUser(id: string): Promise<boolean> {
    try {
      await userApi.remove(id);
      return true;
    } catch (error) {
      console.error('User delete failed:', error);
      return false;
    }
  }
}

export default new UserRepository();
