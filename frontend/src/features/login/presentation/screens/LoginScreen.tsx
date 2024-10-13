import { RouteProp } from '@react-navigation/native';
import React, { useContext, useState } from 'react';
import {
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

import { AuthContext } from '@/src/core/context/AuthContext';
import {
  LoginScreenNavigationProp,
  RootStackParamList,
} from '@/src/core/navigation/types';
import { User } from '@/src/features/shared/data/models/User';

type Props = {
  navigation: LoginScreenNavigationProp;
  route: RouteProp<RootStackParamList, 'Login'>;
};

const LoginScreen: React.FC<Props> = ({ navigation }) => {
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('loginscreen must be used within an authprovider');
  }

  const { login } = context;

  const handleLogin = () => {
    if (username && password) {
      const userData: Omit<User, 'id'> = {
        username: username,
        password: password,
      };
      login(userData);
      navigation.navigate('Home');
    } else {
      console.log('display some alert I guess');
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

      <TouchableOpacity style={styles.button} onPress={handleLogin}>
        <Text style={styles.buttonText}>Entrar</Text>
      </TouchableOpacity>

      <TouchableOpacity
        style={styles.signUpButton}
        onPress={handleSignUpRedirect}
      >
        <Text style={styles.signUpText}>Não é cadastrado? Cadastre-se!</Text>
      </TouchableOpacity>
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

export default LoginScreen;
