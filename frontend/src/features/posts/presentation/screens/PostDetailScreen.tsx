import { Ionicons } from '@expo/vector-icons';
import {
  RouteProp,
  useFocusEffect,
  useNavigation,
  useRoute,
} from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import * as ImageManipulator from 'expo-image-manipulator';
import * as ImagePicker from 'expo-image-picker';
import React, { useContext, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  Image,
  Keyboard,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

import {
  IconDropdown,
  IconDropdownOption,
} from '@/src/core/components/IconDropdown';
import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { FeedStackParamList } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { Post } from '@/src/features/shared/data/models/Post';
import { User } from '@/src/features/shared/data/models/User';

import { PostRepositoryImpl } from '../../data/repositories/PostRepositoryImpl';
import { PostUsecaseImpl } from '../../domain/usecases/PostUsecase';

export const PostDetailScreen: React.FC = () => {
  const navigation =
    useNavigation<StackNavigationProp<FeedStackParamList, 'UserProfile'>>();

  const route = useRoute<RouteProp<FeedStackParamList, 'PostDetail'>>();
  const { postId, editing } = route.params;

  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [post, setPost] = useState<Post>();
  const [isLiked, setIsLiked] = useState<boolean>(false);
  const [isEditing, setIsEditing] = useState<boolean>(editing ?? false);
  const [image, setImage] = useState<string | null>(null);
  const [imageUri, setImageUri] = useState<string | null>(null);
  const [description, setDescription] = useState<string>('');
  const postRepository = new PostRepositoryImpl();
  const postUsecase = new PostUsecaseImpl(postRepository);

  const canGoBack = navigation.canGoBack();

  const isDisabled: boolean = image === null;

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('postdetailscreen must be used within an authprovider');
  }

  const { authUser } = context;

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  useFocusEffect(
    React.useCallback(() => {
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
            if (mainPost.user && authUser.id === mainPost.user.id) {
              if (mainPost.image) {
                setImage(mainPost.image);
                setImageUri(`data:image/jpeg;base64,${mainPost.image}`);
              }
              if (mainPost.description) {
                setDescription(mainPost.description);
              }
            }
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
    }, []),
  );

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
        throw new Error('missing authuser or post');
      }
    } catch (_error) {
      Alert.alert('Oops, algo deu errado');
    }
  };

  const goToUser = async (id: string) => {
    navigation.push('UserProfile', { userId: id });
  };

  const goToLikes = async () => {
    if (post) {
      if ((post.likeCount ?? 0) > 0) {
        try {
          const likes: User[] = await postUsecase.getLikes(post.id);

          if (likes) {
            navigation.push('UserList', { users: likes, title: 'Curtidas' });
          }
        } catch (_error) {
          Alert.alert('Oops, algo deu errado');
        }
      }
    }
  };

  const goToComments = async (id: string) => {
    navigation.push('PostComments', { postId: id });
  };

  const handlePost = async () => {
    setIsLoading(true);
    Keyboard.dismiss();
    try {
      if (post) {
        if (image) {
          const editedPost: Post = {
            id: post.id,
            image: image,
            createdAt: post.createdAt,
          };
          if (authUser) {
            editedPost.user = authUser;

            if (description) {
              editedPost.description = description.trim();
            }

            const didUpdate: boolean = await postUsecase.updatePost(editedPost);

            if (didUpdate) {
              setPost(editedPost);
              setIsEditing(false);
              setIsLoading(false);
            } else {
              setIsLoading(false);
            }
          } else {
            throw new Error('missing authuser');
          }
        } else {
          throw new Error('missing post');
        }
      } else {
        setIsLoading(false);
        Alert.alert('Oops, algo deu errado', 'A publicação requer uma imagem');
      }
    } catch (_error) {
      setIsLoading(false);
      Alert.alert('Oops, algo deu errado');
    }
  };

  const selectImageFromLibrary = async () => {
    const permissionResult =
      await ImagePicker.requestMediaLibraryPermissionsAsync();
    if (permissionResult.granted === false) {
      Alert.alert(
        'Oops, algo deu errado',
        'O aplicativo precisa de permissão para acessar a galeria',
      );
      return;
    }

    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ImagePicker.MediaTypeOptions.All,
      quality: 1,
    });

    if (!result.canceled && result.assets && result.assets.length > 0) {
      try {
        await cropImage(result.assets[0].uri);
      } catch (_error) {
        Alert.alert(
          'Oops, algo deu errado',
          'Por favor, tente selecionar uma imagem novamente',
        );
      }
    }
  };

  const takePhoto = async () => {
    const permissionResult = await ImagePicker.requestCameraPermissionsAsync();
    if (permissionResult.granted === false) {
      Alert.alert(
        'Oops, algo deu errado',
        'O aplicativo precisa de permissão para acessar a câmera',
      );
      return;
    }

    const result = await ImagePicker.launchCameraAsync({
      quality: 1,
    });

    if (!result.canceled && result.assets && result.assets.length > 0) {
      try {
        await cropImage(result.assets[0].uri);
      } catch (_error) {
        Alert.alert(
          'Oops, algo deu errado',
          'Por favor, tente tirar uma foto novamente',
        );
      }
    }
  };

  const cropImage = async (uri: string) => {
    let cropWidth = 1080;
    let cropHeight = 1080;

    const { width, height } = await ImageManipulator.manipulateAsync(uri);

    if (width < cropWidth || height < cropHeight) {
      cropWidth = 720;
      cropHeight = 720;
    }

    const cropX = (width - cropWidth) / 2;
    const cropY = (height - cropHeight) / 2;

    const cropData = {
      crop: {
        originX: cropX,
        originY: cropY,
        width: cropWidth,
        height: cropHeight,
      },
    };

    const result = await ImageManipulator.manipulateAsync(uri, [cropData], {
      compress: 0.8,
      format: ImageManipulator.SaveFormat.JPEG,
      base64: true,
    });

    if (result.base64) {
      setImageUri(result.uri);
      setImage(result.base64);
    } else {
      throw new Error('error converting image to base64');
    }
  };

  const clearImage = () => {
    setImage(null);
    setImageUri(null);
  };

  const options: IconDropdownOption[] = [
    {
      label: 'Editar Publicação',
      iconName: 'pencil',
      onSelect: async () => {
        setIsEditing(true);
      },
    },
    {
      label: 'Excluir Publicação',
      iconName: 'trash-outline',
      onSelect: async () => {
        if (post) {
          Alert.alert(
            'Confirmar exclusão',
            'Você tem certeza absoluta de que deseja excluir sua publicação?',
            [
              {
                text: 'Cancelar',
                style: 'cancel',
              },
              {
                text: 'Excluir',
                style: 'destructive',
                onPress: async () => {
                  try {
                    const didDelete: boolean = await postUsecase.deletePost(
                      post.id,
                    );
                    if (didDelete && canGoBack) {
                      navigation.goBack();
                    }
                  } catch (_error) {
                    Alert.alert('Oops, algo deu errado');
                  }
                },
              },
            ],
            { cancelable: true },
          );
        }
      },
    },
  ];

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={currentTheme.container}
    >
      {!isLoading && canGoBack && (
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={currentTheme.backButton}
        >
          <Ionicons name="arrow-back" size={40} color={currentColors.icon} />
        </TouchableOpacity>
      )}

      {!isLoading ? (
        post && !isEditing ? (
          <>
            {authUser && post.user && authUser.id === post.user.id && (
              <View style={currentTheme.topRow}>
                <IconDropdown options={options}></IconDropdown>
              </View>
            )}

            <ScrollView
              contentContainerStyle={styles.containerScroll}
              showsVerticalScrollIndicator={false}
            >
              <TouchableOpacity
                onPress={() => goToUser(post.user!.id)}
                disabled={authUser ? post.user!.id === authUser!.id : true}
              >
                <View style={styles.postRowContainer}>
                  {post.user?.avatar ? (
                    <Image
                      source={{
                        uri: `data:image/jpeg;base64,${post.user!.avatar}`,
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
                </View>
              </TouchableOpacity>

              <Image
                source={{ uri: `data:image/jpeg;base64,${post.image}` }}
                style={styles.image}
                resizeMode="contain"
              />

              <View style={styles.postRowContainer}>
                <Pressable
                  style={styles.row}
                  onPress={() => handleLike()}
                  onLongPress={() => goToLikes()}
                >
                  <Ionicons
                    name={isLiked ? 'heart' : 'heart-outline'}
                    size={34}
                    color={isLiked ? 'red' : currentColors.icon}
                  ></Ionicons>

                  <Text style={currentTheme.textBold}>
                    {` ${post.likeCount ?? 0}    `}
                  </Text>
                </Pressable>

                <TouchableOpacity
                  style={styles.row}
                  onPress={() => goToComments(post.id)}
                >
                  <Ionicons
                    name="chatbubble-outline"
                    size={34}
                    color={currentColors.icon}
                  ></Ionicons>

                  <Text style={currentTheme.textBold}>
                    {` ${post.commentCount ?? 0}`}
                  </Text>
                </TouchableOpacity>
              </View>

              {post.description && (
                <View style={styles.descriptionContainer}>
                  <Text style={currentTheme.text}>
                    <Text
                      style={currentTheme.textBold}
                    >{`${post.user?.username} `}</Text>
                    {post.description}
                  </Text>
                </View>
              )}

              <View style={styles.postRowContainer}>
                <Text style={currentTheme.text}>
                  {`${new Date(
                    post.createdAt?.toLocaleString() ?? '',
                  ).toLocaleDateString()}  ${new Date(
                    post.createdAt?.toLocaleString() ?? '',
                  ).toLocaleTimeString(undefined, {
                    hour: '2-digit',
                    minute: '2-digit',
                  })}`}
                </Text>
              </View>
            </ScrollView>
          </>
        ) : (
          <>
            {isLoading && (
              <View style={currentTheme.loadingOverlay}>
                <ActivityIndicator size="large" color="white" />
              </View>
            )}

            {imageUri ? (
              <View style={styles.imageContainer}>
                <Image
                  source={{ uri: imageUri }}
                  style={styles.editingImage}
                  resizeMode="contain"
                />

                <TouchableOpacity
                  style={styles.trashIconContainer}
                  onPress={() => {
                    clearImage();
                  }}
                >
                  <Ionicons name="trash" size={24} color="red" />
                </TouchableOpacity>
              </View>
            ) : (
              <View style={styles.buttonContainer}>
                <TouchableOpacity
                  style={currentTheme.filledIconButton}
                  onPress={() => takePhoto()}
                >
                  <Ionicons name="camera" size={32} color="white" />
                </TouchableOpacity>

                <TouchableOpacity
                  style={currentTheme.filledIconButton}
                  onPress={() => selectImageFromLibrary()}
                >
                  <Ionicons name="image" size={32} color="white" />
                </TouchableOpacity>
              </View>
            )}

            <TextInput
              style={currentTheme.largeInput}
              multiline
              maxLength={190}
              placeholder="Descrição (opcional)"
              placeholderTextColor={currentColors.placeholderText}
              value={description}
              onChangeText={setDescription}
              textAlignVertical="top"
            />

            {!isLoading && (
              <View style={styles.bottomIconButtons}>
                <TouchableOpacity
                  style={styles.cancelButton}
                  onPress={() => setIsEditing(false)}
                >
                  <Text style={currentTheme.buttonText}>Cancelar</Text>
                </TouchableOpacity>

                <TouchableOpacity
                  style={
                    isDisabled
                      ? currentTheme.buttonDisabled
                      : currentTheme.button
                  }
                  onPress={() => handlePost()}
                  disabled={isDisabled}
                >
                  <Text style={currentTheme.buttonText}>Salvar</Text>
                </TouchableOpacity>
              </View>
            )}
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
  bottomIconButtons: {
    bottom: 16,
    flexDirection: 'row',
    padding: 10,
    position: 'absolute',
    right: 16,
  },
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginVertical: 208,
    width: '76%',
  },
  cancelButton: {
    backgroundColor: 'red' as string,
    borderRadius: 5,
    marginBottom: 20,
    marginRight: 20,
    marginTop: 10,
    paddingHorizontal: 20,
    paddingVertical: 12,
  },
  containerScroll: {
    flexGrow: 1,
    justifyContent: 'flex-start',
    marginTop: '30%',
  },
  descriptionContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    paddingLeft: 10,
    width: 405,
  },
  editingImage: {
    height: '100%',
    width: '100%',
  },
  image: {
    height: 420,
    width: 420,
  },
  imageContainer: {
    height: 400,
    marginVertical: 40,
    position: 'relative',
    width: '90%',
  },
  postRowContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'flex-start',
    paddingBottom: 5,
    paddingLeft: 10,
    paddingTop: 3,
  },
  row: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'flex-start',
  },
  trashIconContainer: {
    backgroundColor: 'white' as string,
    borderRadius: 30,
    padding: 7,
    position: 'absolute',
    right: 2,
    top: 17,
  },
});
