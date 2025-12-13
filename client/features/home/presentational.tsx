import { StyleSheet, Text, View, Animated } from 'react-native';
import { useEffect, useRef, useState } from 'react';
import MapView, { Marker } from 'react-native-maps';
import * as Location from 'expo-location';
import { GestureHandlerRootView } from 'react-native-gesture-handler';
import { TrushRegisterLink } from './components/TrushRegisterLink';
import { TrashModal } from './components/TrashModal';
import { TrashPlot } from '@/components/trash-plot';
import { TrashBin } from '@/types/model';

interface HomePresentationalProps {
  location: Location.LocationObject | null;
  errorMsg: string | null;
  trashBins: TrashBin[];
}

export const HomePresentational = ({
  location,
  errorMsg,
  trashBins,
}: HomePresentationalProps) => {
  const scaleAnim = useRef(new Animated.Value(1)).current;
  const [selectedTrashBin, setSelectedTrashBin] = useState<TrashBin | null>(
    null
  );
  const [modalVisible, setModalVisible] = useState(false);

  const handleTrashBinPress = (trashBin: TrashBin) => {
    setSelectedTrashBin(trashBin);
    if (!modalVisible) {
      setModalVisible(true);
    }
  };

  const handleModalClose = () => {
    setModalVisible(false);
    setSelectedTrashBin(null);
  };

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
    <GestureHandlerRootView style={styles.container}>
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

        {/* ゴミ箱マーカー */}
        {trashBins.map((trashBin) => (
          <TrashPlot
            key={trashBin.id}
            trashBin={trashBin}
            isSelected={selectedTrashBin?.id === trashBin.id}
            onPress={handleTrashBinPress}
          />
        ))}
      </MapView>
      <View style={styles.fabContainer}>
        <TrushRegisterLink />
      </View>

      <TrashModal
        visible={modalVisible}
        trashBin={selectedTrashBin}
        onClose={handleModalClose}
      />
    </GestureHandlerRootView>
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
