import { router } from 'expo-router';
import { TrashboxPresentational } from './presentational';
import { useCamera } from './hooks/useCamera';

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

  const handleRegister = () => {
    console.log('=== ゴミ箱登録データ ===');
    console.log('写真URI:', photo);
    console.log('位置情報:', location);
    console.log('説明:', description);
    console.log('========================');

    // TODO: APIに送信して、返ってきたIDで遷移
    const mockMonsterId = '1';

    reset();
    router.replace(`/monsters/${mockMonsterId}?fromRegister=true`);
  };

  return (
    <TrashboxPresentational
      permission={permission}
      requestPermission={requestPermission}
      photo={photo}
      description={description}
      onDescriptionChange={setDescription}
      cameraRef={cameraRef}
      onTakePhoto={takePhoto}
      onRetake={retake}
      onRegister={handleRegister}
    />
  );
};
