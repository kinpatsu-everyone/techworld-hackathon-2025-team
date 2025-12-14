import { StyleSheet, View, FlatList, Dimensions } from 'react-native';
import { Colors } from '@/constants/theme';
import { MonsterCard } from './monster-card';
import { MonsterItem } from '@/lib/client';
import { TrashType } from './types';
import { TRASH_TYPE_COLORS } from '@/constants/trash';
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
          <MonsterCard
            monster={{
              id: item.id,
              name: item.nickname,
              monsterImage: item.image_url,
              trashType: item.trash_category as TrashType,
            }}
          />
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
    borderRadius: 16,
    overflow: 'hidden',
    shadowColor: Colors.light.text,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.5,
    shadowRadius: 4,
    elevation: 3,
  },
  image: {
    width: '100%',
    aspectRatio: 1,
    resizeMode: 'cover',
    borderRadius: 999,
    borderWidth: 4,
    borderColor: Colors.light.text,
    marginBottom: 12,
  },
  nameContainer: {
    backgroundColor: Colors.light.text,
    borderRadius: 8,
    marginTop: -44,
  },
  name: {
    fontSize: 14,
    fontWeight: '600',
    color: Colors.light.background,
    padding: 12,
    textAlign: 'center',
  },
});
