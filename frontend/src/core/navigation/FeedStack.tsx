import { createStackNavigator } from '@react-navigation/stack';
import React from 'react';

import { PostCommentsScreen } from '@/src/features/comments/presentation/screens/PostCommentsScreen';
import { FeedScreen } from '@/src/features/posts/presentation/screens/FeedScreen';
import { PostDetailScreen } from '@/src/features/posts/presentation/screens/PostDetailScreen';
import { UserProfileScreen } from '@/src/features/users/presentation/screens/UserProfileScreen';

import { FeedStackParamList } from './types';

const Stack = createStackNavigator<FeedStackParamList>();

export const FeedStack = () => {
  return (
    <Stack.Navigator
      screenOptions={{ headerShown: false }}
      initialRouteName="Feed"
    >
      <Stack.Screen name="Feed" component={FeedScreen} />
      <Stack.Screen name="PostComments" component={PostCommentsScreen} />
      <Stack.Screen name="UserProfile" component={UserProfileScreen} />
      <Stack.Screen name="PostDetail" component={PostDetailScreen} />
    </Stack.Navigator>
  );
};
