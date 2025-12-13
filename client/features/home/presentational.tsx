import { StyleSheet, Text, View, Animated } from 'react-native';
import { useEffect, useRef } from 'react';
import MapView, { Marker } from 'react-native-maps';
import * as Location from 'expo-location';
import { TrushRegisterLink } from './components/TrushRegisterLink';

interface HomePresentationalProps {
  location: Location.LocationObject | null;
  errorMsg: string | null;
}

export const HomePresentational = ({
  location,
  errorMsg,
}: HomePresentationalProps) => {
  const scaleAnim = useRef(new Animated.Value(1)).current;

  useEffect(() => {
    const pulseAnimation = Animated.loop(
      Animated.sequence([
        Animated.timing(scaleAnim, {
          toValue: 1.2,
          duration: 1000,
          useNativeDriver: true,
        }),
        Animated.timing(scaleAnim, {
          toValue: 1,
          duration: 1000,
          useNativeDriver: true,
        }),
      ])
    );

    pulseAnimation.start();

    return () => pulseAnimation.stop();
  }, [scaleAnim]);

  if (errorMsg) {
    return (
      <View style={styles.errorContainer}>
        <Text style={styles.errorText}>エラーが発生しました。</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <MapView
        style={styles.map}
        region={
          location
            ? {
                latitude: location.coords.latitude,
                longitude: location.coords.longitude,
                latitudeDelta: 0.01,
                longitudeDelta: 0.01,
              }
            : undefined
        }
      >
        {location && (
          <Marker
            coordinate={{
              latitude: location.coords.latitude,
              longitude: location.coords.longitude,
            }}
            title="Current Location"
          >
            <View style={styles.markerContainer}>
              <Animated.View
                style={[
                  styles.greenCircle,
                  { transform: [{ scale: scaleAnim }] },
                ]}
              />
            </View>
          </Marker>
        )}
      </MapView>
      <View style={styles.fabContainer}>
        <TrushRegisterLink />
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  map: {
    width: '100%',
    height: '100%',
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  errorText: {
    fontSize: 16,
    color: 'red',
    textAlign: 'center',
  },
  markerContainer: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  greenCircle: {
    width: 20,
    height: 20,
    backgroundColor: '#34C759',
    borderRadius: 10,
    borderWidth: 3,
    borderColor: 'white',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.3,
    shadowRadius: 3,
    elevation: 5,
  },
  fabContainer: {
    position: 'absolute',
    bottom: 20,
    right: 20,
  },
});
