import { BottomTabScreenProps } from '@react-navigation/bottom-tabs';
import { CompositeScreenProps } from '@react-navigation/native';
import { StackScreenProps } from '@react-navigation/stack';

import { User } from '@/src/features/shared/data/models/User';

export type RootStackParamList = {
  Login: undefined;
  Register: undefined;
  Tabs: undefined;
};

export type TabParamList = {
  FeedStack: undefined;
  UserSearchStack: undefined;
  CreatePostStack: undefined;
  CurrentUserProfileStack: undefined;
};

export type FeedStackParamList = {
  Feed: undefined;
  PostComments: { postId: string };
  UserProfile: { userId: string };
  UserList: { users: User[]; title?: string };
  PostDetail: { postId: string; editing?: boolean };
};

export type UserSearchStackParamList = {
  UserSearch: undefined;
  UserProfile: { userId: string };
  UserList: { users: User[]; title?: string };
  PostDetail: { postId: string; editing?: boolean };
  PostComments: { postId: string };
};

export type CreatePostStackParamList = {
  CreatePost: undefined;
  PostDetail: { postId: string; editing?: boolean };
  PostComments: { postId: string };
  UserProfile: { userId: string };
  UserList: { users: User[]; title?: string };
};

export type CurrentUserProfileStackParamList = {
  CurrentUserProfile: undefined;
  UserEdit: undefined;
  PostDetail: { postId: string; editing?: boolean };
  PostComments: { postId: string };
  UserProfile: { userId: string };
  UserList: { users: User[]; title?: string };
};

export type RootStackScreenProps<T extends keyof RootStackParamList> =
  StackScreenProps<RootStackParamList, T>;

export type TabScreenProps<T extends keyof TabParamList> = CompositeScreenProps<
  BottomTabScreenProps<TabParamList, T>,
  StackScreenProps<RootStackParamList>
>;

export type FeedStackScreenProps<T extends keyof FeedStackParamList> =
  StackScreenProps<FeedStackParamList, T>;

export type UserSearchStackScreenProps<
  T extends keyof UserSearchStackParamList,
> = StackScreenProps<UserSearchStackParamList, T>;

export type CreatePostStackScreenProps<
  T extends keyof CreatePostStackParamList,
> = StackScreenProps<CreatePostStackParamList, T>;

export type CurrentUserProfileStackScreenProps<
  T extends keyof CurrentUserProfileStackParamList,
> = StackScreenProps<CurrentUserProfileStackParamList, T>;
