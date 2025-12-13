import { StyleSheet, Text, View } from 'react-native';

export default function TrashRegisterScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>ゴミ箱登録</Text>
      <Text style={styles.subtitle}>新しいゴミ箱の位置を登録します</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f5f5f5',
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 10,
    color: '#333',
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    textAlign: 'center',
  },
});
