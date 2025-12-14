import { ReactNode } from 'react';
import { StyleSheet, ViewStyle } from 'react-native';
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withSpring,
} from 'react-native-reanimated';
import {
  Gesture,
  GestureDetector,
  GestureHandlerRootView,
} from 'react-native-gesture-handler';

type Props = {
  children: ReactNode;
  style?: ViewStyle;
};

const ROTATION_FACTOR = 15; // 最大回転角度

export function TiltCard({ children, style }: Props) {
  const rotateX = useSharedValue(0);
  const rotateY = useSharedValue(0);
  const scale = useSharedValue(1);

  const gesture = Gesture.Pan()
    .onBegin(() => {
      scale.value = withSpring(1.02, { damping: 15 });
    })
    .onUpdate((event) => {
      // タッチ位置に応じて回転（中心からの距離で計算）
      rotateY.value = (event.translationX / 100) * ROTATION_FACTOR;
      rotateX.value = -(event.translationY / 100) * ROTATION_FACTOR;
    })
    .onEnd(() => {
      // 離すと元に戻る
      rotateX.value = withSpring(0, { damping: 10, stiffness: 100 });
      rotateY.value = withSpring(0, { damping: 10, stiffness: 100 });
      scale.value = withSpring(1, { damping: 15 });
    });

  const animatedStyle = useAnimatedStyle(() => ({
    transform: [
      { perspective: 1000 },
      { rotateX: `${rotateX.value}deg` },
      { rotateY: `${rotateY.value}deg` },
      { scale: scale.value },
    ],
  }));

  return (
    <GestureHandlerRootView style={styles.gestureRoot}>
      <GestureDetector gesture={gesture}>
        <Animated.View style={[styles.card, style, animatedStyle]}>
          {children}
        </Animated.View>
      </GestureDetector>
    </GestureHandlerRootView>
  );
}

const styles = StyleSheet.create({
  gestureRoot: {
    width: '100%',
  },
  card: {
    width: '100%',
  },
});
