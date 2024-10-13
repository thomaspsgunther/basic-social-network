import { Ionicons } from '@expo/vector-icons';
import { RouteProp } from '@react-navigation/native';
import * as ImagePicker from 'expo-image-picker';
import React, { useContext, useState } from 'react';
import {
  Image,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

import { AuthContext } from '@/src/core/context/AuthContext';
import {
  RegisterScreenNavigationProp,
  RootStackParamList,
} from '@/src/core/navigation/types';
import { User } from '@/src/features/shared/data/models/User';

type Props = {
  navigation: RegisterScreenNavigationProp;
  route: RouteProp<RootStackParamList, 'Register'>;
};

export const RegisterScreen: React.FC<Props> = ({ navigation }) => {
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  const [fullName, setFullName] = useState<string>('');
  const [avatar, setAvatar] = useState<string | null>(null);
  const [avatarUri, setAvatarUri] = useState<string | null>(null);

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('feedscreen must be used within an authprovider');
  }

  const { register } = context;

  const handleRegister = () => {
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

      register(userData);
      navigation.navigate('Tabs');
    } else {
      console.log('display some alert I guess');
    }
  };

  const selectImageFromLibrary = async () => {
    const permissionResult =
      await ImagePicker.requestMediaLibraryPermissionsAsync();

    if (permissionResult.granted === false) {
      alert('permission to access camera roll is required!');
      return;
    }

    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ImagePicker.MediaTypeOptions.All,
      base64: true,
      quality: 1,
    });

    if (!result.canceled) {
      const base64String = result.assets[0].base64;
      const imageType = result.assets[0].type || 'image/jpeg';
      setAvatar(base64String || null);
      setAvatarUri(`data:${imageType};base64,${base64String}`);
    }
  };

  const takePhoto = async () => {
    const permissionResult = await ImagePicker.requestCameraPermissionsAsync();

    if (permissionResult.granted === false) {
      alert('permission to access the camera is required!');
      return;
    }

    const result = await ImagePicker.launchCameraAsync({
      base64: true,
      quality: 1,
    });

    if (!result.canceled) {
      const base64String = result.assets[0].base64;
      const imageType = result.assets[0].type || 'image/jpeg';
      setAvatar(base64String || null);
      setAvatarUri(`data:${imageType};base64,${base64String}`);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.logo}>y</Text>

      <TouchableOpacity onPress={takePhoto}>
        {avatarUri ? (
          <Image source={{ uri: avatarUri }} style={styles.avatar} />
        ) : (
          <View style={styles.avatarPlaceholder}>
            <Text style={styles.avatarPlaceholderText}>Tire uma foto!</Text>
          </View>
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

      <TextInput
        style={styles.input}
        placeholder="Senha"
        placeholderTextColor="#DDD"
        secureTextEntry={true}
        value={password}
        onChangeText={setPassword}
      />

      <TextInput
        style={styles.input}
        placeholder="Email (opcional)"
        placeholderTextColor="#DDD"
        value={email}
        onChangeText={setEmail}
      />

      <TextInput
        style={styles.input}
        placeholder="Nome Completo (opcional)"
        placeholderTextColor="#DDD"
        value={fullName}
        onChangeText={setFullName}
      />

      <TouchableOpacity style={styles.button} onPress={handleRegister}>
        <Text style={styles.buttonText}>Cadastrar</Text>
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  avatar: {
    borderRadius: 50,
    height: 100,
    marginBottom: 20,
    width: 100,
  },
  avatarPlaceholder: {
    alignItems: 'center',
    backgroundColor: '#ccc' as string,
    borderRadius: 50,
    height: 100,
    justifyContent: 'center',
    marginBottom: 20,
    width: 100,
  },
  avatarPlaceholderText: {
    color: '#777' as string,
  },
  button: {
    backgroundColor: '#8A2BE2' as string,
    borderRadius: 5,
    marginBottom: 15,
    paddingHorizontal: 20,
    paddingVertical: 10,
  },
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: 20,
    width: '80%',
  },
  buttonText: {
    color: '#fff' as string,
    fontSize: 18,
    textAlign: 'center',
  },
  container: {
    alignItems: 'center',
    backgroundColor: '#310d6b' as string,
    flex: 1,
    justifyContent: 'center',
  },
  iconButton: {
    alignItems: 'center',
    backgroundColor: '#8A2BE2' as string,
    borderRadius: 5,
    flex: 1,
    justifyContent: 'center',
    margin: 10,
    padding: 10,
  },
  input: {
    borderColor: '#9b59b6' as string,
    borderRadius: 5,
    borderWidth: 1,
    color: '#fff' as string,
    marginBottom: 20,
    padding: 10,
    width: '80%',
  },
  logo: {
    color: '#fff' as string,
    fontSize: 32,
    fontWeight: 'bold',
    marginBottom: 40,
  },
});
