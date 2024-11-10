import { Ionicons } from '@expo/vector-icons';
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

import {
  IconDropdown,
  IconDropdownOption,
} from '@/src/core/components/IconDropdown';
import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { CurrentUserProfileStackScreenProps } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { Post } from '@/src/features/shared/data/models/Post';
import { User } from '@/src/features/shared/data/models/User';

import { UserRepositoryImpl } from '../../data/repositories/UserRepositoryImpl';
import { UserUsecaseImpl } from '../../domain/usecases/UserUsecase';

export const CurrentUserProfileScreen: React.FC<
  CurrentUserProfileStackScreenProps<'CurrentUserProfile'>
> = ({ navigation }) => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isLoadingFollowers, setIsLoadingFollowers] = useState<boolean>(false);
  const [isLoadingFollowed, setIsLoadingFollowed] = useState<boolean>(false);
  const [user, setUser] = useState<User>();
  const [isLoadingPosts, setIsLoadingPosts] = useState<boolean>(false);
  const [noMorePosts, setNoMorePosts] = useState<boolean>(false);
  const [posts, setPosts] = useState<Post[]>([]);
  const userRepository = new UserRepositoryImpl();
  const userUsecase = new UserUsecaseImpl(userRepository);

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error(
      'currentuserprofilescreen must be used within an authprovider',
    );
  }

  const { authUser, logoutAndLeave } = context;

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  useEffect(() => {
    if (!noMorePosts && posts && posts.length === 0) {
      initPosts();
    }
  }, [posts]);

  const initPosts = async () => {
    if (isLoading) {
      return;
    }
    setIsLoading(true);
    try {
      if (authUser) {
        const users: User[] = await userUsecase.getUsersById(authUser.id);

        if (users) {
          setUser(users[0]);

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
        }
      } else {
        throw new Error('missing authuser');
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
    setPosts([]);
    try {
      if (authUser) {
        const users: User[] = await userUsecase.getUsersById(authUser.id);

        if (users) {
          setUser(users[0]);

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
          throw new Error('missing authuser');
        }
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

      const newPosts: Post[] = await userUsecase.listUserPosts(
        user.id,
        15,
        cursor,
      );

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

  const goToFollowers = async () => {
    if (user) {
      setIsLoadingFollowers(true);
      try {
        const followers: User[] = await userUsecase.getUserFollowers(user.id);

        if (followers) {
          setIsLoadingFollowers(false);
          navigation.push('UserList', {
            users: followers,
            title: 'Seguidores',
          });
        } else {
          setIsLoadingFollowers(false);
        }
      } catch (_error) {
        setIsLoadingFollowers(false);
        Alert.alert('Oops, algo deu errado');
      }
    }
  };

  const goToFollowed = async () => {
    if (user) {
      setIsLoadingFollowed(true);
      try {
        const followed: User[] = await userUsecase.getUserFollowed(user.id);

        if (followed) {
          setIsLoadingFollowed(false);
          navigation.push('UserList', { users: followed, title: 'Seguindo' });
        } else {
          setIsLoadingFollowed(false);
        }
      } catch (_error) {
        setIsLoadingFollowed(false);
        Alert.alert('Oops, algo deu errado');
      }
    }
  };

  const goToPost = async (id: string) => {
    navigation.push('PostDetail', { postId: id });
  };

  const goToEdit = async () => {
    navigation.push('UserEdit');
  };

  const options: IconDropdownOption[] = [
    {
      label: 'Sair',
      iconName: 'log-out-outline',
      onSelect: async () => {
        logoutAndLeave();
      },
    },
  ];

  return (
    <View style={currentTheme.container}>
      {!isLoading ? (
        user &&
        authUser && (
          <>
            <View style={styles.listHeaderTopRow}>
              <Text style={currentTheme.titleText}>{user.username}</Text>

              <View style={currentTheme.row}>
                <View style={styles.icon}>
                  <TouchableOpacity onPress={() => handleReload()}>
                    <Ionicons
                      name="reload"
                      size={34}
                      color={currentColors.icon}
                    ></Ionicons>
                  </TouchableOpacity>
                </View>

                <IconDropdown options={options}></IconDropdown>
              </View>
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
                    {authUser.avatar ? (
                      <Image
                        source={{
                          uri: `data:image/jpeg;base64,${authUser.avatar}`,
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

                    {!isLoadingFollowers ? (
                      <TouchableOpacity
                        style={styles.infoColumn}
                        onPress={() => goToFollowers()}
                        disabled={(user.followerCount ?? 0) === 0}
                      >
                        <Text
                          style={currentTheme.textBold}
                        >{`${user.followerCount ?? 0}`}</Text>

                        <Text style={currentTheme.text}>seguidores</Text>
                      </TouchableOpacity>
                    ) : (
                      <ActivityIndicator
                        size="large"
                        color={currentColors.icon}
                      />
                    )}

                    {!isLoadingFollowed ? (
                      <TouchableOpacity
                        style={styles.infoColumn}
                        onPress={() => goToFollowed()}
                        disabled={(user.followedCount ?? 0) === 0}
                      >
                        <Text
                          style={currentTheme.textBold}
                        >{`${user.followedCount ?? 0}`}</Text>

                        <Text style={currentTheme.text}>seguindo</Text>
                      </TouchableOpacity>
                    ) : (
                      <ActivityIndicator
                        size="large"
                        color={currentColors.icon}
                      />
                    )}
                  </View>

                  {authUser.fullName && (
                    <Text style={currentTheme.textBold}>
                      {authUser.fullName}
                    </Text>
                  )}

                  {authUser.description && (
                    <Text style={currentTheme.text}>
                      {authUser.description}
                    </Text>
                  )}

                  <View style={styles.buttonContainer}>
                    <TouchableOpacity
                      style={currentTheme.button}
                      onPress={() => goToEdit()}
                    >
                      <Text style={currentTheme.buttonText}>Editar perfil</Text>
                    </TouchableOpacity>
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
    paddingTop: 10,
  },
  icon: {
    paddingRight: 20,
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
    marginTop: 50,
    paddingBottom: 10,
    paddingLeft: 20,
    paddingRight: 10,
    width: '100%',
  },
  loadingContainer: {
    paddingVertical: 5,
  },
  postContainer: {
    margin: 1,
  },
  userInfoRow: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingBottom: 10,
  },
});
