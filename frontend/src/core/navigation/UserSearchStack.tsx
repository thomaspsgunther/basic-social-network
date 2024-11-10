import { createStackNavigator } from '@react-navigation/stack';
import React from 'react';

import { PostCommentsScreen } from '@/src/features/comments/presentation/screens/PostCommentsScreen';
import { PostDetailScreen } from '@/src/features/posts/presentation/screens/PostDetailScreen';
import { UserListScreen } from '@/src/features/users/presentation/screens/UserListScreen';
import { UserProfileScreen } from '@/src/features/users/presentation/screens/UserProfileScreen';
import { UserSearchScreen } from '@/src/features/users/presentation/screens/UserSearchScreen';

import { UserSearchStackParamList } from './types';

const Stack = createStackNavigator<UserSearchStackParamList>();

export const UserSearchStack = () => {
  return (
    <Stack.Navigator
      screenOptions={{ headerShown: false }}
      initialRouteName="UserSearch"
    >
      <Stack.Screen name="UserSearch" component={UserSearchScreen} />
      <Stack.Screen name="UserProfile" component={UserProfileScreen} />
      <Stack.Screen name="UserList" component={UserListScreen} />
      <Stack.Screen name="PostDetail" component={PostDetailScreen} />
      <Stack.Screen name="PostComments" component={PostCommentsScreen} />
    </Stack.Navigator>
  );
};
