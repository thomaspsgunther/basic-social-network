import { Ionicons } from '@expo/vector-icons';
import { RouteProp, useNavigation, useRoute } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import React, { useContext, useEffect, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  FlatList,
  Image,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';

import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { FeedStackParamList } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { Post } from '@/src/features/shared/data/models/Post';
import { User } from '@/src/features/shared/data/models/User';

import { UserRepositoryImpl } from '../../data/repositories/UserRepositoryImpl';
import { UserUsecaseImpl } from '../../domain/usecases/UserUsecase';

export const UserProfileScreen: React.FC = () => {
  const navigation =
    useNavigation<StackNavigationProp<FeedStackParamList, 'UserProfile'>>();

  const route = useRoute<RouteProp<FeedStackParamList, 'UserProfile'>>();
  const { userId } = route.params;

  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [user, setUser] = useState<User>();
  const [isFollowing, setIsFollowing] = useState<boolean>(false);
  const [isLoadingFollow, setIsLoadingFollow] = useState<boolean>(false);
  const [isLoadingPosts, setIsLoadingPosts] = useState<boolean>(false);
  const [noMorePosts, setNoMorePosts] = useState<boolean>(false);
  const [posts, setPosts] = useState<Post[]>([]);
  const userRepository = new UserRepositoryImpl();
  const userUsecase = new UserUsecaseImpl(userRepository);

  const canGoBack = navigation.canGoBack();

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('userprofilescreen must be used within an authprovider');
  }

  const { authUser } = context;

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  useEffect(() => {
    if (!noMorePosts && posts && posts.length === 0) {
      initProfile();
    }
  }, [posts]);

  const initProfile = async () => {
    if (isLoading) {
      return;
    }
    setIsLoading(true);
    try {
      const users: User[] = await userUsecase.getUsersById(userId);

      if (users && authUser) {
        setUser(users[0]);

        const doesFollow: boolean = await userUsecase.userFollowsUser(
          authUser!.id,
          users[0].id,
        );

        if (doesFollow) {
          setIsFollowing(true);
        }

        const initialPosts: Post[] = await userUsecase.listUserPosts(
          users[0].id,
          15,
        );

        if (initialPosts && initialPosts.length > 0) {
          setIsLoading(false);
          setPosts(initialPosts);
        } else {
          setIsLoading(false);
          setNoMorePosts(true);
        }
      } else {
        throw new Error('missing user or authuser');
      }
    } catch (_error) {
      setIsLoading(false);
      setNoMorePosts(true);
      Alert.alert('Oops, algo deu errado');
    }
  };

  const handleReload = async () => {
    if (isLoading) {
      return;
    }
    setIsLoading(true);
    setNoMorePosts(false);
    setUser(undefined);
    setIsFollowing(false);
    setPosts([]);
    try {
      const users: User[] = await userUsecase.getUsersById(userId);

      if (users && authUser) {
        setUser(users[0]);

        const doesFollow: boolean = await userUsecase.userFollowsUser(
          authUser!.id,
          users[0].id,
        );

        if (doesFollow) {
          setIsFollowing(true);
        }

        const initialPosts: Post[] = await userUsecase.listUserPosts(
          users[0].id,
          15,
        );

        if (initialPosts && initialPosts.length > 0) {
          setIsLoading(false);
          setPosts(initialPosts);
          setNoMorePosts(false);
        } else {
          setIsLoading(false);
          setNoMorePosts(true);
        }
      } else {
        throw new Error('missing user or authuser');
      }
    } catch (_error) {
      setIsLoading(false);
      setNoMorePosts(false);
      Alert.alert('Oops, algo deu errado');
    }
  };

  const loadPosts = async () => {
    if (!user || isLoading || isLoadingPosts || noMorePosts) {
      return;
    }
    setIsLoadingPosts(true);
    try {
      const currentPosts: Post[] = [...posts];

      const lastPost = currentPosts[currentPosts.length - 1];

      const cursor: string = Buffer.from(
        `${lastPost.createdAt},${lastPost.id}`,
      ).toString('base64');

      const newPosts = await userUsecase.listUserPosts(user.id, 15, cursor);

      if (newPosts && newPosts.length > 0) {
        setPosts((prevPosts) => [...prevPosts, ...newPosts]);
        setTimeout(() => {
          setIsLoadingPosts(false);
        }, 5);
      } else {
        setNoMorePosts(true);
        setTimeout(() => {
          setIsLoadingPosts(false);
        }, 5);
      }
    } catch (_error) {
      setNoMorePosts(true);
      Alert.alert('Oops, algo deu errado');
      setTimeout(() => {
        setIsLoadingPosts(false);
      }, 5);
    }
  };

  const handleFollow = async () => {
    setIsLoadingFollow(true);
    try {
      if (user && authUser) {
        if (!isFollowing) {
          const didFollow: boolean = await userUsecase.followUser(
            authUser.id,
            user.id,
          );

          if (didFollow) {
            user.followerCount = (user.followerCount ?? 0) + 1;
            setIsLoadingFollow(false);
            setIsFollowing(true);
          }
        } else {
          const didUnfollow: boolean = await userUsecase.unfollowUser(
            authUser.id,
            user.id,
          );

          if (didUnfollow) {
            user.followerCount = (user.followerCount ?? 0) - 1;
            setIsLoadingFollow(false);
            setIsFollowing(false);
          }
        }
      } else {
        throw new Error('missing user or authuser');
      }
    } catch (_error) {
      setIsLoadingFollow(false);
      Alert.alert('Oops, algo deu errado');
    }
  };

  const goToPost = async (id: string) => {
    navigation.push('PostDetail', { postId: id });
  };

  return (
    <View style={currentTheme.container}>
      {!isLoading && !user && canGoBack && (
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={currentTheme.backButton}
        >
          <Ionicons name="arrow-back" size={40} color={currentColors.icon} />
        </TouchableOpacity>
      )}

      {!isLoading ? (
        user && (
          <>
            <View style={styles.listHeaderTopRow}>
              <View style={currentTheme.row}>
                {canGoBack && (
                  <TouchableOpacity onPress={() => navigation.goBack()}>
                    <Ionicons
                      name="arrow-back"
                      size={40}
                      color={currentColors.icon}
                    />
                  </TouchableOpacity>
                )}

                <Text
                  style={currentTheme.titleText}
                >{`   ${user.username}`}</Text>
              </View>

              <TouchableOpacity onPress={() => handleReload()}>
                <Ionicons
                  name="reload"
                  size={34}
                  color={currentColors.icon}
                ></Ionicons>
              </TouchableOpacity>
            </View>

            <FlatList
              data={posts}
              keyExtractor={(post) => post.id}
              renderItem={({ item }: { item: Post }) => (
                <View style={styles.postContainer}>
                  <TouchableOpacity onPress={() => goToPost(item.id)}>
                    <Image
                      source={{ uri: `data:image/jpeg;base64,${item.image}` }}
                      style={styles.image}
                      resizeMode="contain"
                    />
                  </TouchableOpacity>
                </View>
              )}
              numColumns={3}
              onEndReached={() => loadPosts()}
              showsVerticalScrollIndicator={false}
              contentContainerStyle={styles.flatListContainer}
              ListHeaderComponent={
                <View style={currentTheme.userHeader}>
                  <View style={styles.userInfoRow}>
                    {user.avatar ? (
                      <Image
                        source={{
                          uri: `data:image/jpeg;base64,${user.avatar}`,
                        }}
                        style={styles.avatar}
                        resizeMode="contain"
                      />
                    ) : (
                      <View style={styles.avatarPlaceholder}>
                        <Ionicons
                          name="person-circle-outline"
                          size={100}
                          color="black"
                        ></Ionicons>
                      </View>
                    )}

                    <View style={styles.infoColumn}>
                      <Text
                        style={currentTheme.textBold}
                      >{`${user.postCount ?? 0}`}</Text>

                      <Text style={currentTheme.text}>publicações</Text>
                    </View>

                    <View style={styles.infoColumn}>
                      <Text
                        style={currentTheme.textBold}
                      >{`${user.followerCount ?? 0}`}</Text>

                      <Text style={currentTheme.text}>seguidores</Text>
                    </View>

                    <View style={styles.infoColumn}>
                      <Text
                        style={currentTheme.textBold}
                      >{`${user.followedCount ?? 0}`}</Text>

                      <Text style={currentTheme.text}>seguindo</Text>
                    </View>
                  </View>

                  {user.fullName && (
                    <Text style={currentTheme.textBold}>{user.fullName}</Text>
                  )}

                  {user.description && (
                    <Text style={currentTheme.text}>{user.description}</Text>
                  )}

                  <View style={styles.buttonContainer}>
                    {!isLoadingFollow ? (
                      <TouchableOpacity
                        style={
                          isFollowing
                            ? styles.unfollowButton
                            : currentTheme.button
                        }
                        onPress={() => handleFollow()}
                      >
                        <Text style={currentTheme.buttonText}>
                          {isFollowing ? 'Deixar de seguir' : 'Seguir'}
                        </Text>
                      </TouchableOpacity>
                    ) : (
                      <ActivityIndicator
                        size="large"
                        color={currentColors.icon}
                      />
                    )}
                  </View>
                </View>
              }
            ></FlatList>

            {isLoadingPosts && (
              <ActivityIndicator
                size="large"
                style={styles.loadingContainer}
                color={currentColors.icon}
              />
            )}
          </>
        )
      ) : (
        <ActivityIndicator size="large" color={currentColors.icon} />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  avatar: {
    borderRadius: 100,
    height: 100,
    width: 100,
  },
  avatarPlaceholder: {
    alignItems: 'center',
    backgroundColor: '#ccc' as string,
    borderRadius: 100,
    height: 100,
    justifyContent: 'center',
    width: 100,
  },
  buttonContainer: {
    flexDirection: 'row',
    paddingTop: 10,
  },
  flatListContainer: {
    flexGrow: 1,
  },
  image: {
    height: 135,
    width: 135,
  },
  infoColumn: {
    alignItems: 'center',
    flexDirection: 'column',
    justifyContent: 'center',
  },
  listHeaderTopRow: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingBottom: 20,
    paddingLeft: 20,
    paddingRight: 20,
    paddingTop: 50,
    width: '100%',
  },
  loadingContainer: {
    paddingVertical: 5,
  },
  postContainer: {
    margin: 1,
  },
  unfollowButton: {
    backgroundColor: 'red' as string,
    borderRadius: 5,
    marginBottom: 20,
    marginTop: 10,
    paddingHorizontal: 20,
    paddingVertical: 12,
  },
  userInfoRow: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingBottom: 10,
  },
});
