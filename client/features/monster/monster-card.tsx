import { StyleSheet, View, Text, Pressable, Dimensions } from 'react-native';
import { Image } from 'expo-image';
import { router } from 'expo-router';
import { Colors } from '@/constants/theme';
import { TRASH_TYPE_COLORS } from '@/constants/trash';
import type { TrashType } from './types';

const { width: SCREEN_WIDTH } = Dimensions.get('window');
const CARD_MARGIN = 8;
const CARD_WIDTH = (SCREEN_WIDTH - 12 * 2 - CARD_MARGIN * 4) / 2;

type MonsterCard = {
  id: string;
  name: string;
  monsterImage: string;
  trashType: TrashType;
};

type Props = {
  monster: MonsterCard;
};

export const MonsterCard = ({ monster }: Props) => {
  const color = TRASH_TYPE_COLORS[monster.trashType];
  return (
    <Pressable
      style={[styles.card, { shadowColor: color }]}
      onPress={() => router.push(`/monsters/${monster.id}`)}
    >
      <Image
        source={{ uri: monster.monsterImage }}
        style={[styles.image, { borderColor: color }]}
        contentFit="cover"
      />
      <View style={[styles.nameContainer, { backgroundColor: color }]}>
        <Text style={styles.name} numberOfLines={1}>
          {monster.name}
        </Text>
      </View>
    </Pressable>
  );
};

const styles = StyleSheet.create({
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
    marginBottom: 12,
  },
  nameContainer: {
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
