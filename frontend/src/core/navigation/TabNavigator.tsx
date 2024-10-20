import { Ionicons } from '@expo/vector-icons';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import React from 'react';

import { CreatePostStack } from './CreatePostStack';
import CustomTabBar from './CustomTabBar';
import { FeedStack } from './FeedStack';
import { TabParamList } from './types';

const Tab = createBottomTabNavigator<TabParamList>();

export function TabNavigator() {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ color, size }) => {
          let iconName: keyof typeof Ionicons.glyphMap;

          switch (route.name) {
            case 'FeedStack':
              iconName = 'list';
              break;
            case 'CreatePostStack':
              iconName = 'add-circle-outline';
              break;
            default:
              iconName = 'home';
          }

          return <Ionicons name={iconName} size={size} color={color} />;
        },
        headerShown: false,
      })}
      tabBar={(props) => <CustomTabBar {...props} />}
    >
      <Tab.Screen name="FeedStack" component={FeedStack} />
      <Tab.Screen name="CreatePostStack" component={CreatePostStack} />
    </Tab.Navigator>
  );
}
