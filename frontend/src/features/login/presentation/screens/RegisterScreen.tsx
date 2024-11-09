import { Ionicons } from '@expo/vector-icons';
import { CommonActions } from '@react-navigation/native';
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
import { RootStackScreenProps } from '@/src/core/navigation/types';
import { User } from '@/src/features/shared/data/models/User';

export const RegisterScreen: React.FC<RootStackScreenProps<'Register'>> = ({
  navigation,
}) => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [avatar, setAvatar] = useState<string | null>(null);
  const [avatarUri, setAvatarUri] = useState<string | null>(null);
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [isPasswordVisible, setIsPasswordVisible] = useState<boolean>(false);
  const [email, setEmail] = useState<string>('');
  const [fullName, setFullName] = useState<string>('');

  const canGoBack = navigation.canGoBack();

  const isDisabled: boolean = username.trim() === '' || password.trim() === '';

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('registerscreen must be used within an authprovider');
  }

  const { register, logout } = context;

  const handleRegister = async () => {
    setIsLoading(true);
    Keyboard.dismiss();
    try {
      if (username && password) {
        const userData: Omit<User, 'id'> = {
          username: username,
          password: password,
        };
        if (email && !isValidEmail(email)) {
          setIsLoading(false);
          Alert.alert(
            'Oops, algo deu errado',
            'Por favor, insira um email válido',
          );
          return;
        }
        if (email) {
          userData.email = email;
        }
        if (fullName) {
          userData.fullName = fullName;
        }
        if (avatar) {
          userData.avatar = avatar;
        }

        await register(userData);
        setIsLoading(false);
        navigation.dispatch(
          CommonActions.reset({
            index: 0,
            routes: [{ name: 'Tabs' }],
          }),
        );
      } else {
        setIsLoading(false);
        Alert.alert(
          'Oops, algo deu errado',
          'Nome de usuário e senha precisam estar preenchidos',
        );
      }
    } catch (error) {
      if (error instanceof Error) {
        setIsLoading(false);
        logout();
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

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={styles.container}
    >
      {canGoBack && (
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={styles.backButton}
        >
          <Ionicons name="arrow-back" size={40} color="white" />
        </TouchableOpacity>
      )}

      <Text style={styles.logo}>y</Text>

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
        <TouchableOpacity style={styles.iconButton} onPress={() => takePhoto()}>
          <Ionicons name="camera" size={32} color="white" />
        </TouchableOpacity>

        <TouchableOpacity
          style={styles.iconButton}
          onPress={() => selectImageFromLibrary()}
        >
          <Ionicons name="image" size={32} color="white" />
        </TouchableOpacity>
      </View>

      <TextInput
        style={styles.input}
        maxLength={20}
        placeholder="Nome de usuário"
        placeholderTextColor="#DDD"
        value={username}
        onChangeText={setUsername}
      />

      <View style={styles.passwordContainer}>
        <TextInput
          style={styles.inputPassword}
          maxLength={30}
          placeholder="Senha"
          placeholderTextColor="#DDD"
          secureTextEntry={!isPasswordVisible}
          value={password}
          onChangeText={setPassword}
        />
        <TouchableOpacity
          onPress={() => setIsPasswordVisible(!isPasswordVisible)}
        >
          <Ionicons
            name={isPasswordVisible ? 'eye-off' : 'eye'}
            size={26}
            color="#ddd"
            style={styles.icon}
          />
        </TouchableOpacity>
      </View>

      <TextInput
        style={styles.input}
        placeholder="Email (opcional)"
        placeholderTextColor="#DDD"
        value={email}
        onChangeText={setEmail}
      />

      <TextInput
        style={styles.input}
        maxLength={50}
        placeholder="Nome completo (opcional)"
        placeholderTextColor="#DDD"
        value={fullName}
        onChangeText={setFullName}
      />

      {!isLoading ? (
        <TouchableOpacity
          style={isDisabled ? styles.buttonDisabled : styles.button}
          onPress={() => handleRegister()}
          disabled={isDisabled}
        >
          <Text style={styles.buttonText}>Cadastrar</Text>
        </TouchableOpacity>
      ) : (
        <ActivityIndicator size="large" color="#FFFFFF" />
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
  backButton: {
    left: 20,
    position: 'absolute',
    top: 50,
    zIndex: 1,
  },
  button: {
    backgroundColor: '#8A2BE2' as string,
    borderRadius: 5,
    marginTop: 10,
    paddingHorizontal: 20,
    paddingVertical: 12,
  },
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: 20,
    width: '80%',
  },
  buttonDisabled: {
    backgroundColor: 'gray' as string,
    borderRadius: 5,
    marginTop: 10,
    paddingHorizontal: 20,
    paddingVertical: 12,
  },
  buttonText: {
    color: 'white' as string,
    fontSize: 20,
    textAlign: 'center',
  },
  container: {
    alignItems: 'center',
    backgroundColor: '#310d6b' as string,
    flex: 1,
    justifyContent: 'center',
  },
  icon: {
    marginLeft: 35,
  },
  iconButton: {
    alignItems: 'center',
    backgroundColor: '#8A2BE2' as string,
    borderRadius: 5,
    flex: 1,
    justifyContent: 'center',
    margin: 5,
    padding: 10,
  },
  input: {
    backgroundColor: '#250a4e' as string,
    borderColor: '#9b59b6' as string,
    borderRadius: 5,
    borderWidth: 1,
    color: 'white' as string,
    marginBottom: 20,
    padding: 10,
    width: '76%',
  },
  inputPassword: {
    color: 'white' as string,
    padding: 10,
    width: '76%',
  },
  logo: {
    color: 'white' as string,
    fontSize: 50,
    fontWeight: 'bold',
    marginBottom: 30,
  },
  passwordContainer: {
    alignItems: 'center',
    backgroundColor: '#250a4e' as string,
    borderColor: '#9b59b6' as string,
    borderRadius: 5,
    borderWidth: 1,
    flexDirection: 'row',
    marginBottom: 20,
    width: '76%',
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
