import { StyleSheet, View, Pressable } from 'react-native';

type Props = {
  onPress: () => void;
  color?: string;
};

export function CaptureButton({ onPress, color = '#34C759' }: Props) {
  return (
    <Pressable
      style={({ pressed }) => [
        styles.button,
        { borderColor: color },
        pressed && styles.buttonPressed,
      ]}
      onPress={onPress}
    >
      <View style={[styles.inner, { backgroundColor: color }]} />
    </Pressable>
  );
}

const styles = StyleSheet.create({
  button: {
    marginTop: 30,
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: '#fff',
    borderWidth: 5,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.2,
    shadowRadius: 4,
    elevation: 4,
  },
  buttonPressed: {
    transform: [{ scale: 0.95 }],
    opacity: 0.9,
  },
  inner: {
    width: 60,
    height: 60,
    borderRadius: 30,
  },
});
