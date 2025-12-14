import { DefaultTheme, ThemeProvider } from '@react-navigation/native';
import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import 'react-native-reanimated';

export default function RootLayout() {
  return (
    <ThemeProvider value={DefaultTheme}>
      <Stack screenOptions={{ headerBackTitle: '戻る' }}>
        <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
        <Stack.Screen
          name="modal"
          options={{ presentation: 'modal', title: 'Modal' }}
        />
        <Stack.Screen
          name="trash/register"
          options={{ title: 'ゴミ箱を撮影' }}
        />
        <Stack.Screen
          name="monsters/index"
          options={{ title: 'ゴミスター一覧' }}
        />
        <Stack.Screen
          name="monsters/[id]"
          options={{ title: 'ゴミスター詳細' }}
        />
      </Stack>
      <StatusBar style="auto" />
    </ThemeProvider>
  );
}
