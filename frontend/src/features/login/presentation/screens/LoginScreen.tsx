import { Ionicons } from '@expo/vector-icons';
import { CommonActions } from '@react-navigation/native';
import Constants from 'expo-constants';
import React, { useContext, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  Keyboard,
  KeyboardAvoidingView,
  Platform,
  StyleSheet,
  Text,
  TextInput,
  View,
} from 'react-native';
import { TouchableOpacity } from 'react-native-gesture-handler';

import { AuthContext } from '@/src/core/context/AuthContext';
import { RootStackScreenProps } from '@/src/core/navigation/types';
import { User } from '@/src/features/shared/data/models/User';

export const LoginScreen: React.FC<RootStackScreenProps<'Login'>> = ({
  navigation,
}) => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [isPasswordVisible, setIsPasswordVisible] = useState<boolean>(false);

  const isDisabled: boolean = username.trim() === '' || password.trim() === '';

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('loginscreen must be used within an authprovider');
  }

  const { login, logout } = context;

  const handleLogin = async () => {
    setIsLoading(true);
    Keyboard.dismiss();
    try {
      if (username && password) {
        if (!hasSpecialCharacters(username)) {
          const userData: Omit<User, 'id'> = {
            username: username,
            password: password,
          };
          await login(userData);
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
            'Usuário não pode conter caracteres especiais',
          );
        }
      } else {
        setIsLoading(false);
        Alert.alert(
          'Oops, algo deu errado',
          'Usuário e senha precisam estar preenchidos',
        );
      }
    } catch (error) {
      if (error instanceof Error) {
        setIsLoading(false);
        logout();
        if (error.message.trim() === 'wrong username or password') {
          Alert.alert('Oops, algo deu errado', 'Usuário ou senha incorretos');
        } else {
          Alert.alert('Oops, algo deu errado', 'Por favor, tente novamente');
        }
      }
    }
  };

  const handleUsernameChange = (input: string) => {
    const noSpacesInput: string = input.replace(/\s+/g, '');
    const lowercaseUsername: string = noSpacesInput.toLowerCase();
    setUsername(lowercaseUsername);
  };

  const handlePasswordChange = (input: string) => {
    const noSpacesPassword: string = input.replace(/\s+/g, '');
    setPassword(noSpacesPassword);
  };

  const hasSpecialCharacters = (str: string) => /[^a-zA-Z0-9\s]/.test(str);

  const handleSignUpRedirect = () => {
    navigation.navigate('Register');
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={styles.container}
    >
      <Text style={styles.logo}>y</Text>

      <TextInput
        style={styles.input}
        maxLength={20}
        placeholder="Nome de usuário"
        placeholderTextColor="#DDD"
        value={username}
        onChangeText={handleUsernameChange}
      />

      <View style={styles.passwordContainer}>
        <TextInput
          style={styles.inputPassword}
          maxLength={30}
          placeholder="Senha"
          placeholderTextColor="#DDD"
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
            color="#ddd"
            style={styles.icon}
          />
        </TouchableOpacity>
      </View>

      {!isLoading ? (
        <TouchableOpacity
          style={isDisabled ? styles.buttonDisabled : styles.button}
          onPress={() => handleLogin()}
          disabled={isDisabled}
        >
          <Text style={styles.buttonText}>Entrar</Text>
        </TouchableOpacity>
      ) : (
        <ActivityIndicator size="large" color="#FFFFFF" />
      )}

      <TouchableOpacity
        style={styles.signUpButton}
        onPress={() => handleSignUpRedirect()}
        disabled={isLoading}
      >
        <Text style={styles.signUpText}>
          Ainda não tem conta? Cadastre-se aqui!
        </Text>
      </TouchableOpacity>

      <Text style={styles.versionText}>v{Constants.expoConfig?.version}</Text>
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  button: {
    backgroundColor: '#8A2BE2' as string,
    borderRadius: 5,
    marginBottom: 30,
    marginTop: 10,
    paddingHorizontal: 20,
    paddingVertical: 12,
  },
  buttonDisabled: {
    backgroundColor: 'gray' as string,
    borderRadius: 5,
    marginBottom: 30,
    marginTop: 10,
    paddingHorizontal: 20,
    paddingVertical: 12,
  },
  buttonText: {
    color: 'white' as string,
    fontSize: 20,
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
    fontSize: 100,
    fontWeight: 'bold',
    marginBottom: 50,
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
  signUpButton: {
    marginTop: 10,
  },
  signUpText: {
    color: '#dda0dd' as string,
    fontSize: 18,
    textDecorationLine: 'underline',
  },
  versionText: {
    bottom: 15,
    color: 'white' as string,
    fontSize: 14,
    position: 'absolute',
    right: 20,
  },
});
