import { useState, useRef } from 'react';
import { StyleSheet, View, Pressable, Text } from 'react-native';
import { CameraView, useCameraPermissions } from 'expo-camera';
import { Image } from 'expo-image';
import { ThemedView } from '@/components/themed-view';
import { ThemedText } from '@/components/themed-text';

export default function CameraTestScreen() {
  const [permission, requestPermission] = useCameraPermissions();
  const [photo, setPhoto] = useState<string | null>(null);
  const cameraRef = useRef<CameraView>(null);

  // 権限がまだ読み込まれていない
  if (!permission) {
    return (
      <ThemedView style={styles.container}>
        <ThemedText>カメラ権限を確認中...</ThemedText>
      </ThemedView>
    );
  }

  // 権限がない場合
  if (!permission.granted) {
    return (
      <ThemedView style={styles.container}>
        <ThemedText style={styles.message}>
          カメラを使用するには権限が必要です
        </ThemedText>
        <Pressable style={styles.button} onPress={requestPermission}>
          <Text style={styles.buttonText}>権限を許可する</Text>
        </Pressable>
      </ThemedView>
    );
  }

  // 撮影処理
  const takePhoto = async () => {
    if (cameraRef.current) {
      const result = await cameraRef.current.takePictureAsync();
      if (result) {
        setPhoto(result.uri);
      }
    }
  };

  // 撮り直し
  const retake = () => {
    setPhoto(null);
  };

  // 撮影済みの場合はプレビュー表示
  if (photo) {
    return (
      <ThemedView style={styles.container}>
        <ThemedText type="subtitle" style={styles.title}>
          撮影結果
        </ThemedText>
        <View style={styles.previewContainer}>
          <Image
            source={{ uri: photo }}
            style={styles.preview}
            contentFit="cover"
            contentPosition="center"
          />
        </View>
        <View style={styles.buttonRow}>
          <Pressable
            style={[styles.button, styles.retakeButton]}
            onPress={retake}
          >
            <Text style={styles.buttonText}>撮り直す</Text>
          </Pressable>
          <Pressable
            style={[styles.button, styles.useButton]}
            onPress={() => {
              // ここで画像を使う処理
              alert('この画像を使用します！');
            }}
          >
            <Text style={styles.buttonText}>使用する</Text>
          </Pressable>
        </View>
      </ThemedView>
    );
  }

  // カメラビュー
  return (
    <ThemedView style={styles.container}>
      <ThemedText type="subtitle" style={styles.title}>
        1:1 カメラ
      </ThemedText>

      {/* 1:1のカメラプレビュー */}
      <View style={styles.cameraContainer}>
        <CameraView ref={cameraRef} style={styles.camera} facing="back" />
      </View>

      {/* 撮影ボタン */}
      <Pressable
        style={({ pressed }) => [
          styles.captureButton,
          pressed && styles.captureButtonPressed,
        ]}
        onPress={takePhoto}
      >
        <View style={styles.captureButtonInner} />
      </Pressable>

      <ThemedText style={styles.hint}>タップして撮影</ThemedText>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  title: {
    marginBottom: 20,
  },
  message: {
    textAlign: 'center',
    marginBottom: 20,
  },
  cameraContainer: {
    width: '100%',
    aspectRatio: 1,
    borderRadius: 16,
    overflow: 'hidden',
    backgroundColor: '#000',
  },
  camera: {
    flex: 1,
  },
  captureButton: {
    marginTop: 30,
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: '#fff',
    borderWidth: 5,
    borderColor: '#007AFF',
    alignItems: 'center',
    justifyContent: 'center',
  },
  captureButtonPressed: {
    transform: [{ scale: 0.95 }],
    borderColor: '#0056B3',
  },
  captureButtonInner: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: '#007AFF',
  },
  hint: {
    marginTop: 16,
    opacity: 0.6,
  },
  previewContainer: {
    width: '100%',
    aspectRatio: 1,
    borderRadius: 16,
    overflow: 'hidden',
  },
  preview: {
    flex: 1,
  },
  buttonRow: {
    flexDirection: 'row',
    gap: 16,
    marginTop: 30,
  },
  button: {
    paddingVertical: 14,
    paddingHorizontal: 24,
    borderRadius: 12,
    backgroundColor: '#007AFF',
  },
  retakeButton: {
    backgroundColor: '#8E8E93',
  },
  useButton: {
    backgroundColor: '#34C759',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
