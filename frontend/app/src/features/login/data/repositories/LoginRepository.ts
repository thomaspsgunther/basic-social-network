import loginApi from '../api/loginApi';
import User from '../../../shared/data/models/User';

class LoginRepository {
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

export default new LoginRepository();
