import { createStackNavigator } from '@react-navigation/stack';
import React from 'react';

import { PostCommentsScreen } from '@/src/features/comments/presentation/screens/PostCommentsScreen';
import { PostDetailScreen } from '@/src/features/posts/presentation/screens/PostDetailScreen';
import { CurrentUserProfileScreen } from '@/src/features/users/presentation/screens/CurrentUserProfileScreen';
import { UserEditScreen } from '@/src/features/users/presentation/screens/UserEditScreen';
import { UserListScreen } from '@/src/features/users/presentation/screens/UserListScreen';
import { UserProfileScreen } from '@/src/features/users/presentation/screens/UserProfileScreen';

import { CurrentUserProfileStackParamList } from './types';

const Stack = createStackNavigator<CurrentUserProfileStackParamList>();

export const CurrentUserProfileStack = () => {
  return (
    <Stack.Navigator
      screenOptions={{ headerShown: false }}
      initialRouteName="CurrentUserProfile"
    >
      <Stack.Screen
        name="CurrentUserProfile"
        component={CurrentUserProfileScreen}
      />
      <Stack.Screen name="UserEdit" component={UserEditScreen} />
      <Stack.Screen name="PostDetail" component={PostDetailScreen} />
      <Stack.Screen name="PostComments" component={PostCommentsScreen} />
      <Stack.Screen name="UserProfile" component={UserProfileScreen} />
      <Stack.Screen name="UserList" component={UserListScreen} />
    </Stack.Navigator>
  );
};
