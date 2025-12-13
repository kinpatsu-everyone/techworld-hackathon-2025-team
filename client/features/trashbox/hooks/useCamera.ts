import { useState, useRef } from 'react';
import { CameraView, useCameraPermissions } from 'expo-camera';
import * as Location from 'expo-location';

export type LocationData = {
  latitude: number;
  longitude: number;
} | null;

export const useCamera = () => {
  const [permission, requestPermission] = useCameraPermissions();
  const [photo, setPhoto] = useState<string | null>(null);
  const [location, setLocation] = useState<LocationData>(null);
  const [description, setDescription] = useState('');
  const cameraRef = useRef<CameraView>(null);

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

  const retake = () => {
    setPhoto(null);
    setLocation(null);
    setDescription('');
  };

  const reset = () => {
    setPhoto(null);
    setLocation(null);
    setDescription('');
  };

  return {
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
  };
};
