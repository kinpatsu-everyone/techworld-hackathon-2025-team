import { useState } from 'react';
import { Image } from 'expo-image';
import { Platform, StyleSheet, Pressable, Alert } from 'react-native';
import * as ImagePicker from 'expo-image-picker';

import { Collapsible } from '@/components/ui/collapsible';
import { ExternalLink } from '@/components/external-link';
import ParallaxScrollView from '@/components/parallax-scroll-view';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { Fonts } from '@/constants/theme';

export default function TabTwoScreen() {
  const [selectedImage, setSelectedImage] = useState<string | null>(null);

  const takePhoto = async () => {
    const permissionResult = await ImagePicker.requestCameraPermissionsAsync();

    if (!permissionResult.granted) {
      Alert.alert(
        '権限が必要です',
        '写真を撮るにはカメラへのアクセスを許可してください。'
      );
      return;
    }

    const result = await ImagePicker.launchCameraAsync({
      allowsEditing: true,
      aspect: [1, 1],
      quality: 1,
    });

    if (!result.canceled) {
      setSelectedImage(result.assets[0].uri);
    }
  };
  return (
    <ParallaxScrollView
      headerBackgroundColor={{ light: '#D0D0D0', dark: '#353636' }}
      headerImage={
        <IconSymbol
          size={310}
          color="#808080"
          name="chevron.left.forwardslash.chevron.right"
          style={styles.headerImage}
        />
      }
    >
      <ThemedView style={styles.titleContainer}>
        <ThemedText
          type="title"
          style={{
            fontFamily: Fonts.rounded,
          }}
        >
          Camera Test
        </ThemedText>
      </ThemedView>
      {/* 画像選択セクション */}
      <ThemedView style={styles.imagePickerSection}>
        <Pressable
          style={({ pressed }) => [
            styles.pickButton,
            pressed && styles.pickButtonPressed,
          ]}
          onPress={takePhoto}
        >
          <ThemedText style={styles.pickButtonText}>📷 カメラで撮影</ThemedText>
        </Pressable>

        {selectedImage && (
          <Image
            source={{ uri: selectedImage }}
            style={styles.selectedImage}
            contentFit="cover"
            contentPosition="center"
          />
        )}
      </ThemedView>
    </ParallaxScrollView>
  );
}

const styles = StyleSheet.create({
  headerImage: {
    color: '#808080',
    bottom: -90,
    left: -35,
    position: 'absolute',
  },
  titleContainer: {
    flexDirection: 'row',
    gap: 8,
  },
  imagePickerSection: {
    marginVertical: 16,
    alignItems: 'center',
    gap: 16,
  },
  pickButton: {
    backgroundColor: '#007AFF',
    paddingVertical: 14,
    paddingHorizontal: 24,
    borderRadius: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.15,
    shadowRadius: 4,
    elevation: 3,
  },
  pickButtonPressed: {
    backgroundColor: '#0056B3',
    transform: [{ scale: 0.97 }],
  },
  pickButtonText: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '600',
  },
  selectedImage: {
    width: '100%',
    aspectRatio: 1,
    borderRadius: 12,
  },
});
