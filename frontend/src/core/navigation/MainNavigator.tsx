import {
  NavigationContainer,
  NavigationContainerRef,
} from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import React, { useContext } from 'react';

import { LoadingScreen } from '@/src/features/login/presentation/screens/LoadingScreen';
import { LoginScreen } from '@/src/features/login/presentation/screens/LoginScreen';
import { RegisterScreen } from '@/src/features/login/presentation/screens/RegisterScreen';

import { AuthContext } from '../context/AuthContext';
import { TabNavigator } from './TabNavigator';
import { RootStackParamList } from './types';

const Stack = createStackNavigator<RootStackParamList>();

interface MainNavigatorProps {
  navigationRef: React.RefObject<NavigationContainerRef<RootStackParamList>>;
}

export const MainNavigator: React.FC<MainNavigatorProps> = ({
  navigationRef,
}) => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('mainnavigator must be used within an authprovider');
  }

  const { isAuthenticated } = context;

  if (isAuthenticated === null) {
    return <LoadingScreen />;
  }

  return (
    <NavigationContainer ref={navigationRef}>
      <Stack.Navigator
        screenOptions={{ headerShown: false }}
        initialRouteName={isAuthenticated ? 'Tabs' : 'Login'}
      >
        <Stack.Screen name="Tabs" component={TabNavigator} />
        <Stack.Screen name="Login" component={LoginScreen} />
        <Stack.Screen name="Register" component={RegisterScreen} />
      </Stack.Navigator>
    </NavigationContainer>
  );
};
