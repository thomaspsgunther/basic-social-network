import { Ionicons } from '@expo/vector-icons';
import React, { useState } from 'react';
import {
  Modal,
  StyleSheet,
  Text,
  TouchableOpacity,
  TouchableWithoutFeedback,
  View,
} from 'react-native';

// import { useTheme } from '../context/ThemeContext';
// import { appColors } from '../theme/appColors';

export interface IconDropdownOption {
  label: string;
  iconName?: keyof typeof Ionicons.glyphMap;
  onSelect?: () => void;
}

interface IconDropdownProps {
  options: IconDropdownOption[];
}

// const { isDarkMode } = useTheme();
// const currentColors = isDarkMode ? appColors.dark : appColors.light;

export const IconDropdown: React.FC<IconDropdownProps> = ({ options }) => {
  const [modalVisible, setModalVisible] = useState(false);

  return (
    <View style={styles.container}>
      <TouchableOpacity onPress={() => setModalVisible(true)}>
        <Ionicons name="ellipsis-vertical" size={24} color="#000" />
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
        <View style={styles.modalContainer}>
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
                  size={20}
                  color="#000"
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

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    flex: 1,
    justifyContent: 'center',
  },
  icon: {
    marginRight: 10,
  },
  item: {
    alignItems: 'center',
    flexDirection: 'row',
    padding: 10,
  },
  itemText: {
    fontSize: 16,
  },
  modalContainer: {
    // backgroundColor: currentColors.background,
    borderRadius: 8,
    elevation: 5,
    padding: 10,
    position: 'absolute',
    right: 20,
    // shadowColor: currentColors.shadow,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.25,
    shadowRadius: 4,
    top: 40,
  },
  overlay: {
    // backgroundColor: currentColors.overlay,
    flex: 1,
  },
});
