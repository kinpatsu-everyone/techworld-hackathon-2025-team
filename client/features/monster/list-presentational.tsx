import {
  StyleSheet,
  View,
  Text,
  Pressable,
  FlatList,
  Dimensions,
} from 'react-native';
import { Image } from 'expo-image';
import { router } from 'expo-router';
import { MonsterItem } from '@/lib/client';
const { width: SCREEN_WIDTH } = Dimensions.get('window');
const CARD_MARGIN = 8;
const CARD_WIDTH = (SCREEN_WIDTH - 12 * 2 - CARD_MARGIN * 4) / 2;

type Props = {
  monsters: MonsterItem[];
};

export const MonsterListPresentational = ({ monsters }: Props) => {
  return (
    <View style={styles.container}>
      <FlatList
        data={monsters}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContent}
        numColumns={2}
        renderItem={({ item }) => (
          <Pressable
            style={styles.card}
            onPress={() => router.push(`/monsters/${item.id}`)}
          >
            <Image
              source={{ uri: item.image_url }}
              style={styles.image}
              contentFit="cover"
            />
            <Text style={styles.name} numberOfLines={1}>
              {item.nickname}
            </Text>
          </Pressable>
        )}
      />
    </View>
  );
};

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
