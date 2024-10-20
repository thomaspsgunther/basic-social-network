import { createStackNavigator } from '@react-navigation/stack';
import React from 'react';

import { CreatePostScreen } from '@/src/features/posts/presentation/screens/CreatePostScreen';
import { PostDetailScreen } from '@/src/features/posts/presentation/screens/PostDetailScreen';

import { CreatePostStackParamList } from './types';

const Stack = createStackNavigator<CreatePostStackParamList>();

export const CreatePostStack = () => {
  return (
    <Stack.Navigator
      screenOptions={{ headerShown: false }}
      initialRouteName="CreatePost"
    >
      <Stack.Screen name="CreatePost" component={CreatePostScreen} />
      <Stack.Screen name="PostDetail" component={PostDetailScreen} />
    </Stack.Navigator>
  );
};
