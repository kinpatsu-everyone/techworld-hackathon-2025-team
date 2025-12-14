import { MonsterListContainer } from '@/features/monster';
import { SafeAreaView } from 'react-native-safe-area-context';

export default function MonstersTabScreen() {
  return (
    <SafeAreaView style={{ flex: 1 }}>
      <MonsterListContainer />
    </SafeAreaView>
  );
}
