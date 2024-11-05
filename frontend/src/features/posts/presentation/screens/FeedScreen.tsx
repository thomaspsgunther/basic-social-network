import React, { useContext, useEffect, useState } from 'react';
import { ActivityIndicator, Alert, Text, View } from 'react-native';

import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { FeedStackScreenProps } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { Post } from '@/src/features/shared/data/models/Post';

import { PostRepositoryImpl } from '../../data/repositories/PostRepositoryImpl';
import { PostUsecaseImpl } from '../../domain/usecases/PostUsecase';

export const FeedScreen: React.FC<FeedStackScreenProps<'Feed'>> = () =>
  //   {
  //   navigation,
  //   },
  {
    const [isLoading, setIsLoading] = useState<boolean>(false);
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
      if (!posts) {
        loadPosts();
        checkLikes(posts);
      }
    }, [posts]);

    const loadPosts = async () => {
      setIsLoading(true);
      try {
        const initialPosts = await postUsecase.listPosts(3);

        setIsLoading(false);
        setPosts(initialPosts);
      } catch (_error) {
        setIsLoading(false);
        Alert.alert('Oops, algo deu errado');
      }
    };

    const checkLikes = async (posts: Post[]) => {
      try {
        if (authUser) {
          const newLikedPostIds: string[] = likedPostIds;
          for (const post of posts) {
            const didLike: boolean = await postUsecase.checkIfUserLikedPost(
              authUser.id,
              post.id,
            );
            if (didLike) {
              if (!newLikedPostIds.includes(post.id)) {
                newLikedPostIds.push(post.id);
              }
            }
          }
          if (newLikedPostIds !== likedPostIds) {
            setLikedPostIds(newLikedPostIds);
          }
        } else {
          throw new Error('missing authuser');
        }
      } catch (_error) {
        Alert.alert('Oops, algo deu errado');
      }
    };

    return (
      <View style={currentTheme.container}>
        {!isLoading ? (
          posts && (
            <>
              <Text style={currentTheme.text}>Feed</Text>
            </>
          )
        ) : (
          <ActivityIndicator size="large" color={currentColors.primary} />
        )}
      </View>
    );
  };
