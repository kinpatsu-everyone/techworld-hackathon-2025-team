import { StyleSheet, View, Text } from 'react-native';

export default function BadgeTabScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.emoji}>🏅</Text>
      <Text style={styles.title}>COMING SOON</Text>
      <Text style={styles.subtitle}>バッジ機能は準備中です</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f5f5f5',
    padding: 40,
  },
  emoji: {
    fontSize: 120,
    marginBottom: 24,
  },
  title: {
    fontSize: 48,
    fontWeight: '900',
    color: '#333',
    letterSpacing: 4,
    marginBottom: 16,
  },
  subtitle: {
    fontSize: 18,
    color: '#888',
  },
});
