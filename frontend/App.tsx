import { NavigationContainerRef } from '@react-navigation/native';
import React from 'react';

import { AuthProvider } from './src/core/context/AuthContext';
import { ThemeProvider } from './src/core/context/ThemeContext';
import { MainNavigator } from './src/core/navigation/MainNavigator';
import { RootStackParamList } from './src/core/navigation/types';

const navigationRef =
  React.createRef<NavigationContainerRef<RootStackParamList>>();

const App: React.FC = () => {
  return (
    <AuthProvider navigationRef={navigationRef}>
      <ThemeProvider>
        <MainNavigator navigationRef={navigationRef} />
      </ThemeProvider>
    </AuthProvider>
  );
};

export default App;
