import { createStackNavigator } from '@react-navigation/stack';
import React from 'react';

import { FeedScreen } from '@/src/features/shared/presentation/screens/FeedScreen';

import { FeedStackParamList } from './types';

const Stack = createStackNavigator<FeedStackParamList>();

export const FeedStackNavigator = () => {
  return (
    <Stack.Navigator>
      <Stack.Screen name="Feed" component={FeedScreen} />
    </Stack.Navigator>
  );
};
