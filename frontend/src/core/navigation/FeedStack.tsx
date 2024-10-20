import { createStackNavigator } from '@react-navigation/stack';
import React from 'react';

import { FeedScreen } from '@/src/features/posts/presentation/screens/FeedScreen';

import { FeedStackParamList } from './types';

const Stack = createStackNavigator<FeedStackParamList>();

export const FeedStack = () => {
  return (
    <Stack.Navigator
      screenOptions={{ headerShown: false }}
      initialRouteName="Feed"
    >
      <Stack.Screen name="Feed" component={FeedScreen} />
    </Stack.Navigator>
  );
};
