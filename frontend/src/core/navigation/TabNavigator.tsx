import { Ionicons } from '@expo/vector-icons';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import React from 'react';

import { CreatePostStack } from './CreatePostStack';
import { CurrentUserProfileStack } from './CurrentUserProfileStack';
import { CustomTabBar } from './CustomTabBar';
import { FeedStack } from './FeedStack';
import { TabParamList } from './types';

const Tab = createBottomTabNavigator<TabParamList>();

export function TabNavigator() {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: keyof typeof Ionicons.glyphMap;

          switch (route.name) {
            case 'FeedStack':
              iconName = focused ? 'home' : 'home-outline';
              break;
            case 'CreatePostStack':
              iconName = focused ? 'add-circle' : 'add-circle-outline';
              break;
            case 'CurrentUserProfileStack':
              iconName = focused ? 'person-circle' : 'person-circle-outline';
              break;
            default:
              iconName = focused ? 'home' : 'home-outline';
          }

          return <Ionicons name={iconName} size={size} color={color} />;
        },
        headerShown: false,
      })}
      tabBar={(props) => <CustomTabBar {...props} />}
    >
      <Tab.Screen name="FeedStack" component={FeedStack} />
      <Tab.Screen name="CreatePostStack" component={CreatePostStack} />
      <Tab.Screen
        name="CurrentUserProfileStack"
        component={CurrentUserProfileStack}
      />
    </Tab.Navigator>
  );
}
