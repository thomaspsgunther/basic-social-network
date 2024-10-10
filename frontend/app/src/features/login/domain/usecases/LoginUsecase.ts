import { ILoginRepository } from '../repositories/LoginRepository';
import { User } from '../../../shared/data/models/User';

interface ILoginUsecase {
  registerUser(userData: Omit<User, 'id'>): Promise<string>;
  loginUser(userData: Omit<User, 'id'>): Promise<string>;
  refreshToken(token: string): Promise<string>;
}

export class LoginUsecaseImpl implements ILoginUsecase {
  private repository: ILoginRepository;

  constructor(repository: ILoginRepository) {
    this.repository = repository;
  }

  async registerUser(userData: Omit<User, 'id'>): Promise<string> {
    this.validateUserData(userData);

    const token: string = await this.repository.registerUser(userData);

    return token;
  }

  async loginUser(userData: Omit<User, 'id'>): Promise<string> {
    this.validateUserData(userData);

    const token: string = await this.repository.loginUser(userData);

    return token;
  }

  async refreshToken(token: string): Promise<string> {
    const newToken: string = await this.repository.refreshToken(token);

    return newToken;
  }

  private validateUserData(userData: Omit<User, 'id'>) {
    if (!userData.email || !userData.password) {
      throw new Error('email and password are required');
    }
  }
}
