import { BottomTabScreenProps } from '@react-navigation/bottom-tabs';
import { CompositeScreenProps } from '@react-navigation/native';
import { StackScreenProps } from '@react-navigation/stack';

export type RootStackParamList = {
  Login: undefined;
  Register: undefined;
  Tabs: undefined;
};

export type TabParamList = {
  FeedStack: undefined;
  CreatePostStack: undefined;
  SearchUserStack: undefined;
  CurrentUserProfileStack: undefined;
};

export type FeedStackParamList = {
  Feed: undefined;
  PostDetail: { postId: string };
  UserProfile: undefined;
};

export type CreatePostStackParamList = {
  CreatePost: undefined;
  PostDetail: { postId: string };
  UserProfile: undefined;
};

export type RootStackScreenProps<T extends keyof RootStackParamList> =
  StackScreenProps<RootStackParamList, T>;

export type TabScreenProps<T extends keyof TabParamList> = CompositeScreenProps<
  BottomTabScreenProps<TabParamList, T>,
  StackScreenProps<RootStackParamList>
>;

export type FeedStackScreenProps<T extends keyof FeedStackParamList> =
  StackScreenProps<FeedStackParamList, T>;

export type CreatePostStackScreenProps<
  T extends keyof CreatePostStackParamList,
> = StackScreenProps<CreatePostStackParamList, T>;
