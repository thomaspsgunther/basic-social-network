import { StackNavigationProp } from '@react-navigation/stack';

export type RootStackParamList = {
  Login: undefined;
  Register: undefined;
  Tabs: undefined;
};

export type TabParamList = {
  FeedStack: undefined;
};

export type FeedStackParamList = {
  Feed: undefined;
};

export type LoginScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'Login'
>;

export type RegisterScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'Register'
>;

export type FeedScreenNavigationProp = StackNavigationProp<
  FeedStackParamList,
  'Feed'
>;
