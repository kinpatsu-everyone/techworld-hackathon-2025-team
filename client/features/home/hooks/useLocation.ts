import { useEffect, useState } from 'react';
import * as Location from 'expo-location';

export const useLocation = () => {
  const [location, setLocation] = useState<Location.LocationObject | null>(
    null
  );
  const [heading, setHeading] = useState<number>(0);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  useEffect(() => {
    let subscription: Location.LocationSubscription | null = null;
    let headingSubscription: Location.LocationSubscription | null = null;

    const startLocationTracking = async () => {
      let { status } = await Location.requestForegroundPermissionsAsync();
      if (status !== 'granted') {
        setErrorMsg('Permission to access location was denied');
        return;
      }

      subscription = await Location.watchPositionAsync(
        {
          accuracy: Location.Accuracy.High,
          timeInterval: 1000,
          distanceInterval: 10,
        },
        (location) => {
          setLocation(location);
        }
      );

      headingSubscription = await Location.watchHeadingAsync((headingData) => {
        setHeading(headingData.trueHeading || 0);
      });
    };

    startLocationTracking();

    return () => {
      if (subscription) {
        subscription.remove();
      }
      if (headingSubscription) {
        headingSubscription.remove();
      }
    };
  }, []);

  return { location, heading, errorMsg };
};
