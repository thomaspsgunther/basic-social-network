import { Ionicons } from '@expo/vector-icons';
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
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { CreatePostStackScreenProps } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { Post } from '@/src/features/shared/data/models/Post';

import { PostRepositoryImpl } from '../../data/repositories/PostRepositoryImpl';
import { PostUsecaseImpl } from '../../domain/usecases/PostUsecase';

export const CreatePostScreen: React.FC<
  CreatePostStackScreenProps<'CreatePost'>
> = ({ navigation }) => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [image, setImage] = useState<string | null>(null);
  const [imageUri, setImageUri] = useState<string | null>(null);
  const [description, setDescription] = useState<string>('');
  const postRepository = new PostRepositoryImpl();
  const postUsecase = new PostUsecaseImpl(postRepository);

  const isDisabled: boolean = image === null;

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('createpostscreen must be used within an authprovider');
  }

  const { authUser } = context;

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  const handlePost = async () => {
    setIsLoading(true);
    Keyboard.dismiss();
    try {
      if (image) {
        const post: Omit<Post, 'id'> = {
          image: image,
        };
        if (authUser) {
          post.user = authUser;

          if (description) {
            post.description = description.trim();
          }

          const newPost: Post = await postUsecase.createPost(post);

          if (newPost) {
            clearImage();
            setDescription('');
            setIsLoading(false);
            navigation.push('PostDetail', { postId: newPost.id });
          } else {
            setIsLoading(false);
          }
        } else {
          throw new Error('missing authuser');
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

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={currentTheme.container}
    >
      {isLoading && (
        <View style={currentTheme.loadingOverlay}>
          <ActivityIndicator size="large" color="white" />
        </View>
      )}

      <Text style={currentTheme.titleTopLeftAligned}>Nova Publicação</Text>

      {imageUri ? (
        <View style={styles.imageContainer}>
          <Image
            source={{ uri: imageUri }}
            style={styles.image}
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
        <TouchableOpacity
          style={styles.bottomIconButton}
          onPress={() => handlePost()}
          disabled={isDisabled}
        >
          <Ionicons
            name="send"
            size={45}
            color={isDisabled ? currentColors.disabled : currentColors.icon}
          />
        </TouchableOpacity>
      )}
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  bottomIconButton: {
    bottom: 16,
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
  image: {
    height: '100%',
    width: '100%',
  },
  imageContainer: {
    height: 400,
    marginVertical: 40,
    position: 'relative',
    width: '90%',
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
