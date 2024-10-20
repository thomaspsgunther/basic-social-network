import React from 'react';
import { Text, View } from 'react-native';

import { useTheme } from '@/src/core/context/ThemeContext';
import { CreatePostStackScreenProps } from '@/src/core/navigation/types';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';

export const CreatePostScreen: React.FC<
  CreatePostStackScreenProps<'CreatePost'>
> = () => {
  const { isDarkMode } = useTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;

  return (
    <View style={currentTheme.container}>
      <Text style={currentTheme.text}>Criar Post</Text>
    </View>
  );
};
