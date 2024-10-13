import { RouteProp } from '@react-navigation/native';
import React from 'react';
import { StyleSheet, Text, View } from 'react-native';

import {
  FeedScreenNavigationProp,
  FeedStackParamList,
} from '@/src/core/navigation/types';

type Props = {
  navigation: FeedScreenNavigationProp;
  route: RouteProp<FeedStackParamList, 'Feed'>;
};

export const FeedScreen: React.FC<Props> = () => {
  return (
    <View style={styles.container}>
      <Text style={styles.text}>Welcome to the Feed Screen!</Text>
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
