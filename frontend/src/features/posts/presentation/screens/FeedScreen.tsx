import React from 'react';
import { Text, View } from 'react-native';

import { useAppTheme } from '@/src/core/context/ThemeContext';
import { FeedStackScreenProps } from '@/src/core/navigation/types';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';

export const FeedScreen: React.FC<FeedStackScreenProps<'Feed'>> = () => {
  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;

  return (
    <View style={currentTheme.container}>
      <Text style={currentTheme.text}>Feed</Text>
    </View>
  );
};
