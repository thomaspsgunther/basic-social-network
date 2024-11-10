import { Ionicons } from '@expo/vector-icons';
import React, { useEffect, useRef, useState } from 'react';
import {
  Dimensions,
  Modal,
  StyleSheet,
  Text,
  TouchableOpacity,
  TouchableWithoutFeedback,
  View,
} from 'react-native';

import { useAppTheme } from '../context/ThemeContext';
import { appColors } from '../theme/appColors';

export interface IconDropdownOption {
  label: string;
  iconName?: keyof typeof Ionicons.glyphMap;
  onSelect?: () => void;
}

interface IconDropdownProps {
  options: IconDropdownOption[];
}

const { height: screenHeight, width: screenWidth } = Dimensions.get('window');

export const IconDropdown: React.FC<IconDropdownProps> = ({ options }) => {
  const [modalVisible, setModalVisible] = useState(false);
  const [modalPosition, setModalPosition] = useState({ top: 0, left: 0 });
  const [modalSize, setModalSize] = useState({ width: 0, height: 0 });

  const { isDarkMode } = useAppTheme();
  const currentColors = isDarkMode ? appColors.dark : appColors.light;

  const styles = StyleSheet.create({
    icon: {
      marginRight: 10,
    },
    item: {
      alignItems: 'center',
      flexDirection: 'row',
      padding: 10,
    },
    itemText: {
      color: currentColors.text,
      fontSize: 16,
    },
    modalContainer: {
      backgroundColor: currentColors.background,
      borderRadius: 8,
      elevation: 5,
      padding: 10,
      position: 'absolute',
      shadowColor: currentColors.shadow,
      shadowOffset: { width: 0, height: 2 },
      shadowOpacity: 0.25,
      shadowRadius: 4,
    },
    overlay: {
      flex: 1,
    },
  });

  const iconRef = useRef<TouchableOpacity | null>(null);
  const modalRef = useRef<View | null>(null);

  const handlePress = () => {
    if (iconRef.current) {
      iconRef.current.measure((_, __, ___, height, px, py) => {
        setModalPosition({ top: py + height - 20, left: px + 60 });
        setModalVisible(true);
      });
    }
  };

  useEffect(() => {
    if (modalVisible && modalSize.height > 0) {
      const calculatedTop = Math.min(
        modalPosition.top,
        screenHeight - modalSize.height - 20,
      );
      const calculatedLeft = Math.min(
        modalPosition.left,
        screenWidth - modalSize.width - 20,
      );
      setModalPosition({ top: calculatedTop, left: calculatedLeft });
    }
  }, [modalVisible, modalSize]);

  return (
    <View>
      <TouchableOpacity ref={iconRef} onPress={() => handlePress()}>
        <Ionicons
          name="ellipsis-vertical"
          size={34}
          color={currentColors.icon}
        />
      </TouchableOpacity>

      <Modal
        transparent={true}
        animationType="fade"
        visible={modalVisible}
        onRequestClose={() => setModalVisible(false)}
      >
        <TouchableWithoutFeedback onPress={() => setModalVisible(false)}>
          <View style={styles.overlay} />
        </TouchableWithoutFeedback>
        <View
          ref={modalRef}
          style={[
            styles.modalContainer,
            { top: modalPosition.top, left: modalPosition.left },
          ]}
          onLayout={(event) => {
            const { width, height } = event.nativeEvent.layout;
            setModalSize({ width, height });
          }}
        >
          {options.map((option, index) => (
            <TouchableOpacity
              key={index}
              onPress={() => {
                if (option.onSelect) {
                  option.onSelect();
                }
                setModalVisible(false);
              }}
              style={styles.item}
            >
              {option.iconName && (
                <Ionicons
                  name={option.iconName}
                  size={24}
                  color={currentColors.icon}
                  style={styles.icon}
                />
              )}
              <Text style={styles.itemText}>{option.label}</Text>
            </TouchableOpacity>
          ))}
        </View>
      </Modal>
    </View>
  );
};
