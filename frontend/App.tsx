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
import ErrorAlert from './src/core/errors/ErrorAlert';
import { FeedStack } from './src/core/navigation/FeedStack';
import { RootStackParamList, TabParamList } from './src/core/navigation/types';
import { store } from './src/core/redux/store';
import { LoadingScreen } from './src/features/login/presentation/screens/LoadingScreen';
import { LoginScreen } from './src/features/login/presentation/screens/LoginScreen';
import { RegisterScreen } from './src/features/login/presentation/screens/RegisterScreen';

const navigationRef =
  React.createRef<NavigationContainerRef<RootStackParamList>>();
const Stack = createStackNavigator<RootStackParamList>();
const Tab = createBottomTabNavigator<TabParamList>();

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <ErrorAlert />
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
          name="FeedStack"
          component={FeedStack}
          options={{ headerShown: false }}
        />
      </Tab.Navigator>
    );
  }

  return (
    <NavigationContainer ref={navigationRef}>
      <Stack.Navigator
        screenOptions={{ headerShown: false }}
        initialRouteName={isAuthenticated ? 'Tabs' : 'Login'}
      >
        <Stack.Screen name="Tabs" component={TabNavigator} />
        <Stack.Screen
          name="Login"
          component={LoginScreen}
          options={{ headerShown: false }}
        />
        <Stack.Screen
          name="Register"
          component={RegisterScreen}
          options={{ headerShown: false }}
        />
      </Stack.Navigator>
    </NavigationContainer>
  );
}

export default App;
