import { StyleSheet, View, Pressable, Text } from 'react-native';
import { CameraView } from 'expo-camera';
import { Colors } from '@/constants/theme';
import { CameraViewWithScope } from './components/camera-view-with-scope';
import { CaptureButton } from './components/capture-button';
import { PhotoPreview } from './components/photo-preview';

type Props = {
  permission: { granted: boolean } | null;
  requestPermission: () => void;
  photo: string | null;
  description: string;
  onDescriptionChange: (text: string) => void;
  cameraRef: React.RefObject<CameraView | null>;
  onTakePhoto: () => void;
  onRetake: () => void;
  onRegister: () => void;
};

export const TrashboxPresentational = ({
  permission,
  requestPermission,
  photo,
  description,
  onDescriptionChange,
  cameraRef,
  onTakePhoto,
  onRetake,
  onRegister,
}: Props) => {
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

  // 撮影済みの場合はプレビュー表示
  if (photo) {
    return (
      <PhotoPreview
        photoUri={photo}
        description={description}
        onDescriptionChange={onDescriptionChange}
        onRetake={onRetake}
        onRegister={onRegister}
      />
    );
  }

  // カメラビュー
  return (
    <View style={styles.container}>
      <CameraViewWithScope ref={cameraRef} />
      <Text style={styles.subtitle}>ゴミ箱が写るように撮影してください</Text>
      <CaptureButton onPress={onTakePhoto} />
    </View>
  );
};

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
