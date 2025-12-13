export type TrashType =
  | '燃えるゴミ'
  | '燃えないゴミ'
  | 'プラスチック'
  | '缶・ビン'
  | 'ペットボトル'
  | '紙類'
  | 'その他';

export type Monster = {
  id: string;
  name: string;
  trashTypes: TrashType[];
  latitude: number;
  longitude: number;
  description: string;
  trashImage: string;
  monsterImage: string;
};
