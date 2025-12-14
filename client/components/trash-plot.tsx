import { Pressable, View, StyleSheet } from 'react-native';
import { SvgXml } from 'react-native-svg';
import { Marker } from 'react-native-maps';
import { TrashItem } from '@/lib/client';

interface TrashPlotProps {
  trashBin: TrashItem;
  onPress: (trashBin: TrashItem) => void;
}

const trashIconSvg = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#ffffff"><path d="M280-120q-33 0-56.5-23.5T200-200v-520h-40v-80h200v-40h240v40h200v80h-40v520q0 33-23.5 56.5T680-120H280Zm400-600H280v520h400v-520ZM360-280h80v-360h-80v360Zm160 0h80v-360h-80v360ZM280-720v520-520Z"/></svg>`;

export function TrashPlot({ trashBin, onPress }: TrashPlotProps) {
  return (
    <Marker
      key={trashBin.id}
      coordinate={{
        latitude: trashBin.latitude,
        longitude: trashBin.longitude,
      }}
      title={trashBin.nickname}
      onPress={() => onPress(trashBin)}
    >
      <Pressable style={styles.trashMarkerContainer}>
        <View style={styles.trashMarker}>
          <SvgXml xml={trashIconSvg} width={20} height={20} />
        </View>
      </Pressable>
    </Marker>
  );
}

const styles = StyleSheet.create({
  trashMarkerContainer: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  trashMarker: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: '#34C759',
    borderWidth: 3,
    borderColor: 'white',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.3,
    shadowRadius: 3,
    elevation: 5,
  },
  trashMarkerPressed: {
    transform: [{ scale: 1.1 }],
  },
  selectedIndicator: {
    position: 'absolute',
    bottom: -4,
    right: -4,
    width: 12,
    height: 12,
    backgroundColor: '#007AFF',
    borderRadius: 6,
    borderWidth: 2,
    borderColor: 'white',
  },
});
