import { BottomTabBarProps } from '@react-navigation/bottom-tabs';
import React from 'react';
import { StyleSheet, TouchableOpacity, View } from 'react-native';

import { useAppTheme } from '../context/ThemeContext';
import { appColors } from '../theme/appColors';
import { darkTheme, lightTheme } from '../theme/appTheme';

const CustomTabBar = ({
  state,
  descriptors,
  navigation,
}: BottomTabBarProps) => {
  const { isDarkMode } = useAppTheme();
  const currentTheme = isDarkMode ? darkTheme : lightTheme;
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  return (
    <View style={currentTheme.bottomTabBar}>
      {state.routes.map((route, index) => {
        const { options } = descriptors[route.key];
        const isFocused = state.index === index;

        const iconName =
          options.tabBarIcon &&
          options.tabBarIcon({
            focused: isFocused,
            color: isFocused ? currentColors.primary : currentColors.unfocused,
            size: 32,
          });
        const onPress = () => {
          const event = navigation.emit({
            type: 'tabPress',
            target: route.key,
            canPreventDefault: true,
          });

          if (!isFocused && !event.defaultPrevented) {
            navigation.navigate(route.name);
          }
        };

        return (
          <TouchableOpacity
            key={route.key}
            onPress={onPress}
            style={styles.tab}
          >
            {iconName}
          </TouchableOpacity>
        );
      })}
    </View>
  );
};

const styles = StyleSheet.create({
  tab: {
    alignItems: 'center',
    flex: 1,
  },
});

export default CustomTabBar;
