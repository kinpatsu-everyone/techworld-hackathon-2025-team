import { useLocalSearchParams } from 'expo-router';
import { MonsterDetailContainer } from '@/features/monster/container';

export default function MonsterDetailScreen() {
  const { id, fromRegister } = useLocalSearchParams<{
    id: string;
    fromRegister?: string;
  }>();

  return (
    <MonsterDetailContainer
      monsterId={id ?? ''}
      isFromRegister={fromRegister === 'true'}
    />
  );
}
