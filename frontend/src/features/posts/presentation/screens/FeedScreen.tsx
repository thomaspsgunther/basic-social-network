import { Ionicons } from '@expo/vector-icons';
import { Buffer } from 'buffer';
import React, { useContext, useEffect, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  Image,
  Pressable,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';
import { FlatList } from 'react-native-gesture-handler';

import {
  IconDropdown,
  IconDropdownOption,
} from '@/src/core/components/IconDropdown';
import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { FeedStackScreenProps } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { Post } from '@/src/features/shared/data/models/Post';

import { PostRepositoryImpl } from '../../data/repositories/PostRepositoryImpl';
import { PostUsecaseImpl } from '../../domain/usecases/PostUsecase';

export const FeedScreen: React.FC<FeedStackScreenProps<'Feed'>> = ({
  navigation,
}) => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isLoadingPosts, setIsLoadingPosts] = useState<boolean>(false);
  const [noMorePosts, setNoMorePosts] = useState<boolean>(false);
  const [posts, setPosts] = useState<Post[]>([]);
  const [likedPostIds, setLikedPostIds] = useState<string[]>([]);
  const postRepository = new PostRepositoryImpl();
  const postUsecase = new PostUsecaseImpl(postRepository);

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('feedscreen must be used within an authprovider');
  }

  const { authUser } = context;

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  useEffect(() => {
    if (posts.length <= 5) {
      checkLikes(posts);
    }

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
      const initialPosts: Post[] = await postUsecase.listPosts(5);

      if (initialPosts && initialPosts.length > 0) {
        setIsLoading(false);
        setPosts(initialPosts);
      } else {
        setIsLoading(false);
        setNoMorePosts(true);
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
    setPosts([]);
    setLikedPostIds([]);
    try {
      const initialPosts: Post[] = await postUsecase.listPosts(5);

      if (initialPosts && initialPosts.length > 0) {
        setIsLoading(false);
        setPosts(initialPosts);
        setNoMorePosts(false);
      } else {
        setIsLoading(false);
        setNoMorePosts(true);
      }
    } catch (_error) {
      setIsLoading(false);
      setNoMorePosts(true);
      Alert.alert('Oops, algo deu errado');
    }
  };

  const loadPosts = async () => {
    if (isLoading || isLoadingPosts || noMorePosts) {
      return;
    }
    setIsLoadingPosts(true);
    try {
      const currentPosts: Post[] = [...posts];

      const lastPost = currentPosts[currentPosts.length - 1];

      const cursor: string = Buffer.from(
        `${lastPost.createdAt},${lastPost.id}`,
      ).toString('base64');

      const newPosts = await postUsecase.listPosts(5, cursor);

      if (newPosts && newPosts.length > 0) {
        checkLikes(newPosts);
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

  const checkLikes = async (posts: Post[]) => {
    try {
      if (authUser) {
        let newLikedPostIds: string[] = [...likedPostIds];

        for (const post of posts) {
          const didLike: boolean = await postUsecase.checkIfUserLikedPost(
            authUser.id,
            post.id,
          );

          if (didLike && !newLikedPostIds.includes(post.id)) {
            newLikedPostIds.push(post.id);
          } else if (!didLike && newLikedPostIds.includes(post.id)) {
            newLikedPostIds = newLikedPostIds.filter(
              (postId) => postId !== post.id,
            );
          }
        }

        if (newLikedPostIds.length !== likedPostIds.length) {
          setLikedPostIds(newLikedPostIds);
        }
      } else {
        throw new Error('missing authuser');
      }
    } catch (_error) {
      Alert.alert('Oops, algo deu errado');
    }
  };

  const handleLike = async (post: Post) => {
    try {
      if (authUser) {
        let newLikedPostIds: string[] = [];

        if (likedPostIds.includes(post.id)) {
          const didUnlike: boolean = await postUsecase.unlikePost(
            authUser.id,
            post.id,
          );
          if (didUnlike) {
            post.likeCount = (post.likeCount ?? 0) - 1;
            newLikedPostIds = likedPostIds.filter(
              (postId) => postId !== post.id,
            );
          }
        } else {
          const didLike: boolean = await postUsecase.likePost(
            authUser.id,
            post.id,
          );
          if (didLike) {
            post.likeCount = (post.likeCount ?? 0) + 1;
            newLikedPostIds = [...likedPostIds, post.id];
          }
        }

        setLikedPostIds(newLikedPostIds);
      } else {
        throw new Error('missing authuser');
      }
    } catch (_error) {
      Alert.alert('Oops, algo deu errado');
    }
  };

  const goToUser = async (id: string) => {
    if (authUser && authUser.id != id) {
      navigation.push('UserProfile', { userId: id });
    }
  };

  return (
    <View style={currentTheme.container}>
      {!isLoading ? (
        <>
          {posts && (
            <FlatList
              data={posts}
              keyExtractor={(post) => post.id}
              renderItem={({ item }: { item: Post }) => {
                const options: IconDropdownOption[] = [
                  {
                    label: 'Excluir Publicação',
                    iconName: 'trash-outline',
                    onSelect: async () => {
                      if (item) {
                        try {
                          const didDelete: boolean =
                            await postUsecase.deletePost(item.id);

                          if (didDelete) {
                            const newPosts: Post[] = posts.filter(
                              (post) => post.id !== item.id,
                            );

                            setPosts(newPosts);
                          }
                        } catch (_error) {
                          Alert.alert('Oops, algo deu errado');
                        }
                      }
                    },
                  },
                ];

                return (
                  <View style={styles.postContainer}>
                    <View style={styles.topPostRowContainer}>
                      <TouchableOpacity
                        style={styles.postRowContainer}
                        onPress={() => goToUser(item.user!.id)}
                      >
                        {item.user?.avatar ? (
                          <Image
                            source={{
                              uri: `data:image/jpeg;base64,${item.user!.avatar}`,
                            }}
                            style={styles.avatar}
                            resizeMode="contain"
                          />
                        ) : (
                          <View style={styles.avatarPlaceholder}>
                            <Ionicons
                              name="person-circle-outline"
                              size={45}
                              color="black"
                            ></Ionicons>
                          </View>
                        )}

                        <Text
                          style={currentTheme.textBold}
                        >{`   ${item.user?.username}`}</Text>
                      </TouchableOpacity>

                      {authUser &&
                        item.user &&
                        authUser.id === item.user.id && (
                          <View>
                            <IconDropdown options={options}></IconDropdown>
                          </View>
                        )}
                    </View>

                    <Image
                      source={{ uri: `data:image/jpeg;base64,${item.image}` }}
                      style={styles.image}
                      resizeMode="contain"
                    />

                    <View style={styles.postRowContainer}>
                      <Pressable onPress={() => handleLike(item)}>
                        <Ionicons
                          name={
                            likedPostIds.includes(item.id)
                              ? 'heart'
                              : 'heart-outline'
                          }
                          size={34}
                          color={
                            likedPostIds.includes(item.id)
                              ? 'red'
                              : currentColors.icon
                          }
                        ></Ionicons>
                      </Pressable>

                      <Text style={currentTheme.textBold}>
                        {` ${item.likeCount ?? 0}    `}
                      </Text>

                      <Pressable>
                        <Ionicons
                          name="chatbubble-outline"
                          size={34}
                          color={currentColors.icon}
                        ></Ionicons>
                      </Pressable>

                      <Text style={currentTheme.textBold}>
                        {` ${item.commentCount ?? 0}`}
                      </Text>
                    </View>

                    {item.description && (
                      <View style={styles.descriptionContainer}>
                        <Text
                          style={currentTheme.textBold}
                        >{`${item.user?.username}  `}</Text>

                        <Text style={currentTheme.text}>
                          {item.description}
                        </Text>
                      </View>
                    )}

                    <View style={styles.postRowContainer}>
                      <Text style={currentTheme.text}>
                        {`${new Date(
                          item.createdAt?.toLocaleString() ?? '',
                        ).toLocaleDateString()}  ${new Date(
                          item.createdAt?.toLocaleString() ?? '',
                        ).toLocaleTimeString(undefined, {
                          hour: '2-digit',
                          minute: '2-digit',
                        })}`}
                      </Text>
                    </View>
                  </View>
                );
              }}
              onEndReached={() => loadPosts()}
              showsVerticalScrollIndicator={false}
              contentContainerStyle={styles.flatListContainer}
              ListHeaderComponent={
                <View style={styles.listHeader}>
                  <Text style={currentTheme.logo}>y</Text>

                  <View style={styles.icon}>
                    <TouchableOpacity onPress={() => handleReload()}>
                      <Ionicons
                        name="reload"
                        size={34}
                        color={currentColors.icon}
                      ></Ionicons>
                    </TouchableOpacity>
                  </View>
                </View>
              }
            ></FlatList>
          )}

          {posts && posts.length === 0 && (
            <View style={styles.centerText}>
              <Text style={currentTheme.text}>Nenhuma publição disponível</Text>
            </View>
          )}

          {isLoadingPosts && (
            <ActivityIndicator
              size="large"
              style={styles.loadingContainer}
              color={currentColors.icon}
            />
          )}
        </>
      ) : (
        <ActivityIndicator size="large" color={currentColors.icon} />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  avatar: {
    borderRadius: 100,
    height: 45,
    width: 45,
  },
  avatarPlaceholder: {
    alignItems: 'center',
    backgroundColor: '#ccc' as string,
    borderRadius: 100,
    height: 45,
    justifyContent: 'center',
    width: 45,
  },
  centerText: {
    position: 'absolute',
    top: '50%',
  },
  descriptionContainer: {
    alignItems: 'flex-start',
    flexDirection: 'row',
    paddingLeft: 10,
    width: 365,
  },
  flatListContainer: {
    flexGrow: 1,
    paddingTop: 25,
  },
  icon: {
    marginTop: 23,
  },
  image: {
    height: 420,
    width: 420,
  },
  listHeader: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingBottom: 30,
    paddingLeft: 20,
    paddingRight: 20,
    width: 420,
  },
  loadingContainer: {
    paddingVertical: 5,
  },
  postContainer: {
    paddingBottom: 45,
  },
  postRowContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'flex-start',
    paddingBottom: 5,
    paddingLeft: 10,
    paddingTop: 3,
  },
  topPostRowContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingRight: 20,
  },
});
