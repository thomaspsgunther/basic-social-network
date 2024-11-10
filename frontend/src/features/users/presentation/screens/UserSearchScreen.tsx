import { Ionicons } from '@expo/vector-icons';
import { useContext, useEffect, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  Image,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import { FlatList } from 'react-native-gesture-handler';

import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { UserSearchStackScreenProps } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { User } from '@/src/features/shared/data/models/User';

import { UserRepositoryImpl } from '../../data/repositories/UserRepositoryImpl';
import { UserUsecaseImpl } from '../../domain/usecases/UserUsecase';

export const UserSearchScreen: React.FC<
  UserSearchStackScreenProps<'UserSearch'>
> = ({ navigation }) => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [debouncedSearchTerm, setDebouncedSearchTerm] = useState<string>('');
  const [users, setUsers] = useState<User[]>([]);
  const userRepository = new UserRepositoryImpl();
  const userUsecase = new UserUsecaseImpl(userRepository);

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('usersearchscreen must be used within an authprovider');
  }

  const { authUser } = context;

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedSearchTerm(searchTerm);
    }, 500);

    return () => {
      clearTimeout(handler);
    };
  }, [searchTerm]);

  useEffect(() => {
    const findUsers = async (searchStr: string) => {
      if (searchStr.length > 2) {
        setIsLoading(true);
        try {
          const foundUsers: User[] =
            await userUsecase.getUsersBySearch(searchStr);

          setIsLoading(false);
          setUsers(foundUsers);
        } catch (_error) {
          setIsLoading(false);
          Alert.alert('Oops, algo deu errado');
        }
      }
    };

    if (debouncedSearchTerm) {
      findUsers(debouncedSearchTerm);
    }
  }, [debouncedSearchTerm]);

  const handleSearchTermChange = (input: string) => {
    setSearchTerm(input);

    if (input.length <= 2) {
      setUsers([]);
    }
  };
  const clearSearch = async () => {
    setSearchTerm('');
    setUsers([]);
  };

  const goToUser = async (id: string) => {
    navigation.push('UserProfile', { userId: id });
  };

  return (
    <View style={currentTheme.container}>
      <View style={styles.listHeader}>
        <View style={currentTheme.inputIconContainer}>
          <Ionicons
            name="search"
            size={26}
            color={currentColors.icon}
            style={styles.icon}
          />

          <TextInput
            style={currentTheme.inputIcon}
            maxLength={50}
            placeholder="Encontre um usuário"
            placeholderTextColor={currentColors.placeholderText}
            value={searchTerm}
            onChangeText={handleSearchTermChange}
          />

          {searchTerm.length > 0 && (
            <TouchableOpacity onPress={() => clearSearch()}>
              <Ionicons
                name="close-outline"
                size={26}
                color={currentColors.icon}
                style={styles.clearIcon}
              />
            </TouchableOpacity>
          )}
        </View>
      </View>

      {!isLoading ? (
        <FlatList
          data={users}
          keyExtractor={(user) => user.id}
          renderItem={({ item }: { item: User }) => {
            return (
              <View style={styles.row}>
                <TouchableOpacity
                  style={styles.row}
                  onPress={() => goToUser(item.id)}
                  disabled={authUser ? item.id === authUser!.id : true}
                >
                  {item.avatar ? (
                    <Image
                      source={{
                        uri: `data:image/jpeg;base64,${item.avatar}`,
                      }}
                      style={styles.avatar}
                      resizeMode="contain"
                    />
                  ) : (
                    <View style={styles.avatarPlaceholder}>
                      <Ionicons
                        name="person-circle-outline"
                        size={55}
                        color="black"
                      />
                    </View>
                  )}

                  <View>
                    <Text style={currentTheme.textBold}>
                      {`   ${item.username}`}
                    </Text>

                    {item.fullName && (
                      <Text style={currentTheme.text}>
                        {`   ${item.fullName}`}
                      </Text>
                    )}
                  </View>
                </TouchableOpacity>
              </View>
            );
          }}
          showsVerticalScrollIndicator={false}
          contentContainerStyle={styles.flatListContainer}
        />
      ) : (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={currentColors.icon} />
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  avatar: {
    borderRadius: 100,
    height: 55,
    width: 55,
  },
  avatarPlaceholder: {
    alignItems: 'center',
    backgroundColor: '#ccc' as string,
    borderRadius: 100,
    height: 55,
    justifyContent: 'center',
    width: 55,
  },
  clearIcon: {
    marginRight: 15,
    padding: 10,
  },
  flatListContainer: {
    flexGrow: 1,
  },
  icon: {
    marginLeft: 8,
  },
  listHeader: {
    alignItems: 'center',
    marginTop: 50,
    paddingBottom: 20,
    width: '100%',
  },
  loadingContainer: {
    alignItems: 'center',
    flex: 1,
    justifyContent: 'center',
  },
  row: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'flex-start',
    paddingBottom: 15,
    paddingLeft: 10,
    width: 420,
  },
});
