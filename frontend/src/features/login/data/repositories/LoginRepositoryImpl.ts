import { User } from '@/src/features/shared/data/models/User';

import { ILoginRepository } from '../../domain/repositories/LoginRepository';
import { loginApi } from '../api/loginApi';

export class LoginRepositoryImpl implements ILoginRepository {
  async registerUser(userData: Omit<User, 'id'>): Promise<string> {
    const response = await loginApi.register(userData);
    const token: string = response.data.token;

    return token;
  }

  async loginUser(userData: Omit<User, 'id'>): Promise<string> {
    const response = await loginApi.login(userData);
    const token: string = response.data.token;

    return token;
  }

  async refreshToken(token: string): Promise<string> {
    const response = await loginApi.refreshToken(token);
    const newToken: string = response.data.token;

    return newToken;
  }
}
