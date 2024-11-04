import { Ionicons } from '@expo/vector-icons';
import { RouteProp, useNavigation, useRoute } from '@react-navigation/native';
import React, { useContext, useEffect, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  Image,
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
        if (mainPost) {
          setIsLoading(false);
          setPost(mainPost);
        } else {
          setIsLoading(false);
          throw new Error('missing post');
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

            <View style={styles.postUserContainer}>
              {post.user?.avatar && (
                <Image
                  source={{
                    uri: `data:image/jpeg;base64,${post.user?.avatar}`,
                  }}
                ></Image>
              )}
            </View>

            <Image
              source={{ uri: `data:image/jpeg;base64,${post.image}` }}
              style={styles.image}
              resizeMode="contain"
            />

            <View style={styles.descriptionContainer}>
              <Text
                style={currentTheme.textBold}
              >{`${post.user?.username}   `}</Text>

              <Text style={currentTheme.text}>{post.description}</Text>
            </View>
          </>
        )
      ) : (
        <ActivityIndicator size="large" color={currentColors.primary} />
      )}
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  descriptionContainer: {
    flexDirection: 'row',
    paddingLeft: 10,
    paddingTop: 5,
  },
  image: {
    height: '50%',
    width: '100%',
  },
  postUserContainer: {
    flexDirection: 'row',
    paddingLeft: 10,
    paddingTop: 5,
  },
});
