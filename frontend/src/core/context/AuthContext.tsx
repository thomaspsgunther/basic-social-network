import {
  CommonActions,
  NavigationContainerRef,
} from '@react-navigation/native';
import * as SecureStore from 'expo-secure-store';
import { jwtDecode } from 'jwt-decode';
import React, { createContext, ReactNode, useEffect, useState } from 'react';
import { Alert, Platform } from 'react-native';

import { DecodedToken } from '@/src/features/login/data/models/DecodedToken';
import { LoginRepositoryImpl } from '@/src/features/login/data/repositories/LoginRepositoryImpl';
import { LoginUsecaseImpl } from '@/src/features/login/domain/usecases/LoginUsecase';
import { User } from '@/src/features/shared/data/models/User';
import { UserRepositoryImpl } from '@/src/features/users/data/repositories/UserRepositoryImpl';
import { UserUsecaseImpl } from '@/src/features/users/domain/usecases/UserUsecase';

import { setAuthToken } from '../axios/axiosInstance';
import { RootStackParamList } from '../navigation/types';

interface AuthContextType {
  token: string | null;
  authUser: User | null;
  isAuthenticated: boolean | null;
  register: (userData: Omit<User, 'id'>) => Promise<void>;
  login: (userData: Omit<User, 'id'>) => Promise<void>;
  logout: () => Promise<void>;
  logoutAndLeave: () => Promise<void>;
  setAuthUser: (user: User) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
  navigationRef: React.RefObject<NavigationContainerRef<RootStackParamList>>;
}

const AuthProvider: React.FC<AuthProviderProps> = ({
  children,
  navigationRef,
}) => {
  const [token, setToken] = useState<string | null>(null);
  const [authUser, setAuthUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const [refreshTimer, setRefreshTimer] = useState<NodeJS.Timeout | null>(null);
  const loginRepository = new LoginRepositoryImpl();
  const loginUsecase = new LoginUsecaseImpl(loginRepository);
  const userRepository = new UserRepositoryImpl();
  const userUsecase = new UserUsecaseImpl(userRepository);

  useEffect(() => {
    if (Platform.OS !== 'web') {
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

              setAuthToken(storedToken);
              await fetchUser(storedToken);
              setToken(storedToken);
              setRefreshTimerLogic(storedToken);
              setIsAuthenticated(true);
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
    } else {
      setIsAuthenticated(false);
    }

    return () => {
      if (refreshTimer) {
        clearTimeout(refreshTimer);
      }
    };
  }, []);

  const register = async (userData: Omit<User, 'id'>) => {
    const newToken = await loginUsecase.registerUser(userData);
    setAuthToken(newToken);
    await fetchUser(newToken);
    if (Platform.OS !== 'web') {
      await SecureStore.setItemAsync('token', newToken);
    }
    setToken(newToken);
    setRefreshTimerLogic(newToken);
    setIsAuthenticated(true);
  };

  const login = async (userData: Omit<User, 'id'>) => {
    const newToken = await loginUsecase.loginUser(userData);
    setAuthToken(newToken);
    await fetchUser(newToken);
    if (Platform.OS !== 'web') {
      await SecureStore.setItemAsync('token', newToken);
    }
    setToken(newToken);
    setRefreshTimerLogic(newToken);
    setIsAuthenticated(true);
  };

  const refreshToken = async (token: string) => {
    try {
      const newToken = await loginUsecase.refreshToken(token);
      if (Platform.OS !== 'web') {
        await SecureStore.setItemAsync('token', newToken);
      }
      setToken(newToken);
      setAuthToken(newToken);
      setIsAuthenticated(true);
    } catch (_error) {
      Alert.alert(
        'Oops, algo deu errado',
        '',
        [{ text: 'OK', onPress: () => logoutAndLeave() }],
        { cancelable: false },
      );
    }
  };

  const logout = async () => {
    if (Platform.OS !== 'web') {
      await SecureStore.deleteItemAsync('token');
    }
    setToken(null);
    setAuthToken(null);
    setRefreshTimer(null);
    setIsAuthenticated(false);
    setAuthUser(null);
  };

  const logoutAndLeave = async () => {
    logout();
    navigationRef.current?.dispatch(
      CommonActions.reset({
        index: 0,
        routes: [{ name: 'Login' }],
      }),
    );
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
    const decodedToken = decodeToken(token);
    if (decodedToken != null) {
      const userList: User[] = await userUsecase.getUsersById(decodedToken.id);

      if (!userList || userList.length === 0) {
        throw new Error('user not found');
      }

      setAuthUser(userList[0]);
    } else {
      throw new Error('invalid token');
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
      value={{
        token,
        authUser,
        isAuthenticated,
        register,
        login,
        logout,
        logoutAndLeave,
        setAuthUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export { AuthContext, AuthProvider };
