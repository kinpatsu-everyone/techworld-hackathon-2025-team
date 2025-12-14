import { View, StyleSheet, ActivityIndicator, Text } from 'react-native';
import { MonsterDetailPresentational } from './presentational';
import { useApi } from '@/hooks/use-api';
import { MonsterItem } from '@/lib/client';
import type { Monster, TrashType } from './types';

type Props = {
  monsterId: string;
  isFromRegister?: boolean;
};

// APIのtrash_categoryをTrashType配列に変換
function convertTrashCategory(trashCategory: string): TrashType[] {
  const categoryMap: Record<string, TrashType> = {
    燃えるゴミ: '燃えるゴミ',
    不燃ごみ: '燃えないゴミ',
    燃えないゴミ: '燃えないゴミ',
    プラスチック: 'プラスチック',
    缶: '缶・ビン',
    瓶: '缶・ビン',
    '缶・ビン': '缶・ビン',
    ペットボトル: 'ペットボトル',
    紙類: '紙類',
    指定なし: 'その他',
  };
  return [categoryMap[trashCategory] || 'その他'];
}

// APIのMonsterItemをフロントエンドのMonster型に変換
function convertToMonster(item: MonsterItem): Monster {
  return {
    id: item.id,
    name: item.nickname,
    trashTypes: convertTrashCategory(item.trash_category),
    latitude: item.latitude,
    longitude: item.longitude,
    description: '', // APIに対応するフィールドがないため空
    trashImage: item.image_url, // 現状はmonsterImageと同じ
    monsterImage: item.image_url,
  };
}

export function MonsterDetailContainer({ monsterId }: Props) {
  const { data, isLoading, error } = useApi('/monster/v1/GetMonster', {
    id: monsterId,
  });

  const mockMonster = {
    id: '1',
    nickname: 'ゴミスター1',
    trash_category: '燃えるゴミ',
    image_url: 'https://avatars.githubusercontent.com/u/248258447?s=200&v=4',
    latitude: 35.681236,
    longitude: 139.767125,
  };

  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  //   if (error || !data) {
  //     return (
  //       <View style={styles.errorContainer}>
  //         <Text style={styles.errorText}>
  //           モンスターの読み込みに失敗しました
  //         </Text>
  //       </View>
  //     );
  //   }

  const monster = convertToMonster(mockMonster);

  return (
    <MonsterDetailPresentational monster={monster} isFromRegister={false} />
  );
}

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  errorText: {
    fontSize: 16,
    color: '#666',
    textAlign: 'center',
  },
});
