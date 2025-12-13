import { MonsterDetailPresentational } from './presentational';
import type { Monster } from './types';

type Props = {
  monsterId: string;
  isFromRegister?: boolean;
};

// TODO: 実際のAPI連携時に置き換え
const MOCK_MONSTER: Monster = {
  id: '1',
  name: 'バーニングゴミスター',
  trashTypes: ['燃えるゴミ', 'プラスチック'],
  latitude: 35.6762,
  longitude: 139.6503,
  description: '築地松竹ビル 5階 エレベーターホール横',
  trashImage: 'https://picsum.photos/400',
  monsterImage: 'https://picsum.photos/401',
};

export function MonsterDetailContainer({
  monsterId,
  isFromRegister = false,
}: Props) {
  // TODO: useQueryなどでAPIからデータ取得
  const monster = MOCK_MONSTER;

  return (
    <MonsterDetailPresentational
      monster={monster}
      isFromRegister={isFromRegister}
    />
  );
}
