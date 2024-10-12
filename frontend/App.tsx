import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import {
  NavigationContainer,
  NavigationContainerRef,
} from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import React, { useContext } from 'react';
import { Provider } from 'react-redux';

import { AuthContext, AuthProvider } from './src/core/context/AuthContext';
import { ThemeProvider } from './src/core/context/ThemeContext';
import {
  HomeTabParamList,
  RootStackParamList,
} from './src/core/navigation/types';
import store from './src/core/redux/store';
import LoginScreen from './src/features/login/presentation/screens/LoginScreen';
import FeedScreen from './src/features/shared/presentation/screens/FeedScreen';
import LoadingScreen from './src/features/shared/presentation/screens/LoadingScreen';

const navigationRef =
  React.createRef<NavigationContainerRef<RootStackParamList>>();
const Stack = createStackNavigator<RootStackParamList>();
const Tab = createBottomTabNavigator<HomeTabParamList>();

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <AuthProvider navigationRef={navigationRef}>
        <ThemeProvider>
          <MainNavigator />
        </ThemeProvider>
      </AuthProvider>
    </Provider>
  );
};

function MainNavigator() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error('mainnavigator must be used within an authprovider');
  }

  const { isAuthenticated } = context;

  if (isAuthenticated === null) {
    return <LoadingScreen />;
  }

  function TabNavigator() {
    return (
      <Tab.Navigator>
        <Tab.Screen
          name="Feed"
          component={FeedScreen}
          options={{ headerShown: false }}
        />
      </Tab.Navigator>
    );
  }

  return (
    <NavigationContainer ref={navigationRef}>
      <Stack.Navigator
        screenOptions={{ headerShown: false }}
        initialRouteName={isAuthenticated ? 'Home' : 'Login'}
      >
        <Stack.Screen
          name="Home"
          component={TabNavigator}
          options={{ headerShown: false }}
        />
        <Stack.Screen
          name="Login"
          component={LoginScreen}
          options={{ headerShown: false }}
        />
      </Stack.Navigator>
    </NavigationContainer>
  );
}

export default App;
