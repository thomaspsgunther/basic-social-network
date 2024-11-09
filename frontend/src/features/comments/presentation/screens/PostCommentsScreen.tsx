import { Ionicons } from '@expo/vector-icons';
import { RouteProp, useNavigation, useRoute } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { useContext, useEffect, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  Image,
  Keyboard,
  KeyboardAvoidingView,
  Platform,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';
import { FlatList, TextInput } from 'react-native-gesture-handler';

import {
  IconDropdown,
  IconDropdownOption,
} from '@/src/core/components/IconDropdown';
import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { FeedStackParamList } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';

import { Comment } from '../../data/models/Comment';
import { CommentRepositoryImpl } from '../../data/repositories/CommentRepositoryImpl';
import { CommentUsecaseImpl } from '../../domain/usecases/CommentUsecase';

export const PostCommentsScreen: React.FC = () => {
  const navigation =
    useNavigation<StackNavigationProp<FeedStackParamList, 'PostComments'>>();

  const route = useRoute<RouteProp<FeedStackParamList, 'PostComments'>>();
  const { postId } = route.params;

  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isLoadingComment, setIsLoadingComment] = useState<boolean>(false);
  const [message, setMessage] = useState<string>('');
  const [comments, setComments] = useState<Comment[]>();
  const commentRepository = new CommentRepositoryImpl();
  const commentUsecase = new CommentUsecaseImpl(commentRepository);

  const canGoBack = navigation.canGoBack();

  const isDisabled: boolean = message.trim() === '';

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('postcommentsscreen must be used within an authprovider');
  }

  const { authUser } = context;

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  useEffect(() => {
    const loadComments = async () => {
      setIsLoading(true);
      try {
        const postComments: Comment[] =
          await commentUsecase.getCommentsFromPost(postId);
        if (postComments) {
          setIsLoading(false);
          setComments(postComments);
        } else {
          setIsLoading(false);
          throw new Error('missing comments');
        }
      } catch (_error) {
        setIsLoading(false);
        Alert.alert('Oops, algo deu errado');
      }
    };

    if (!comments) {
      loadComments();
    }
  }, [comments]);

  const handleComment = async () => {
    setIsLoadingComment(true);
    Keyboard.dismiss();
    try {
      if (message) {
        const comment: Omit<Comment, 'id'> = {
          postId: postId,
          message: message,
        };
        if (authUser) {
          comment.user = authUser;

          const newComment = await commentUsecase.createComment(comment);

          if (newComment) {
            newComment.user = comment.user;
            newComment.message = comment.message;

            if (comments) {
              const newComments: Comment[] = [...comments];

              newComments.unshift(newComment);
              setComments(newComments);
            }

            setMessage('');
            setIsLoadingComment(false);
          }
        } else {
          throw new Error('missing authuser');
        }
      } else {
        setIsLoadingComment(false);
        Alert.alert(
          'Oops, algo deu errado',
          'O comentário requer uma mensagem',
        );
      }
    } catch (_error) {
      setIsLoadingComment(false);
      Alert.alert('Oops, algo deu errado');
    }
  };

  const goToUser = async (id: string) => {
    if (authUser && authUser.id != id) {
      navigation.push('UserProfile', { userId: id });
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={currentTheme.container}
    >
      {!isLoading && !comments && canGoBack && (
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={currentTheme.backButton}
        >
          <Ionicons name="arrow-back" size={40} color={currentColors.icon} />
        </TouchableOpacity>
      )}

      {!isLoading ? (
        comments && (
          <>
            <View style={styles.listHeader}>
              {canGoBack && (
                <TouchableOpacity onPress={() => navigation.goBack()}>
                  <Ionicons
                    name="arrow-back"
                    size={40}
                    color={currentColors.icon}
                  />
                </TouchableOpacity>
              )}
            </View>

            <FlatList
              data={comments}
              keyExtractor={(comment) => comment.id}
              renderItem={({ item }: { item: Comment }) => {
                const options: IconDropdownOption[] = [
                  {
                    label: 'Excluir Comentário',
                    iconName: 'trash-outline',
                    onSelect: async () => {
                      if (item) {
                        try {
                          const didDelete: boolean =
                            await commentUsecase.deleteComment(item.id);

                          if (didDelete) {
                            const newComments: Comment[] = comments.filter(
                              (comment) => comment.id !== item.id,
                            );

                            setComments(newComments);
                          }
                        } catch (_error) {
                          Alert.alert('Oops, algo deu errado');
                        }
                      }
                    },
                  },
                ];

                return (
                  <View style={styles.commentRowContainer}>
                    <View style={styles.commentUserRowContainer}>
                      <TouchableOpacity
                        style={styles.commentUserRowContainer}
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
                      </TouchableOpacity>

                      <View style={styles.messageContainer}>
                        <Text style={currentTheme.text}>
                          <Text
                            style={currentTheme.textBold}
                          >{`${item.user?.username} `}</Text>
                          {item.message}
                        </Text>
                      </View>
                    </View>

                    {authUser && item.user && authUser.id === item.user.id && (
                      <View style={styles.iconDropdown}>
                        <IconDropdown options={options}></IconDropdown>
                      </View>
                    )}
                  </View>
                );
              }}
              contentContainerStyle={styles.flatListContainer}
            ></FlatList>

            <View style={styles.commentInputContainer}>
              <TextInput
                style={currentTheme.input}
                multiline
                maxLength={200}
                placeholder="Escreva um comentário"
                placeholderTextColor={currentColors.placeholderText}
                value={message}
                onChangeText={setMessage}
              />

              <View style={styles.icon}>
                {!isLoadingComment ? (
                  <TouchableOpacity
                    onPress={() => handleComment()}
                    disabled={isDisabled}
                  >
                    <Ionicons
                      name="send"
                      size={45}
                      color={
                        isDisabled ? currentColors.disabled : currentColors.icon
                      }
                    />
                  </TouchableOpacity>
                ) : (
                  <ActivityIndicator size="large" color={currentColors.icon} />
                )}
              </View>
            </View>
          </>
        )
      ) : (
        <ActivityIndicator size="large" color={currentColors.icon} />
      )}
    </KeyboardAvoidingView>
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
  commentInputContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  commentRowContainer: {
    alignItems: 'flex-start',
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingBottom: 15,
    paddingRight: 20,
  },
  commentUserRowContainer: {
    alignItems: 'flex-start',
    flexDirection: 'row',
    justifyContent: 'flex-start',
    paddingBottom: 5,
    paddingLeft: 10,
    paddingTop: 3,
  },
  flatListContainer: {
    flexGrow: 1,
    paddingTop: 20,
  },
  icon: {
    marginBottom: 19,
    marginLeft: 18,
  },
  iconDropdown: {
    marginTop: 9,
  },
  listHeader: {
    alignItems: 'center',
    flexDirection: 'row',
    paddingBottom: 10,
    paddingLeft: 20,
    paddingRight: 20,
    paddingTop: 50,
    width: 420,
  },
  messageContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginLeft: 10,
    marginTop: 14,
    width: 288,
  },
});
