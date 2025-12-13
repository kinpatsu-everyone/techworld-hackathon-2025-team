import {
  StyleSheet,
  View,
  Text,
  Pressable,
  FlatList,
  Dimensions,
} from 'react-native';

const { width: SCREEN_WIDTH } = Dimensions.get('window');
const CARD_MARGIN = 8;
const CARD_WIDTH = (SCREEN_WIDTH - 12 * 2 - CARD_MARGIN * 4) / 2;
import { Image } from 'expo-image';
import { router } from 'expo-router';

// TODO: APIから取得
const MOCK_MONSTERS = [
  {
    id: '1',
    name: 'バーニングゴミスター',
    monsterImage: 'https://picsum.photos/200',
  },
  {
    id: '2',
    name: 'リサイクルモンスター',
    monsterImage: 'https://picsum.photos/201',
  },
  {
    id: '3',
    name: 'エコファイター',
    monsterImage: 'https://picsum.photos/202',
  },
];

export default function MonstersListScreen() {
  return (
    <View style={styles.container}>
      <FlatList
        data={MOCK_MONSTERS}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContent}
        numColumns={2}
        renderItem={({ item }) => (
          <Pressable
            style={styles.card}
            onPress={() => router.push(`/monsters/${item.id}`)}
          >
            <Image
              source={{ uri: item.monsterImage }}
              style={styles.image}
              contentFit="cover"
            />
            <Text style={styles.name} numberOfLines={1}>
              {item.name}
            </Text>
          </Pressable>
        )}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  listContent: {
    padding: 12,
  },
  card: {
    width: CARD_WIDTH,
    margin: CARD_MARGIN,
    backgroundColor: '#fff',
    borderRadius: 16,
    overflow: 'hidden',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  image: {
    width: '100%',
    aspectRatio: 1,
  },
  name: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333',
    padding: 12,
    textAlign: 'center',
  },
});
