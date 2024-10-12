import React from 'react';
import { StyleSheet, Text, View } from 'react-native';

const LoadingScreen: React.FC = () => {
  return (
    <View style={styles.container}>
      <Text style={styles.text}>Loading...</Text>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    backgroundColor: '#310d6b' as string,
    flex: 1,
    justifyContent: 'center',
  },
  text: {
    color: '#fff' as string,
    fontSize: 24,
  },
});

export default LoadingScreen;
