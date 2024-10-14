import React from 'react';
import { Image, StyleSheet, View } from 'react-native';

import splash from '../../../../../assets/images/splash.png';

export const LoadingScreen: React.FC = () => {
  return (
    <View style={styles.container}>
      <Image source={splash} style={styles.image}></Image>
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
  image: {
    resizeMode: 'center',
  },
});
