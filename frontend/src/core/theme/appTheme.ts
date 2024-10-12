import { StyleSheet } from 'react-native';

const lightTheme = StyleSheet.create({
  background: {
    backgroundColor: '#F7F7F7' as string,
    color: '#000000' as string,
  },
});

const darkTheme = StyleSheet.create({
  background: {
    backgroundColor: '#1C1C1C' as string,
    color: '#FFFFFF' as string,
  },
});

export { darkTheme, lightTheme };
