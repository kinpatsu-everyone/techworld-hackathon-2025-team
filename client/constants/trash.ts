import type { TrashType } from '@/features/monster/types';

export const TRASH_TYPE_COLORS: Record<TrashType, string> = {
  燃えるゴミ: '#FF9500',
  燃えないゴミ: '#007AFF',
  プラスチック: '#34C759',
  '缶・ビン': '#FFCC00',
  ペットボトル: '#34C759',
  紙類: '#AF52DE',
  その他: '#C7C7CC',
};