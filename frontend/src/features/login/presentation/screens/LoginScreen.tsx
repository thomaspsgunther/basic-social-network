import { Ionicons } from '@expo/vector-icons'; // Import Ionicons
import React, { useContext, useState } from 'react';
import {
  Alert,
  Keyboard,
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
  const [loading, setLoading] = useState<boolean>(false);
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [isPasswordVisible, setIsPasswordVisible] = useState<boolean>(false);

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
    marginBottom: 20,
    marginTop: 10,
    paddingHorizontal: 20,
    paddingVertical: 12,
  },
  buttonDisabled: {
    backgroundColor: 'gray' as string,
    borderRadius: 5,
    marginBottom: 20,
    marginTop: 10,
    paddingHorizontal: 20,
    paddingVertical: 12,
  },
  buttonText: {
    color: '#fff' as string,
    fontSize: 20,
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
  signUpButton: {
    marginTop: 10,
  },
  signUpText: {
    color: '#dda0dd' as string,
    fontSize: 16,
    textDecorationLine: 'underline',
  },
});
