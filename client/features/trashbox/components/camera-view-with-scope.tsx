import { forwardRef } from 'react';
import { StyleSheet, View } from 'react-native';
import { CameraView } from 'expo-camera';
import { Colors } from '@/constants/theme';

type Props = {
  facing?: 'front' | 'back';
};

export const CameraViewWithScope = forwardRef<CameraView, Props>(
  ({ facing = 'back' }, ref) => {
    return (
      <View style={styles.container}>
        <CameraView ref={ref} style={styles.camera} facing={facing} />
        <View style={styles.targetOverlay}>
          <View style={styles.scopeOuter}>
            <View style={[styles.scopeLine, styles.scopeLineTop]} />
            <View style={[styles.scopeLine, styles.scopeLineBottom]} />
            <View style={[styles.scopeLine, styles.scopeLineLeft]} />
            <View style={[styles.scopeLine, styles.scopeLineRight]} />
            <View style={styles.centerDot} />
          </View>
        </View>
      </View>
    );
  }
);

CameraViewWithScope.displayName = 'CameraViewWithScope';

const styles = StyleSheet.create({
  container: {
    width: '100%',
    aspectRatio: 1,
    borderRadius: 16,
    overflow: 'hidden',
    backgroundColor: '#000',
  },
  camera: {
    flex: 1,
  },
  targetOverlay: {
    ...StyleSheet.absoluteFillObject,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scopeOuter: {
    width: 180,
    height: 180,
    borderRadius: 90,
    borderWidth: 2,
    borderColor: Colors.light.tint,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scopeLine: {
    position: 'absolute',
    backgroundColor: Colors.light.tint,
  },
  scopeLineTop: {
    width: 1.5,
    height: 60,
    top: -20,
  },
  scopeLineBottom: {
    width: 1.5,
    height: 60,
    bottom: -20,
  },
  scopeLineLeft: {
    width: 60,
    height: 1.5,
    left: -20,
  },
  scopeLineRight: {
    width: 60,
    height: 1.5,
    right: -20,
  },
  centerDot: {
    position: 'absolute',
    width: 6,
    height: 6,
    borderRadius: 3,
    backgroundColor: Colors.light.tint,
  },
});
