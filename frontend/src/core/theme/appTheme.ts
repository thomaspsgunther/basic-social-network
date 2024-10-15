import { StyleSheet } from 'react-native';

const lightTheme = StyleSheet.create({
  button: {
    backgroundColor: '#310d6b' as string,
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
    backgroundColor: '#F7F7F7' as string,
    flex: 1,
    justifyContent: 'center',
    padding: 20,
  },
  text: {
    color: 'black' as string,
  },
});

const darkTheme = StyleSheet.create({
  button: {
    backgroundColor: '#310d6b' as string,
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
    backgroundColor: '#1C1C1C' as string,
    flex: 1,
    justifyContent: 'center',
    padding: 20,
  },
  text: {
    color: 'white' as string,
  },
});

export { darkTheme, lightTheme };
