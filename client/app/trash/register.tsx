import { useState, useRef } from 'react';
import {
  StyleSheet,
  View,
  Pressable,
  Text,
  TextInput,
  KeyboardAvoidingView,
  ScrollView,
  Platform,
} from 'react-native';
import { CameraView, useCameraPermissions } from 'expo-camera';
import { Image } from 'expo-image';
import * as Location from 'expo-location';
import { Colors } from '@/constants/theme';

type LocationData = {
  latitude: number;
  longitude: number;
} | null;

export default function TrashRegisterScreen() {
  const [permission, requestPermission] = useCameraPermissions();
  const [photo, setPhoto] = useState<string | null>(null);
  const [location, setLocation] = useState<LocationData>(null);
  const [description, setDescription] = useState('');
  const cameraRef = useRef<CameraView>(null);

  // 権限がまだ読み込まれていない
  if (!permission) {
    return (
      <View style={styles.container}>
        <Text style={styles.subtitle}>カメラ権限を確認中...</Text>
      </View>
    );
  }

  // 権限がない場合
  if (!permission.granted) {
    return (
      <View style={styles.container}>
        <Text style={styles.title}>ゴミ箱登録</Text>
        <Text style={styles.subtitle}>
          ゴミ箱を撮影するにはカメラの権限が必要です
        </Text>
        <Pressable style={styles.button} onPress={requestPermission}>
          <Text style={styles.buttonText}>カメラを許可する</Text>
        </Pressable>
      </View>
    );
  }

  // 撮影処理
  const takePhoto = async () => {
    if (cameraRef.current) {
      // 位置情報を取得
      const { status } = await Location.requestForegroundPermissionsAsync();
      if (status === 'granted') {
        const loc = await Location.getCurrentPositionAsync({});
        setLocation({
          latitude: loc.coords.latitude,
          longitude: loc.coords.longitude,
        });
      }

      const result = await cameraRef.current.takePictureAsync();
      if (result) {
        setPhoto(result.uri);
      }
    }
  };

  // 撮り直し
  const retake = () => {
    setPhoto(null);
    setLocation(null);
    setDescription('');
  };

  // 登録処理
  const handleRegister = () => {
    console.log('=== ゴミ箱登録データ ===');
    console.log('写真URI:', photo);
    console.log('位置情報:', location);
    console.log('説明:', description);
    console.log('========================');

    alert('ゴミ箱を登録しました！（コンソールを確認）');
    setPhoto(null);
    setLocation(null);
    setDescription('');
  };

  // 撮影済みの場合はプレビュー表示
  if (photo) {
    return (
      <KeyboardAvoidingView
        style={styles.flex1}
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        keyboardVerticalOffset={Platform.OS === 'ios' ? 100 : 0}
      >
        <ScrollView
          contentContainerStyle={styles.scrollContainer}
          keyboardShouldPersistTaps="handled"
          showsVerticalScrollIndicator={false}
        >
          <Text style={styles.title}>この写真で登録しますか？</Text>
          <View style={styles.previewContainer}>
            <Image
              source={{ uri: photo }}
              style={styles.preview}
              contentFit="cover"
              contentPosition="center"
            />
          </View>

          {/* 説明入力欄 */}
          <View style={styles.inputContainer}>
            <Text style={styles.inputLabel}>説明</Text>
            <TextInput
              style={styles.descriptionInput}
              placeholder="例：築地松竹ビルの5階"
              placeholderTextColor="#999"
              value={description}
              onChangeText={setDescription}
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
              style={[styles.button, styles.registerButton]}
              onPress={handleRegister}
            >
              <Text style={styles.buttonText}>登録する</Text>
            </Pressable>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    );
  }

  // カメラビュー
  return (
    <View style={styles.container}>
      {/* 1:1のカメラプレビュー */}
      <View style={styles.cameraContainer}>
        <CameraView ref={cameraRef} style={styles.camera} facing="back" />
        {/* スコープ風オーバーレイ */}
        <View style={styles.targetOverlay}>
          <View style={styles.scopeOuter}>
            {/* 十字線 */}
            <View style={[styles.scopeLine, styles.scopeLineTop]} />
            <View style={[styles.scopeLine, styles.scopeLineBottom]} />
            <View style={[styles.scopeLine, styles.scopeLineLeft]} />
            <View style={[styles.scopeLine, styles.scopeLineRight]} />

            {/* 中心点 */}
            <View style={styles.centerDot} />
          </View>
        </View>
      </View>

      <Text style={styles.subtitle}>ゴミ箱が写るように撮影してください</Text>

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
    </View>
  );
}

const styles = StyleSheet.create({
  flex1: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  scrollContainer: {
    flexGrow: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#f5f5f5',
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 8,
    color: '#333',
  },
  subtitle: {
    fontSize: 16,
    color: Colors.light.text,
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
  targetOverlay: {
    ...StyleSheet.absoluteFillObject,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scopeOuter: {
    width: 180,
    height: 180,
    borderRadius: 90,
    borderWidth: 2,
    borderColor: Colors.light.tint,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scopeLine: {
    position: 'absolute',
    backgroundColor: Colors.light.tint,
  },
  scopeLineTop: {
    width: 1.5,
    height: 60,
    top: -20,
  },
  scopeLineBottom: {
    width: 1.5,
    height: 60,
    bottom: -20,
  },
  scopeLineLeft: {
    width: 60,
    height: 1.5,
    left: -20,
  },
  scopeLineRight: {
    width: 60,
    height: 1.5,
    right: -20,
  },
  centerDot: {
    position: 'absolute',
    width: 6,
    height: 6,
    borderRadius: 3,
    backgroundColor: Colors.light.tint,
  },
  captureButton: {
    marginTop: 30,
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: '#fff',
    borderWidth: 5,
    borderColor: '#34C759',
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.2,
    shadowRadius: 4,
    elevation: 4,
  },
  captureButtonPressed: {
    transform: [{ scale: 0.95 }],
    borderColor: '#2DA44E',
  },
  captureButtonInner: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: '#34C759',
  },
  previewContainer: {
    width: '100%',
    aspectRatio: 1,
    borderRadius: 16,
    overflow: 'hidden',
    marginTop: 16,
  },
  preview: {
    flex: 1,
  },
  inputContainer: {
    width: '100%',
    marginTop: 16,
  },
  inputLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 6,
  },
  descriptionInput: {
    width: '100%',
    padding: 12,
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 12,
    backgroundColor: '#fff',
    fontSize: 16,
    color: '#333',
  },
  locationText: {
    marginTop: 12,
    fontSize: 14,
    color: '#666',
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
  registerButton: {
    backgroundColor: '#34C759',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
