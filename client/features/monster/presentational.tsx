import { useState } from 'react';
import {
  StyleSheet,
  View,
  Text,
  ScrollView,
  Pressable,
  Alert,
} from 'react-native';
import { Image } from 'expo-image';
import { router } from 'expo-router';
import * as Linking from 'expo-linking';
import { Ionicons } from '@expo/vector-icons';
import type { Monster } from './types';
import { TRASH_TYPE_COLORS } from '@/constants/trash';
import { EggHatch } from './components/egg-hatch';
import { TiltCard } from './components/tilt-card';

type Props = {
  monster: Monster;
  isFromRegister?: boolean;
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

        {/* 2. 画像（3D傾きエフェクト） */}
        <TiltCard>
          <View style={styles.imageContainer}>
            <Image
              source={{
                uri: showMonster ? monster.monsterImage : monster.trashImage,
              }}
              style={styles.image}
              contentFit="cover"
            />
            {/* リボン（登録直後のみ） */}
            {isFromRegister && (
              <View style={styles.ribbon}>
                <Text style={styles.ribbonText}>NEW</Text>
              </View>
            )}
          </View>
        </TiltCard>

        {/* 3. モンスター名 */}
        <Text style={styles.monsterName}>{monster.name}</Text>

        {/* 4. ゴミ種別とアイコン */}
        <View style={styles.trashTypesContainer}>
          {monster.trashTypes.map((type) => (
            <View key={type} style={styles.trashTypeTag}>
              <View
                style={[
                  styles.trashTypeColorDot,
                  { backgroundColor: TRASH_TYPE_COLORS[type] || '#C7C7CC' },
                ]}
              />
              <Text style={styles.trashTypeText}>{type}</Text>
            </View>
          ))}
        </View>

        {/* 5. 詳細 */}
        <View style={styles.descriptionContainer}>
          <View style={styles.descriptionHeader}>
            <Text style={styles.descriptionLabel}>📍 場所の詳細</Text>
            <Pressable
              style={styles.mapButton}
              onPress={() => {
                Alert.alert(
                  'Google Mapで開く',
                  'この場所をGoogle Mapで表示しますか？',
                  [
                    { text: 'キャンセル', style: 'cancel' },
                    {
                      text: '開く',
                      onPress: () => {
                        const url = `https://www.google.com/maps?q=${monster.latitude},${monster.longitude}`;
                        Linking.openURL(url);
                      },
                    },
                  ]
                );
              }}
            >
              <Ionicons name="map-outline" size={20} color="#007AFF" />
              <Text style={styles.mapButtonText}>地図で見る</Text>
            </Pressable>
          </View>
          {monster.description ? (
            <Text style={styles.descriptionText}>{monster.description}</Text>
          ) : (
            <Text style={styles.descriptionText}>
              緯度: {monster.latitude.toFixed(6)}
              {'\n'}
              経度: {monster.longitude.toFixed(6)}
            </Text>
          )}
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
  ribbon: {
    position: 'absolute',
    top: 15,
    right: -45,
    backgroundColor: '#34C759',
    paddingVertical: 14,
    paddingHorizontal: 60,
    transform: [{ rotate: '45deg' }],
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 3 },
    shadowOpacity: 0.25,
    shadowRadius: 4,
    elevation: 4,
  },
  ribbonText: {
    color: '#fff',
    fontSize: 20,
    fontWeight: 'bold',
    textAlign: 'center',
    letterSpacing: 2,
  },
  monsterName: {
    fontSize: 28,
    fontWeight: '900',
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
    paddingVertical: 10,
    paddingHorizontal: 16,
    borderRadius: 20,
    gap: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.08,
    shadowRadius: 3,
    elevation: 2,
  },
  trashTypeColorDot: {
    width: 16,
    height: 16,
    borderRadius: 8,
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
  descriptionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  descriptionLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: '#888',
  },
  mapButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    paddingVertical: 6,
    paddingHorizontal: 12,
    backgroundColor: '#E8F4FD',
    borderRadius: 8,
    borderWidth: 2,
    borderColor: '#007AFF',
  },
  mapButtonText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#007AFF',
  },
  descriptionText: {
    fontSize: 16,
    color: '#333',
    lineHeight: 24,
  },
  listButton: {
    width: '100%',
    backgroundColor: '#34C759',
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
