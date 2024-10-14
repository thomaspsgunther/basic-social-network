import { Ionicons } from '@expo/vector-icons';
import * as ImagePicker from 'expo-image-picker';
import React, { useContext, useState } from 'react';
import {
  Alert,
  Image,
  Keyboard,
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
  const [loading, setLoading] = useState<boolean>(false);
  const [avatar, setAvatar] = useState<string | null>(null);
  const [avatarUri, setAvatarUri] = useState<string | null>(null);
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [isPasswordVisible, setIsPasswordVisible] = useState<boolean>(false);
  const [email, setEmail] = useState<string>('');
  const [fullName, setFullName] = useState<string>('');

  const isDisabled = username.trim() === '' || password.trim() === '';

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('registerscreen must be used within an authprovider');
  }

  const { register, logout } = context;

  const handleRegister = async () => {
    setLoading(true);
    Keyboard.dismiss();
    try {
      if (username && password) {
        const userData: Omit<User, 'id'> = {
          username: username,
          password: password,
        };
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
        navigation.navigate('Tabs');
      } else {
        setLoading(false);
        Alert.alert(
          'Oops, algo deu errado',
          'Nome de usuário e senha precisam estar preenchidos',
        );
      }
    } catch (error) {
      if (error instanceof Error) {
        setLoading(false);
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
      base64: true,
      quality: 1,
    });

    if (!result.canceled) {
      const { base64, type } = result.assets[0];
      setAvatar(base64 || null);
      setAvatarUri(`data:${type};base64,${base64}`);
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
      base64: true,
      quality: 1,
    });

    if (!result.canceled) {
      const { base64, type } = result.assets[0];
      setAvatar(base64 || null);
      setAvatarUri(`data:${type};base64,${base64}`);
    }
  };

  const clearAvatar = () => {
    setAvatar(null);
    setAvatarUri(null);
  };

  return (
    <View style={styles.container}>
      <TouchableOpacity
        onPress={() => navigation.goBack()}
        style={styles.backButton}
      >
        <Ionicons name="arrow-back" size={34} color="#fff" />
      </TouchableOpacity>

      <Text style={styles.logo}>y</Text>

      <TouchableOpacity onPress={takePhoto} style={styles.avatarContainer}>
        {avatarUri ? (
          <Image source={{ uri: avatarUri }} style={styles.avatar} />
        ) : (
          <View style={styles.avatarPlaceholder}>
            <Text style={styles.avatarPlaceholderText}>Tire uma foto!</Text>
          </View>
        )}
        {avatarUri && (
          <TouchableOpacity
            onPress={clearAvatar}
            style={styles.trashIconContainer}
          >
            <Ionicons name="trash" size={24} color="red" />
          </TouchableOpacity>
        )}
      </TouchableOpacity>

      <View style={styles.buttonContainer}>
        <TouchableOpacity style={styles.iconButton} onPress={takePhoto}>
          <Ionicons name="camera" size={32} color="#fff" />
        </TouchableOpacity>
        <TouchableOpacity
          style={styles.iconButton}
          onPress={selectImageFromLibrary}
        >
          <Ionicons name="image" size={32} color="#fff" />
        </TouchableOpacity>
      </View>

      <TextInput
        style={styles.input}
        placeholder="Nome de usuário"
        placeholderTextColor="#DDD"
        value={username}
        onChangeText={setUsername}
      />

      <View style={styles.passwordContainer}>
        <TextInput
          style={styles.input}
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
        placeholder="Nome completo (opcional)"
        placeholderTextColor="#DDD"
        value={fullName}
        onChangeText={setFullName}
      />

      {!loading ? (
        <TouchableOpacity
          style={isDisabled ? styles.buttonDisabled : styles.button}
          onPress={handleRegister}
          disabled={isDisabled}
        >
          <Text style={styles.buttonText}>Cadastrar</Text>
        </TouchableOpacity>
      ) : (
        <Text style={styles.buttonText}>Cadastrando...</Text>
      )}
    </View>
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
    top: 40,
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
    color: '#fff' as string,
    fontSize: 20,
    textAlign: 'center',
  },
  container: {
    alignItems: 'center',
    backgroundColor: '#310d6b' as string,
    flex: 1,
    justifyContent: 'center',
    padding: 20,
  },
  icon: {
    marginBottom: 15,
    marginLeft: 12,
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
    borderColor: '#9b59b6' as string,
    borderRadius: 5,
    borderWidth: 1,
    color: '#fff' as string,
    marginBottom: 20,
    padding: 10,
    width: '85%',
  },
  logo: {
    color: '#fff' as string,
    fontSize: 50,
    fontWeight: 'bold',
    marginBottom: 50,
  },
  passwordContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    width: '85%',
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
