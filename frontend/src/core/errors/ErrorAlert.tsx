import React from 'react';
import { Alert } from 'react-native';
import { useDispatch, useSelector } from 'react-redux';

import { RootState } from '../redux/store';
import { clearError } from './errorsSlice';

const ErrorAlert: React.FC = () => {
  const error = useSelector((state: RootState) => state.errors);
  const dispatch = useDispatch();

  React.useEffect(() => {
    if (error && error.message) {
      Alert.alert('Oops, algo deu errado', '', [
        { text: 'OK', onPress: () => dispatch(clearError()) },
      ]);
    }
  }, [error, dispatch]);

  return null;
};

export default ErrorAlert;
