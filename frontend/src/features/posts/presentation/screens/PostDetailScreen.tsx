import { Ionicons } from '@expo/vector-icons';
import { RouteProp, useNavigation, useRoute } from '@react-navigation/native';
import React, { useContext, useEffect, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  Image,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';

import {
  IconDropdown,
  IconDropdownOption,
} from '@/src/core/components/iconDropdown';
import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { FeedStackParamList } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { Post } from '@/src/features/shared/data/models/Post';

import { PostRepositoryImpl } from '../../data/repositories/PostRepositoryImpl';
import { PostUsecaseImpl } from '../../domain/usecases/PostUsecase';

export const PostDetailScreen: React.FC = () => {
  const navigation = useNavigation();

  const route = useRoute<RouteProp<FeedStackParamList, 'PostDetail'>>();
  const { postId } = route.params;

  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [post, setPost] = useState<Post>();
  const [isLiked, setIsLiked] = useState<boolean>(false);
  const postRepository = new PostRepositoryImpl();
  const postUsecase = new PostUsecaseImpl(postRepository);

  const canGoBack = navigation.canGoBack();

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('postdetailscreen must be used within an authprovider');
  }

  const { authUser } = context;

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  useEffect(() => {
    const loadPost = async () => {
      setIsLoading(true);
      try {
        const mainPost: Post = await postUsecase.getPost(postId);
        if (mainPost && authUser) {
          const didLike: boolean = await postUsecase.checkIfUserLikedPost(
            authUser.id,
            mainPost.id,
          );
          if (didLike) {
            setIsLiked(true);
          }

          setIsLoading(false);
          setPost(mainPost);
        } else {
          setIsLoading(false);
          throw new Error('missing post or authuser');
        }
      } catch (_error) {
        setIsLoading(false);
        Alert.alert('Oops, algo deu errado');
      }
    };

    if (!post) {
      loadPost();
    }
  }, [post]);

  const handleLike = async () => {
    try {
      if (authUser && post) {
        if (isLiked) {
          const didUnlike: boolean = await postUsecase.unlikePost(
            authUser.id,
            post.id,
          );
          if (didUnlike) {
            post.likeCount = (post.likeCount ?? 0) - 1;
            setIsLiked(false);
          }
        } else {
          const didLike: boolean = await postUsecase.likePost(
            authUser.id,
            post.id,
          );
          if (didLike) {
            post.likeCount = (post.likeCount ?? 0) + 1;
            setIsLiked(true);
          }
        }
      } else {
        throw new Error('no authuser or post');
      }
    } catch (_error) {
      Alert.alert('Oops, algo deu errado');
    }
  };

  const options: IconDropdownOption[] = [
    {
      label: 'Excluir Publicação',
      iconName: 'trash',
      onSelect: async () => {
        if (post) {
          try {
            await postUsecase.deletePost(post.id);
            if (canGoBack) {
              navigation.goBack();
            }
          } catch (_error) {
            Alert.alert('Oops, algo deu errado');
          }
        }
      },
    },
  ];

  return (
    <ScrollView contentContainerStyle={currentTheme.containerLeftAligned}>
      {canGoBack && (
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={currentTheme.backButton}
        >
          <Ionicons name="arrow-back" size={34} color={currentColors.icon} />
        </TouchableOpacity>
      )}
      {!isLoading ? (
        post && (
          <>
            {authUser && post.user && authUser.id === post.user.id && (
              <View style={currentTheme.topRow}>
                <IconDropdown options={options}></IconDropdown>
              </View>
            )}

            <View style={styles.postRowContainer}>
              <TouchableOpacity>
                {post.user?.avatar ? (
                  <Image
                    source={{
                      uri: `data:image/jpeg;base64,${post.user?.avatar}`,
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
                >{`   ${post.user?.username}`}</Text>
              </TouchableOpacity>
            </View>

            <Image
              source={{ uri: `data:image/jpeg;base64,${post.image}` }}
              style={styles.image}
              resizeMode="contain"
            />

            <View style={styles.postRowContainer}>
              <Pressable onPress={handleLike}>
                <Ionicons
                  name={isLiked ? 'heart' : 'heart-outline'}
                  size={34}
                  color={isLiked ? 'red' : currentColors.icon}
                ></Ionicons>
              </Pressable>

              <Text style={currentTheme.textBold}>
                {` ${post.likeCount ?? 0}     `}
              </Text>

              <Pressable>
                <Ionicons
                  name="chatbubble-outline"
                  size={34}
                  color={currentColors.icon}
                ></Ionicons>
              </Pressable>

              <Text style={currentTheme.textBold}>
                {` ${post.commentCount ?? 0}`}
              </Text>
            </View>

            {post.description && (
              <View style={styles.descriptionContainer}>
                <Text
                  style={currentTheme.textBold}
                >{`${post.user?.username}   `}</Text>

                <Text style={currentTheme.text}>{post.description}</Text>
              </View>
            )}
          </>
        )
      ) : (
        <ActivityIndicator size="large" color={currentColors.primary} />
      )}
    </ScrollView>
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
  descriptionContainer: {
    flexDirection: 'row',
    paddingLeft: 10,
    paddingTop: 5,
  },
  image: {
    height: '50%',
    width: '100%',
  },
  postRowContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    paddingLeft: 10,
    paddingTop: 3,
  },
});
