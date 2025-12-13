import { useEffect, useState } from 'react';
import { StyleSheet, View, Dimensions } from 'react-native';
import { Image } from 'expo-image';
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withTiming,
  withSequence,
  withDelay,
  withSpring,
  Easing,
  runOnJS,
} from 'react-native-reanimated';

const { width: SCREEN_WIDTH, height: SCREEN_HEIGHT } = Dimensions.get('window');

// 卵の画像
const EGG_BEFORE = require('@/assets/images/eggs-before.png');
const EGG_AFTER = require('@/assets/images/eggs-after.png');

type Props = {
  onHatchComplete: () => void;
};

// カラフルな紙吹雪パーティクル
const CONFETTI_COLORS = [
  '#FF6B6B',
  '#4ECDC4',
  '#FFE66D',
  '#95E1D3',
  '#FF8E53',
  '#A8E6CF',
  '#DDA0DD',
  '#87CEEB',
];

type ConfettiPieceProps = {
  color: string;
  startX: number;
  startY: number;
  delay: number;
};

function ConfettiPiece({ color, startX, startY, delay }: ConfettiPieceProps) {
  const translateY = useSharedValue(0);
  const translateX = useSharedValue(0);
  const rotate = useSharedValue(0);
  const opacity = useSharedValue(0);
  const scale = useSharedValue(0);

  const targetX = (Math.random() - 0.5) * SCREEN_WIDTH * 1.5;
  const targetY = Math.random() * SCREEN_HEIGHT * 0.8 + 100;

  useEffect(() => {
    opacity.value = withDelay(delay, withTiming(1, { duration: 50 }));
    scale.value = withDelay(delay, withSpring(1, { damping: 8 }));
    translateX.value = withDelay(
      delay,
      withTiming(targetX, { duration: 1500, easing: Easing.out(Easing.cubic) })
    );
    translateY.value = withDelay(
      delay,
      withTiming(targetY, { duration: 1500, easing: Easing.in(Easing.quad) })
    );
    rotate.value = withDelay(
      delay,
      withTiming(360 * (Math.random() > 0.5 ? 3 : -3), { duration: 1500 })
    );
    opacity.value = withDelay(delay + 1000, withTiming(0, { duration: 500 }));
  }, []);

  const style = useAnimatedStyle(() => ({
    transform: [
      { translateX: translateX.value },
      { translateY: translateY.value },
      { rotate: `${rotate.value}deg` },
      { scale: scale.value },
    ],
    opacity: opacity.value,
  }));

  const size = 8 + Math.random() * 8;
  const isSquare = Math.random() > 0.5;

  return (
    <Animated.View
      style={[
        styles.confettiPiece,
        style,
        {
          left: startX,
          top: startY,
          width: size,
          height: isSquare ? size : size * 2,
          backgroundColor: color,
          borderRadius: isSquare ? 2 : size / 2,
        },
      ]}
    />
  );
}

export function EggHatch({ onHatchComplete }: Props) {
  const [showCracked, setShowCracked] = useState(false);
  const [showConfetti, setShowConfetti] = useState(false);

  // 卵のアニメーション値
  const eggRotate = useSharedValue(0);
  const eggScale = useSharedValue(1);
  const eggOpacity = useSharedValue(1);
  const flashOpacity = useSharedValue(0);

  useEffect(() => {
    // Phase 1: 揺れる
    setTimeout(() => {
      eggRotate.value = withSequence(
        withTiming(-10, { duration: 100 }),
        withTiming(10, { duration: 100 }),
        withTiming(-12, { duration: 80 }),
        withTiming(12, { duration: 80 }),
        withTiming(-15, { duration: 50 }),
        withTiming(15, { duration: 50 }),
        withTiming(-18, { duration: 40 }),
        withTiming(18, { duration: 40 }),
        withTiming(0, { duration: 40 })
      );
    }, 600);

    // ひび割れ画像に切り替え
    setTimeout(() => {
      setShowCracked(true);
    }, 1800);

    // Phase 2: 爆発（1.8秒後）
    setTimeout(() => {
      setShowConfetti(true);

      // フラッシュ
      flashOpacity.value = withSequence(withTiming(1, { duration: 500 }));

      // 卵が消える
      eggScale.value = withTiming(1.3, { duration: 500 });
    }, 2300);

    // Phase 3: 完了（3.5秒後）
    setTimeout(() => {
      runOnJS(onHatchComplete)();
    }, 3000);
  }, []);

  const eggStyle = useAnimatedStyle(() => ({
    transform: [{ rotate: `${eggRotate.value}deg` }, { scale: eggScale.value }],
    opacity: eggOpacity.value,
  }));

  const flashStyle = useAnimatedStyle(() => ({
    opacity: flashOpacity.value,
  }));

  const eggCenterX = SCREEN_WIDTH / 2;
  const eggCenterY = SCREEN_HEIGHT / 2 - 50;

  return (
    <View style={styles.container}>
      {/* 背景 */}
      <View style={styles.background} />

      {/* 卵 */}
      <Animated.View style={[styles.eggContainer, eggStyle]}>
        <Image
          source={showCracked ? EGG_AFTER : EGG_BEFORE}
          style={styles.eggImage}
          contentFit="contain"
        />
      </Animated.View>

      {/* フラッシュエフェクト */}
      <Animated.View style={[styles.flash, flashStyle]} />

      {/* 紙吹雪 */}
      {showConfetti && (
        <View style={styles.confettiContainer}>
          {Array.from({ length: 50 }).map((_, i) => (
            <ConfettiPiece
              key={i}
              color={CONFETTI_COLORS[i % CONFETTI_COLORS.length]}
              startX={eggCenterX}
              startY={eggCenterY}
              delay={i * 20}
            />
          ))}
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    ...StyleSheet.absoluteFillObject,
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 1000,
  },
  background: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: '#1a1a2e',
  },
  eggContainer: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  eggImage: {
    width: 200,
    height: 250,
  },
  flash: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: '#fff',
  },
  confettiContainer: {
    ...StyleSheet.absoluteFillObject,
  },
  confettiPiece: {
    position: 'absolute',
  },
});
