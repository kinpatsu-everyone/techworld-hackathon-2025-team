import { useLocalSearchParams } from 'expo-router';
import { MonsterDetailContainer } from '@/features/monster/container';

export default function MonsterDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();

  return <MonsterDetailContainer monsterId={id ?? ''} />;
}
