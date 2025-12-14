import { Platform, StatusBar } from 'react-native';
import Constants from 'expo-constants';

export const STATUSBAR_HEIGHT =
  Platform.OS === 'ios' ? Constants.statusBarHeight : StatusBar.currentHeight || 0;
