import { User } from '../../../shared/data/models/User';

export interface ILoginRepository {
  registerUser(userData: Omit<User, 'id'>): Promise<string>;
  loginUser(userData: Omit<User, 'id'>): Promise<string>;
  refreshToken(token: string): Promise<string>;
}
