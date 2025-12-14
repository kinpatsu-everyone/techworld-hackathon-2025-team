import { useEffect } from 'react';
import { StyleSheet, View, Text, Dimensions } from 'react-native';
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withRepeat,
  withTiming,
  withSequence,
  withDelay,
  Easing,
  interpolate,
  interpolateColor,
} from 'react-native-reanimated';

const { width: SCREEN_WIDTH, height: SCREEN_HEIGHT } = Dimensions.get('window');
const CENTER_SIZE = 80;

const COLORS = {
  purple: '#A855F7',
  pink: '#EC4899',
  cyan: '#06B6D4',
  orange: '#F97316',
  yellow: '#FACC15',
};

const PARTICLE_COLORS = [
  COLORS.purple,
  COLORS.pink,
  COLORS.cyan,
  COLORS.orange,
  COLORS.yellow,
];

const MESSAGES = [
  'エネルギーを集めています...',
  'モンスターの魂を呼び起こしています...',
  'パワーを注入中...',
  'もうすぐ誕生します！',
];

type EnergyParticleProps = {
  index: number;
  total: number;
  delay: number;
  color: string;
};

const EnergyParticle = ({ index, total, delay, color }: EnergyParticleProps) => {
  const progress = useSharedValue(0);
  const angle = (index / total) * 2 * Math.PI;
  const startRadius = Math.max(SCREEN_WIDTH, SCREEN_HEIGHT) * 0.4;

  useEffect(() => {
    progress.value = withDelay(
      delay,
      withRepeat(
        withTiming(1, { duration: 2000, easing: Easing.inOut(Easing.ease) }),
        -1,
        false
      )
    );
  }, []);

  const animatedStyle = useAnimatedStyle(() => {
    const radius = interpolate(progress.value, [0, 1], [startRadius, 0]);
    const x = Math.cos(angle) * radius;
    const y = Math.sin(angle) * radius;
    const scale = interpolate(progress.value, [0, 0.5, 1], [0.3, 1, 0]);
    const opacity = interpolate(progress.value, [0, 0.3, 0.8, 1], [0, 1, 1, 0]);

    return {
      transform: [{ translateX: x }, { translateY: y }, { scale }],
      opacity,
    };
  });

  return (
    <Animated.View
      style={[
        styles.particle,
        { backgroundColor: color, shadowColor: color },
        animatedStyle,
      ]}
    />
  );
};

type OrbitingSparkProps = {
  index: number;
  total: number;
  radius: number;
  duration: number;
  size: number;
  color: string;
};

const OrbitingSpark = ({
  index,
  total,
  radius,
  duration,
  size,
  color,
}: OrbitingSparkProps) => {
  const rotation = useSharedValue(0);
  const initialAngle = (index / total) * 2 * Math.PI;

  useEffect(() => {
    rotation.value = withRepeat(
      withTiming(2 * Math.PI, { duration, easing: Easing.linear }),
      -1,
      false
    );
  }, []);

  const animatedStyle = useAnimatedStyle(() => {
    const currentAngle = rotation.value + initialAngle;
    const x = Math.cos(currentAngle) * radius;
    const y = Math.sin(currentAngle) * radius;

    return {
      transform: [{ translateX: x }, { translateY: y }],
    };
  });

  return (
    <Animated.View
      style={[
        styles.spark,
        {
          width: size,
          height: size,
          borderRadius: size / 2,
          backgroundColor: color,
          shadowColor: color,
        },
        animatedStyle,
      ]}
    />
  );
};

const CentralCore = () => {
  const pulse = useSharedValue(0);
  const colorProgress = useSharedValue(0);

  useEffect(() => {
    pulse.value = withRepeat(
      withSequence(
        withTiming(1, { duration: 800, easing: Easing.inOut(Easing.ease) }),
        withTiming(0, { duration: 800, easing: Easing.inOut(Easing.ease) })
      ),
      -1,
      false
    );

    colorProgress.value = withRepeat(
      withTiming(4, { duration: 4000, easing: Easing.linear }),
      -1,
      false
    );
  }, []);

  const coreStyle = useAnimatedStyle(() => {
    const scale = interpolate(pulse.value, [0, 1], [0.8, 1.2]);

    return {
      transform: [{ scale }],
    };
  });

  const glowStyle = useAnimatedStyle(() => {
    const scale = interpolate(pulse.value, [0, 1], [1, 1.5]);
    const opacity = interpolate(pulse.value, [0, 1], [0.6, 0.2]);

    const bgColor = interpolateColor(
      colorProgress.value,
      [0, 1, 2, 3, 4],
      [COLORS.purple, COLORS.pink, COLORS.cyan, COLORS.orange, COLORS.purple]
    );

    return {
      transform: [{ scale }],
      opacity,
      backgroundColor: bgColor,
    };
  });

  const innerGlowStyle = useAnimatedStyle(() => {
    const bgColor = interpolateColor(
      colorProgress.value,
      [0, 1, 2, 3, 4],
      [COLORS.pink, COLORS.cyan, COLORS.orange, COLORS.purple, COLORS.pink]
    );

    return {
      backgroundColor: bgColor,
    };
  });

  return (
    <View style={styles.coreContainer}>
      <Animated.View style={[styles.glow, glowStyle]} />
      <Animated.View style={[styles.core, coreStyle]}>
        <Animated.View style={[styles.innerCore, innerGlowStyle]} />
      </Animated.View>
    </View>
  );
};

const AnimatedMessage = () => {
  const opacity = useSharedValue(1);
  const messageIndex = useSharedValue(0);

  useEffect(() => {
    const cycleMessage = () => {
      opacity.value = withSequence(
        withTiming(0, { duration: 300 }),
        withTiming(1, { duration: 300 })
      );

      setTimeout(() => {
        messageIndex.value = (messageIndex.value + 1) % MESSAGES.length;
      }, 300);
    };

    const interval = setInterval(cycleMessage, 3000);
    return () => clearInterval(interval);
  }, []);

  const textStyle = useAnimatedStyle(() => ({
    opacity: opacity.value,
  }));

  return (
    <Animated.View style={[styles.messageContainer, textStyle]}>
      <Text style={styles.messageText}>
        {MESSAGES[Math.floor(messageIndex.value)]}
      </Text>
    </Animated.View>
  );
};

export const MonsterGenerationLoading = () => {
  const bgOpacity = useSharedValue(0);

  useEffect(() => {
    bgOpacity.value = withTiming(1, { duration: 500 });
  }, []);

  const containerStyle = useAnimatedStyle(() => ({
    opacity: bgOpacity.value,
  }));

  const particles = Array.from({ length: 12 }, (_, i) => ({
    index: i,
    delay: i * 150,
    color: PARTICLE_COLORS[i % PARTICLE_COLORS.length],
  }));

  const innerSparks = Array.from({ length: 8 }, (_, i) => ({
    index: i,
    radius: 60,
    duration: 2500,
    size: 8,
    color: PARTICLE_COLORS[i % PARTICLE_COLORS.length],
  }));

  const outerSparks = Array.from({ length: 6 }, (_, i) => ({
    index: i,
    radius: 100,
    duration: 4000,
    size: 6,
    color: PARTICLE_COLORS[(i + 2) % PARTICLE_COLORS.length],
  }));

  return (
    <Animated.View style={[styles.container, containerStyle]}>
      <View style={styles.animationArea}>
        {particles.map((p) => (
          <EnergyParticle
            key={`particle-${p.index}`}
            index={p.index}
            total={particles.length}
            delay={p.delay}
            color={p.color}
          />
        ))}

        {innerSparks.map((s) => (
          <OrbitingSpark
            key={`inner-spark-${s.index}`}
            index={s.index}
            total={innerSparks.length}
            radius={s.radius}
            duration={s.duration}
            size={s.size}
            color={s.color}
          />
        ))}

        {outerSparks.map((s) => (
          <OrbitingSpark
            key={`outer-spark-${s.index}`}
            index={s.index}
            total={outerSparks.length}
            radius={s.radius}
            duration={s.duration}
            size={s.size}
            color={s.color}
          />
        ))}

        <CentralCore />
      </View>

      <AnimatedMessage />
    </Animated.View>
  );
};

const styles = StyleSheet.create({
  container: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'rgba(0, 0, 0, 0.9)',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 999,
  },
  animationArea: {
    width: SCREEN_WIDTH,
    height: SCREEN_WIDTH,
    justifyContent: 'center',
    alignItems: 'center',
  },
  particle: {
    position: 'absolute',
    width: 16,
    height: 16,
    borderRadius: 8,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 1,
    shadowRadius: 10,
    elevation: 10,
  },
  spark: {
    position: 'absolute',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.8,
    shadowRadius: 6,
    elevation: 6,
  },
  coreContainer: {
    position: 'absolute',
    width: CENTER_SIZE,
    height: CENTER_SIZE,
    justifyContent: 'center',
    alignItems: 'center',
  },
  glow: {
    position: 'absolute',
    width: CENTER_SIZE * 2,
    height: CENTER_SIZE * 2,
    borderRadius: CENTER_SIZE,
    backgroundColor: COLORS.purple,
  },
  core: {
    width: CENTER_SIZE,
    height: CENTER_SIZE,
    borderRadius: CENTER_SIZE / 2,
    backgroundColor: '#fff',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#fff',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 1,
    shadowRadius: 20,
    elevation: 20,
  },
  innerCore: {
    width: CENTER_SIZE * 0.6,
    height: CENTER_SIZE * 0.6,
    borderRadius: (CENTER_SIZE * 0.6) / 2,
  },
  messageContainer: {
    position: 'absolute',
    bottom: 120,
    paddingHorizontal: 30,
  },
  messageText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: '600',
    textAlign: 'center',
    textShadowColor: COLORS.purple,
    textShadowOffset: { width: 0, height: 0 },
    textShadowRadius: 10,
  },
});
