import { useState } from 'react';
import { Alert } from 'react-native';
import { router } from 'expo-router';
import { TrashboxPresentational } from './presentational';
import { useCamera } from './hooks/useCamera';
import { createMonster } from '@/lib/client';

export const TrashboxContainer = () => {
  const {
    permission,
    requestPermission,
    photo,
    location,
    description,
    setDescription,
    cameraRef,
    takePhoto,
    retake,
    reset,
  } = useCamera();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleRegister = async () => {
    if (!photo || !location) {
      Alert.alert('エラー', '写真と位置情報が必要です');
      return;
    }

    if (!description.trim()) {
      Alert.alert('エラー', 'ニックネームを入力してください');
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await createMonster({
        nickname: description,
        latitude: location.latitude,
        longitude: location.longitude,
        image: photo,
      });

      reset();
      router.push(`/monsters/${response.data.monsterid}?fromRegister=true`);
    } catch (error) {
      console.error('モンスター登録エラー:', error);
      Alert.alert(
        'エラー',
        `登録に失敗しました: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <TrashboxPresentational
      permission={permission}
      requestPermission={requestPermission}
      photo={photo}
      description={description}
      isSubmitting={isSubmitting}
      onDescriptionChange={setDescription}
      cameraRef={cameraRef}
      onTakePhoto={takePhoto}
      onRetake={retake}
      onRegister={handleRegister}
    />
  );
};
