import React, { createContext, useState, useEffect, ReactNode } from 'react';
import * as SecureStore from 'expo-secure-store';
import { User } from '../../features/shared/data/models/User';
import { setAuthToken } from '../axios/axiosInstance';
import { jwtDecode } from 'jwt-decode';
import { DecodedToken } from '../../features/shared/data/models/DecodedToken';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { RootStackParamList } from '../navigation/types';
import { LoginRepositoryImpl } from '../../features/login/data/repositories/LoginRepositoryImpl';
import { LoginUsecaseImpl } from '../../features/login/domain/usecases/LoginUsecase';
import { UserRepositoryImpl } from '../../features/users/data/repositories/UserRepositoryImpl';
import { UserUsecaseImpl } from '../../features/users/domain/usecases/UserUsecase';

interface AuthContextType {
  token: string | null;
  authUser: User | null;
  isAuthenticated: boolean | null;
  register: (userData: Omit<User, 'id'>) => Promise<void>;
  login: (userData: Omit<User, 'id'>) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [token, setToken] = useState<string | null>(null);
  const [authUser, setAuthUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const [refreshTimer, setRefreshTimer] = useState<NodeJS.Timeout | null>(null);
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>();
  const loginRepository = new LoginRepositoryImpl();
  const loginUsecase = new LoginUsecaseImpl(loginRepository);
  const userRepository = new UserRepositoryImpl();
  const userUsecase = new UserUsecaseImpl(userRepository);

  useEffect(() => {
    const loadStoredToken = async () => {
      try {
        const storedToken = await SecureStore.getItemAsync('token');
        if (storedToken) {
          const decodedToken = decodeToken(storedToken);
          if (decodedToken != null && decodedToken.exp) {
            const currentTime = Math.floor(Date.now() / 1000);
            if (decodedToken.exp < currentTime) {
              logout();

              return;
            }

            setToken(storedToken);
            setAuthToken(storedToken);
            setRefreshTimerLogic(storedToken);
            setIsAuthenticated(true);
            await fetchUser(storedToken);
          } else {
            logout();
          }
        } else {
          logout();
        }
      } catch (_error) {
        logout();
      }
    };

    loadStoredToken();

    return () => {
      if (refreshTimer) {
        clearTimeout(refreshTimer);
      }
    };
  }, []);

  const register = async (userData: Omit<User, 'id'>) => {
    const newToken = await loginUsecase.registerUser(userData);
    await SecureStore.setItemAsync('token', newToken);
    setToken(newToken);
    setAuthToken(newToken);
    setRefreshTimerLogic(newToken);
    setIsAuthenticated(true);
    await fetchUser(newToken);
  };

  const login = async (userData: Omit<User, 'id'>) => {
    const newToken = await loginUsecase.loginUser(userData);
    await SecureStore.setItemAsync('token', newToken);
    setToken(newToken);
    setAuthToken(newToken);
    setRefreshTimerLogic(newToken);
    setIsAuthenticated(true);
    await fetchUser(newToken);
  };

  const refreshToken = async (token: string) => {
    try {
      const newToken = await loginUsecase.refreshToken(token);
      await SecureStore.setItemAsync('token', newToken);
      setToken(newToken);
      setAuthToken(newToken);
      setIsAuthenticated(true);
    } catch (_error) {
      logout();
      navigation.navigate('Login');
    }
  };

  const logout = async () => {
    await SecureStore.deleteItemAsync('token');
    setToken(null);
    setAuthToken(null);
    setRefreshTimer(null);
    setIsAuthenticated(false);
    setAuthUser(null);
  };

  const decodeToken = (token: string): DecodedToken | null => {
    try {
      const decoded = jwtDecode<DecodedToken>(token);

      return decoded;
    } catch (_error) {
      return null;
    }
  };

  const fetchUser = async (token: string) => {
    try {
      const decodedToken = decodeToken(token);
      if (decodedToken != null) {
        const userList: User[] = await userUsecase.getUsersById(
          decodedToken.id,
        );

        if (!userList) {
          throw new Error('user not found');
        }

        setAuthUser(userList[0]);
      } else {
        throw new Error('invalid token');
      }
    } catch (_error) {
      throw new Error('failed to fetch user');
    }
  };

  const setRefreshTimerLogic = (token: string) => {
    const decodedToken = decodeToken(token);
    if (decodedToken && decodedToken.exp) {
      const currentTime = Math.floor(Date.now() / 1000);
      const refreshTime = (decodedToken.exp - currentTime - 300) * 1000; // 5 minutes before expiration
      if (refreshTimer) {
        clearTimeout(refreshTimer);
      }
      if (refreshTime > 0) {
        setRefreshTimer(setTimeout(() => refreshToken(token), refreshTime));
      }
    }
  };

  return (
    <AuthContext.Provider
      value={{ token, authUser, isAuthenticated, register, login, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export { AuthContext, AuthProvider };
