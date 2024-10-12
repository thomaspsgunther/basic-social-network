import { BottomTabNavigationProp } from '@react-navigation/bottom-tabs';
import { StackNavigationProp } from '@react-navigation/stack';

export type RootStackParamList = {
  Login: undefined;
  Home: undefined;
};

export type HomeTabParamList = {
  Feed: undefined;
};

export type LoginScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'Login'
>;

export type FeedScreenNavigationProp = BottomTabNavigationProp<
  HomeTabParamList,
  'Feed'
>;
