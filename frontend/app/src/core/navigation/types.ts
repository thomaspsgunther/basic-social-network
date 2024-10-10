import { StackNavigationProp } from '@react-navigation/stack';

// Define your stack navigator types
export type RootStackParamList = {
  Login: undefined; // No parameters for login
  Home: undefined; // Other screens
  Profile: undefined;
  // Add more routes as necessary
};

export type LoginScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'Login'
>;

// If you have other navigators, you can define them here
// export type OtherNavigatorProp = StackNavigationProp<OtherStackParamList, 'SomeScreen'>;

// For now, you may just need the Login screen navigation prop
