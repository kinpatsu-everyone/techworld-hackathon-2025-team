import { useState, useRef } from 'react';
import { StyleSheet, View, Pressable, Text } from 'react-native';
import { CameraView, useCameraPermissions } from 'expo-camera';
import * as Location from 'expo-location';
import { router } from 'expo-router';
import { Colors } from '@/constants/theme';
import { CameraViewWithScope } from '@/features/register/components/camera-view-with-scope';
import { CaptureButton } from '@/features/register/components/capture-button';
import { PhotoPreview } from '@/features/register/components/photo-preview';

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

    // TODO: APIに送信して、返ってきたIDで遷移
    const mockMonsterId = '1';

    setPhoto(null);
    setLocation(null);
    setDescription('');

    router.push(`/monsters/${mockMonsterId}`);
  };

  // 撮影済みの場合はプレビュー表示
  if (photo) {
    return (
      <PhotoPreview
        photoUri={photo}
        description={description}
        onDescriptionChange={setDescription}
        onRetake={retake}
        onRegister={handleRegister}
      />
    );
  }

  // カメラビュー
  return (
    <View style={styles.container}>
      <CameraViewWithScope ref={cameraRef} />
      <Text style={styles.subtitle}>ゴミ箱が写るように撮影してください</Text>
      <CaptureButton onPress={takePhoto} />
    </View>
  );
}

const styles = StyleSheet.create({
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
  button: {
    paddingVertical: 14,
    paddingHorizontal: 24,
    borderRadius: 12,
    backgroundColor: '#007AFF',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
