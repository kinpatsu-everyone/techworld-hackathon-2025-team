import { Tabs } from 'expo-router';
import { View } from 'react-native';
import { HapticTab } from '@/components/haptic-tab';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { Colors } from '@/constants/theme';
import { useColorScheme } from '@/hooks/use-color-scheme';
import { Image } from 'expo-image';

export default function TabLayout() {
  const colorScheme = useColorScheme();

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: Colors[colorScheme ?? 'light'].tint,
        headerShown: false,
        tabBarButton: HapticTab,
      }}
    >
      <Tabs.Screen
        name="monsters"
        options={{
          title: 'MONSTER',
          tabBarIcon: ({ color }) => (
            <View style={{ width: 28, height: 28 }}>
              {color === Colors[colorScheme ?? 'light'].tint ? (
                <Image
                  source={require('@/assets/icons/egg-selected.png')}
                  style={{ width: 28, height: 28, resizeMode: 'contain' }}
                />
              ) : (
                <Image
                  source={require('@/assets/icons/egg-gray.png')}
                  style={{ width: 28, height: 28, resizeMode: 'contain' }}
                />
              )}
            </View>
          ),
        }}
      />
      <Tabs.Screen
        name="index"
        options={{
          title: 'MAP',
          tabBarIcon: ({ color }) => (
            <IconSymbol size={28} name="location.fill" color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="badge"
        options={{
          title: 'BADGE',
          tabBarIcon: ({ color }) => (
            <View style={{ width: 28, height: 28 }}>
              {color === Colors[colorScheme ?? 'light'].tint ? (
                <Image
                  source={require('@/assets/icons/badge-selected.png')}
                  style={{ width: 28, height: 28, resizeMode: 'contain' }}
                />
              ) : (
                <Image
                  source={require('@/assets/icons/badge-gray.png')}
                  style={{ width: 28, height: 28, resizeMode: 'contain' }}
                />
              )}
            </View>
          ),
        }}
      />
    </Tabs>
  );
}
