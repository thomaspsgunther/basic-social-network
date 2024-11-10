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

import {
  IconDropdown,
  IconDropdownOption,
} from '@/src/core/components/IconDropdown';
import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { CurrentUserProfileStackScreenProps } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { User } from '@/src/features/shared/data/models/User';

import { UserRepositoryImpl } from '../../data/repositories/UserRepositoryImpl';
import { UserUsecaseImpl } from '../../domain/usecases/UserUsecase';

export const UserEditScreen: React.FC<
  CurrentUserProfileStackScreenProps<'UserEdit'>
> = ({ navigation }) => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('UserEditscreen must be used within an authprovider');
  }

  const { authUser, setAuthUser, logoutAndLeave } = context;

  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [avatar, setAvatar] = useState<string | null>(
    authUser ? (authUser.avatar ?? null) : null,
  );
  const [avatarUri, setAvatarUri] = useState<string | null>(
    authUser
      ? authUser.avatar
        ? `data:image/jpeg;base64,${authUser.avatar}`
        : null
      : null,
  );
  const [username, setUsername] = useState<string>(
    authUser ? (authUser.username ?? '') : '',
  );
  const [password, setPassword] = useState<string>('');
  const [isPasswordVisible, setIsPasswordVisible] = useState<boolean>(false);
  const [email, setEmail] = useState<string>(
    authUser ? (authUser.email ?? '') : '',
  );
  const [fullName, setFullName] = useState<string>(
    authUser ? (authUser.fullName ?? '') : '',
  );
  const [description, setDescription] = useState<string>(
    authUser ? (authUser.description ?? '') : '',
  );
  const userRepository = new UserRepositoryImpl();
  const userUsecase = new UserUsecaseImpl(userRepository);

  const canGoBack = navigation.canGoBack();

  const isDisabled: boolean = username.trim() === '';

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  const handleEdit = async () => {
    setIsLoading(true);
    Keyboard.dismiss();
    try {
      if (authUser) {
        if (username) {
          if (!hasSpecialCharacters(username)) {
            const userData: User = {
              id: authUser.id,
              username: username,
            };
            if (password) {
              userData.password = password;
            }
            if (email && !isValidEmail(email)) {
              setIsLoading(false);
              Alert.alert(
                'Oops, algo deu errado',
                'Por favor, insira um email válido',
              );
              return;
            } else if (email && isValidEmail(email)) {
              userData.email = email.trim();
            }
            if (fullName) {
              userData.fullName = fullName.trim();
            }
            if (description) {
              userData.description = description.trim();
            }
            if (avatar) {
              userData.avatar = avatar;
            }

            const didUpdate = await userUsecase.updateUser(userData);

            if (didUpdate) {
              userData.password = undefined;

              setIsLoading(false);
              setAuthUser(userData);
              if (canGoBack) {
                navigation.goBack();
              }
            } else {
              setIsLoading(false);
            }
          } else {
            setIsLoading(false);
            Alert.alert(
              'Oops, algo deu errado',
              'Nome de usuário não pode conter caracteres especiais',
            );
          }
        } else {
          setIsLoading(false);
          Alert.alert(
            'Oops, algo deu errado',
            'Nome de usuário e senha precisam estar preenchidos',
          );
        }
      } else {
        throw new Error('missing authuser');
      }
    } catch (error) {
      if (error instanceof Error) {
        setIsLoading(false);
        if (error.message.trim() === 'user already exists') {
          Alert.alert(
            'Oops, algo deu errado',
            'Esse nome de usuário já é utilizado',
          );
        } else {
          Alert.alert('Oops, algo deu errado', 'Por favor, tente novamente');
        }
      }
    }
  };

  const handleUsernameChange = (input: string) => {
    const trimmedUsername: string = input.trim();
    const noSpacesInput: string = trimmedUsername.replace(/\s+/g, '');
    const lowercaseUsername: string = noSpacesInput.toLowerCase();
    setUsername(lowercaseUsername);
  };

  const handlePasswordChange = (input: string) => {
    const trimmedPassword: string = input.trim();
    const noSpacesPassword: string = trimmedPassword.replace(/\s+/g, '');
    setPassword(noSpacesPassword);
  };

  const hasSpecialCharacters = (str: string) => /[^a-zA-Z0-9\s]/.test(str);

  const isValidEmail = (email: string) => {
    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return regex.test(email);
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
      setAvatarUri(result.uri);
      setAvatar(result.base64);
    } else {
      throw new Error('error converting image to base64');
    }
  };

  const clearAvatar = () => {
    setAvatar(null);
    setAvatarUri(null);
  };

  const options: IconDropdownOption[] = [
    {
      label: 'Excluir Conta',
      iconName: 'trash-outline',
      onSelect: async () => {
        if (authUser) {
          Alert.alert(
            'Confirmar exclusão',
            'Você tem certeza absoluta de que deseja excluir sua conta?',
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
                    const didDelete: boolean = await userUsecase.deleteUser(
                      authUser.id,
                    );

                    if (didDelete) {
                      logoutAndLeave();
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
      {canGoBack && (
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={currentTheme.backButton}
        >
          <Ionicons name="arrow-back" size={40} color={currentColors.icon} />
        </TouchableOpacity>
      )}

      <View style={currentTheme.topRow}>
        <IconDropdown options={options}></IconDropdown>
      </View>

      <TouchableOpacity
        onPress={() => takePhoto()}
        style={styles.avatarContainer}
      >
        {avatarUri ? (
          <Image
            source={{ uri: avatarUri }}
            style={styles.avatar}
            resizeMode="contain"
          />
        ) : (
          <View style={styles.avatarPlaceholder}>
            <Text style={styles.avatarPlaceholderText}>Tire uma foto!</Text>
          </View>
        )}
        {avatarUri && (
          <TouchableOpacity
            onPress={() => clearAvatar()}
            style={styles.trashIconContainer}
          >
            <Ionicons name="trash" size={24} color="red" />
          </TouchableOpacity>
        )}
      </TouchableOpacity>

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

      <TextInput
        style={currentTheme.input}
        maxLength={20}
        placeholder="Nome de usuário"
        placeholderTextColor={currentColors.placeholderText}
        value={username}
        onChangeText={handleUsernameChange}
      />

      <View style={currentTheme.inputIconContainer}>
        <TextInput
          style={currentTheme.inputIcon}
          maxLength={30}
          placeholder="Senha"
          placeholderTextColor={currentColors.placeholderText}
          secureTextEntry={!isPasswordVisible}
          value={password}
          onChangeText={handlePasswordChange}
        />
        <TouchableOpacity
          onPress={() => setIsPasswordVisible(!isPasswordVisible)}
        >
          <Ionicons
            name={isPasswordVisible ? 'eye-off' : 'eye'}
            size={26}
            color={currentColors.icon}
            style={styles.icon}
          />
        </TouchableOpacity>
      </View>

      <TextInput
        style={currentTheme.input}
        placeholder="Email (opcional)"
        placeholderTextColor={currentColors.placeholderText}
        value={email}
        onChangeText={setEmail}
      />

      <TextInput
        style={currentTheme.input}
        maxLength={50}
        placeholder="Nome completo (opcional)"
        placeholderTextColor={currentColors.placeholderText}
        value={fullName}
        onChangeText={setFullName}
      />

      <TextInput
        style={currentTheme.largeInput}
        multiline
        maxLength={215}
        placeholder="Descrição (opcional)"
        placeholderTextColor={currentColors.placeholderText}
        value={description}
        onChangeText={setDescription}
        textAlignVertical="top"
      />

      {!isLoading ? (
        <TouchableOpacity
          style={isDisabled ? currentTheme.buttonDisabled : currentTheme.button}
          onPress={() => handleEdit()}
          disabled={isDisabled}
        >
          <Text style={currentTheme.buttonText}>Salvar</Text>
        </TouchableOpacity>
      ) : (
        <ActivityIndicator size="large" color={currentColors.icon} />
      )}
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  avatar: {
    borderRadius: 100,
    height: 140,
    width: 140,
  },
  avatarContainer: {
    marginBottom: 20,
    marginTop: 60,
    position: 'relative',
  },
  avatarPlaceholder: {
    alignItems: 'center',
    backgroundColor: '#ccc' as string,
    borderRadius: 100,
    height: 140,
    justifyContent: 'center',
    width: 140,
  },
  avatarPlaceholderText: {
    color: '#777' as string,
  },
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: 20,
    width: '80%',
  },
  icon: {
    marginLeft: 35,
  },
  trashIconContainer: {
    backgroundColor: 'white' as string,
    borderRadius: 15,
    padding: 5,
    position: 'absolute',
    right: 5,
    top: 5,
  },
});
