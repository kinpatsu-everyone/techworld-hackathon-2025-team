import { useState } from 'react';
import { StyleSheet, View, Text, ScrollView, Pressable } from 'react-native';
import { Image } from 'expo-image';
import { router } from 'expo-router';
import type { Monster, TrashType } from './types';
import { EggHatch } from './components/egg-hatch';

type Props = {
  monster: Monster;
  isFromRegister?: boolean;
};

const TRASH_TYPE_ICONS: Record<TrashType, string> = {
  燃えるゴミ: '🔥',
  燃えないゴミ: '🪨',
  プラスチック: '♻️',
  '缶・ビン': '🥫',
  ペットボトル: '🧴',
  紙類: '📄',
  その他: '📦',
};

export function MonsterDetailPresentational({
  monster,
  isFromRegister = false,
}: Props) {
  const [showMonster, setShowMonster] = useState(true);
  const [isHatching, setIsHatching] = useState(isFromRegister);

  // 卵が割れている間は卵アニメーションを表示
  if (isHatching) {
    return <EggHatch onHatchComplete={() => setIsHatching(false)} />;
  }

  return (
    <View style={styles.wrapper}>
      <ScrollView
        style={styles.container}
        contentContainerStyle={styles.contentContainer}
      >
        {/* 1. トグル */}
        <View style={styles.toggleContainer}>
          <Pressable
            style={[
              styles.toggleButton,
              !showMonster && styles.toggleButtonActive,
            ]}
            onPress={() => setShowMonster(false)}
          >
            <Text
              style={[
                styles.toggleText,
                !showMonster && styles.toggleTextActive,
              ]}
            >
              ゴミ箱
            </Text>
          </Pressable>
          <Pressable
            style={[
              styles.toggleButton,
              showMonster && styles.toggleButtonActive,
            ]}
            onPress={() => setShowMonster(true)}
          >
            <Text
              style={[
                styles.toggleText,
                showMonster && styles.toggleTextActive,
              ]}
            >
              ゴミスター
            </Text>
          </Pressable>
        </View>

        {/* 2. 画像 */}
        <View style={styles.imageContainer}>
          <Image
            source={{
              uri: showMonster ? monster.monsterImage : monster.trashImage,
            }}
            style={styles.image}
            contentFit="cover"
          />
        </View>

        {/* 3. モンスター名 */}
        <Text style={styles.monsterName}>{monster.name}</Text>

        {/* 4. ゴミ種別とアイコン */}
        <View style={styles.trashTypesContainer}>
          {monster.trashTypes.map((type) => (
            <View key={type} style={styles.trashTypeTag}>
              <Text style={styles.trashTypeIcon}>
                {TRASH_TYPE_ICONS[type] || '📦'}
              </Text>
              <Text style={styles.trashTypeText}>{type}</Text>
            </View>
          ))}
        </View>

        {/* 5. 詳細 */}
        <View style={styles.descriptionContainer}>
          <Text style={styles.descriptionLabel}>📍 場所の詳細</Text>
          <Text style={styles.descriptionText}>{monster.description}</Text>
        </View>

        {/* 6. モンスター一覧画面への動線（登録直後のみ表示） */}
        {isFromRegister && (
          <Pressable
            style={styles.listButton}
            onPress={() => router.replace('/monsters')}
          >
            <Text style={styles.listButtonText}>ゴミスター一覧を見る</Text>
          </Pressable>
        )}
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  wrapper: {
    flex: 1,
  },
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  contentContainer: {
    padding: 20,
    alignItems: 'center',
  },
  toggleContainer: {
    flexDirection: 'row',
    backgroundColor: '#e0e0e0',
    borderRadius: 12,
    padding: 4,
    marginBottom: 20,
  },
  toggleButton: {
    paddingVertical: 10,
    paddingHorizontal: 24,
    borderRadius: 10,
  },
  toggleButtonActive: {
    backgroundColor: '#fff',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  toggleText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#888',
  },
  toggleTextActive: {
    color: '#333',
  },
  imageContainer: {
    width: '100%',
    aspectRatio: 1,
    borderRadius: 20,
    overflow: 'hidden',
    backgroundColor: '#fff',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
    marginBottom: 20,
  },
  image: {
    width: '100%',
    height: '100%',
  },
  monsterName: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 16,
  },
  trashTypesContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'center',
    gap: 8,
    marginBottom: 20,
  },
  trashTypeTag: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    paddingVertical: 8,
    paddingHorizontal: 14,
    borderRadius: 20,
    gap: 6,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 2,
    elevation: 1,
  },
  trashTypeIcon: {
    fontSize: 16,
  },
  trashTypeText: {
    fontSize: 14,
    fontWeight: '500',
    color: '#555',
  },
  descriptionContainer: {
    width: '100%',
    backgroundColor: '#fff',
    padding: 16,
    borderRadius: 16,
    marginBottom: 24,
  },
  descriptionLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: '#888',
    marginBottom: 8,
  },
  descriptionText: {
    fontSize: 16,
    color: '#333',
    lineHeight: 24,
  },
  listButton: {
    width: '100%',
    backgroundColor: '#007AFF',
    paddingVertical: 16,
    borderRadius: 14,
    alignItems: 'center',
  },
  listButtonText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#fff',
  },
});
