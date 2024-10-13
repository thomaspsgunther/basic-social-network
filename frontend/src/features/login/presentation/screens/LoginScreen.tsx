import React, { useContext, useState } from 'react';
import {
  Alert,
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

export const LoginScreen: React.FC<RootStackScreenProps<'Login'>> = ({
  navigation,
}) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');

  const isDisabled = username.trim() === '' || password.trim() === '';

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('loginscreen must be used within an authprovider');
  }

  const { login, logout } = context;

  const handleLogin = async () => {
    setLoading(true);
    Keyboard.dismiss();
    try {
      if (username && password) {
        const userData: Omit<User, 'id'> = {
          username: username,
          password: password,
        };
        await login(userData);
        navigation.navigate('Tabs');
      } else {
        setLoading(false);
        Alert.alert(
          'Oops, algo deu errado',
          'Usuário e senha precisam estar preenchidos',
        );
      }
    } catch (error) {
      if (error instanceof Error) {
        setLoading(false);
        logout();
        if (error.message.trim() === 'wrong username or password') {
          Alert.alert('Oops, algo deu errado', 'Usuário ou senha incorretos');
        } else {
          Alert.alert('Oops, algo deu errado', 'Por favor, tente novamente');
        }
      }
    }
  };

  const handleSignUpRedirect = () => {
    navigation.navigate('Register');
  };

  return (
    <View style={styles.container}>
      <Text style={styles.logo}>y</Text>

      <TextInput
        style={styles.input}
        placeholder="Usuário"
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

      {!loading ? (
        <TouchableOpacity
          style={isDisabled ? styles.buttonDisabled : styles.button}
          onPress={handleLogin}
          disabled={isDisabled}
        >
          <Text style={styles.buttonText}>Entrar</Text>
        </TouchableOpacity>
      ) : (
        <Text style={styles.buttonText}>Entrando...</Text>
      )}

      {!loading ? (
        <TouchableOpacity
          style={styles.signUpButton}
          onPress={handleSignUpRedirect}
        >
          <Text style={styles.signUpText}>
            Ainda não tem conta? Cadastre-se aqui!
          </Text>
        </TouchableOpacity>
      ) : null}
    </View>
  );
};

const styles = StyleSheet.create({
  button: {
    backgroundColor: '#8A2BE2' as string,
    borderRadius: 5,
    marginBottom: 15,
    paddingHorizontal: 20,
    paddingVertical: 10,
  },
  buttonDisabled: {
    backgroundColor: 'gray' as string,
    borderRadius: 5,
    marginBottom: 15,
    paddingHorizontal: 20,
    paddingVertical: 10,
  },
  buttonText: {
    color: '#fff' as string,
    fontSize: 18,
  },
  container: {
    alignItems: 'center',
    backgroundColor: '#310d6b' as string,
    flex: 1,
    justifyContent: 'center',
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
  signUpButton: {
    marginTop: 10,
  },
  signUpText: {
    color: '#dda0dd' as string,
    fontSize: 16,
    textDecorationLine: 'underline',
  },
});
