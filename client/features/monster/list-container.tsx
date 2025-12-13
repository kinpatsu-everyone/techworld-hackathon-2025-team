import { MonsterListPresentational } from './list-presentational';

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

export const MonsterListContainer = () => {
  // TODO: useQueryなどでAPIからデータ取得
  const monsters = MOCK_MONSTERS;

  return <MonsterListPresentational monsters={monsters} />;
};
