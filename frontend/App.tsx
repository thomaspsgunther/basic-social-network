import {
  NavigationContainer,
  NavigationContainerRef,
} from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import React, { useContext } from 'react';
import { Provider } from 'react-redux';

import { AuthContext, AuthProvider } from './src/core/context/AuthContext';
import { RootStackParamList } from './src/core/navigation/types';
import store from './src/core/redux/store';
import LoginScreen from './src/features/login/presentation/screens/LoginScreen';
import HomeScreen from './src/features/shared/presentation/screens/HomeScreen';
import LoadingScreen from './src/features/shared/presentation/screens/LoadingScreen';

const navigationRef =
  React.createRef<NavigationContainerRef<RootStackParamList>>();
const Stack = createStackNavigator<RootStackParamList>();

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <AuthProvider navigationRef={navigationRef}>
        <MainNavigator />
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

  return (
    <NavigationContainer ref={navigationRef}>
      <Stack.Navigator initialRouteName={isAuthenticated ? 'Home' : 'Login'}>
        <Stack.Screen name="Home" component={HomeScreen} />
        <Stack.Screen name="Login" component={LoginScreen} />
      </Stack.Navigator>
    </NavigationContainer>
  );
}

export default App;
