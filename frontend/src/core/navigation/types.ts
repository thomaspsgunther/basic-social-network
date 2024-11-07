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
  SearchUserStack: undefined;
  CreatePostStack: undefined;
  CurrentUserProfileStack: undefined;
};

export type FeedStackParamList = {
  Feed: undefined;
  PostDetail: { postId: string };
  UserProfile: { userId: string };
};

export type CreatePostStackParamList = {
  CreatePost: undefined;
  PostDetail: { postId: string };
  UserProfile: { userId: string };
};

export type CurrentUserProfileStackParamList = {
  CurrentUserProfile: undefined;
  EditUser: undefined;
  PostDetail: { postId: string };
  UserProfile: { userId: string };
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

export type CurrentUserProfileStackScreenProps<
  T extends keyof CurrentUserProfileStackParamList,
> = StackScreenProps<CurrentUserProfileStackParamList, T>;
