import { ILoginRepository } from '../../domain/repositories/LoginRepository';
import { User } from '../../../shared/data/models/User';
import { loginApi } from '../api/loginApi';

export class LoginRepositoryImpl implements ILoginRepository {
  async registerUser(userData: Omit<User, 'id'>): Promise<string> {
    const response = await loginApi.register(userData);
    const token: string = response.data;

    return token;
  }

  async loginUser(userData: Omit<User, 'id'>): Promise<string> {
    const response = await loginApi.login(userData);
    const token: string = response.data;

    return token;
  }

  async refreshToken(token: string): Promise<string> {
    const response = await loginApi.refreshToken(token);
    const newToken: string = response.data;

    return newToken;
  }
}
