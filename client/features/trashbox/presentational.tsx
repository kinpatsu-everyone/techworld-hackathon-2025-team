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
  isSubmitting?: boolean;
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
  isSubmitting,
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
        isLoading={isSubmitting}
        onDescriptionChange={onDescriptionChange}
        onRetake={onRetake}
        onRegister={onRegister}
      />
    );
  }

  // カメラビュー
  return (
    <View style={styles.cameraContainer}>
      <View style={styles.spacer} />
      <CameraViewWithScope ref={cameraRef} />
      <View style={styles.footer}>
        <Text style={styles.subtitle}>ゴミ箱が写るように撮影してください</Text>
        <CaptureButton onPress={onTakePhoto} />
      </View>
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
  cameraContainer: {
    height: '100%',
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    padding: 16,
  },
  spacer: {
    flexGrow: 1,
  },
  footer: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    gap: 16,
    alignItems: 'center',
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
  },
  button: {
    borderRadius: 12,
    borderWidth: 1,
    borderColor: '#007AFF',
    backgroundColor: '#007AFF',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
