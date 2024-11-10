import { Ionicons } from '@expo/vector-icons';
import { RouteProp, useNavigation, useRoute } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { useContext } from 'react';
import { Image, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { FlatList } from 'react-native-gesture-handler';

import { AuthContext } from '@/src/core/context/AuthContext';
import { useAppTheme } from '@/src/core/context/ThemeContext';
import { FeedStackParamList } from '@/src/core/navigation/types';
import { appColors } from '@/src/core/theme/appColors';
import { darkTheme, lightTheme } from '@/src/core/theme/appTheme';
import { User } from '@/src/features/shared/data/models/User';

export const UserListScreen: React.FC = () => {
  const navigation =
    useNavigation<StackNavigationProp<FeedStackParamList, 'UserList'>>();

  const route = useRoute<RouteProp<FeedStackParamList, 'UserList'>>();
  const { users, title } = route.params;

  const canGoBack = navigation.canGoBack();

  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('userlistscreen must be used within an authprovider');
  }

  const { authUser } = context;

  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  const goToUser = async (id: string) => {
    navigation.push('UserProfile', { userId: id });
  };

  return (
    <View style={currentTheme.container}>
      {!users && canGoBack && (
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={currentTheme.backButton}
        >
          <Ionicons name="arrow-back" size={40} color={currentColors.icon} />
        </TouchableOpacity>
      )}

      {users && (
        <>
          <View style={styles.listHeader}>
            {canGoBack && (
              <TouchableOpacity onPress={() => navigation.goBack()}>
                <Ionicons
                  name="arrow-back"
                  size={40}
                  color={currentColors.icon}
                />
              </TouchableOpacity>
            )}

            {title && (
              <Text style={currentTheme.titleText}>{`   ${title}`}</Text>
            )}
          </View>

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
                          size={45}
                          color="black"
                        ></Ionicons>
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
          ></FlatList>
        </>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  avatar: {
    borderRadius: 100,
    height: 45,
    width: 45,
  },
  avatarPlaceholder: {
    alignItems: 'center',
    backgroundColor: '#ccc' as string,
    borderRadius: 100,
    height: 45,
    justifyContent: 'center',
    width: 45,
  },
  flatListContainer: {
    flexGrow: 1,
    paddingTop: 20,
  },
  listHeader: {
    alignItems: 'center',
    flexDirection: 'row',
    marginTop: 50,
    paddingBottom: 16,
    paddingLeft: 20,
    paddingRight: 20,
    width: '100%',
  },
  row: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'flex-start',
    paddingBottom: 10,
    paddingLeft: 10,
    width: '100%',
  },
});
