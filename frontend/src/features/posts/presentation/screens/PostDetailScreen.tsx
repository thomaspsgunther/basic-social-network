import { useNavigation } from '@react-navigation/native';
import React from 'react';
import { Text, View } from 'react-native';

import { useTheme } from '@/src/core/context/ThemeContext';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';

export const PostDetailScreen: React.FC = () => {
  const _navigation = useNavigation();
  const { isDarkMode } = useTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;

  return (
    <View style={currentTheme.container}>
      <Text style={currentTheme.text}>Publicação</Text>
    </View>
  );
};
