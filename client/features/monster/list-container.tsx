import { useApi } from '@/hooks/use-api';
import { useEffect, useState } from 'react';
import { MonsterListPresentational } from './list-presentational';
import { MonsterItem } from '@/lib/client';
import { View, StyleSheet } from 'react-native';
import { ActivityIndicator } from 'react-native';

export const MonsterListContainer = () => {
  const { data, isLoading } = useApi('/monster/v1/GetMonsters', {});
  const [monsters, setMonsters] = useState<MonsterItem[]>([]);

  useEffect(() => {
    if (data) {
      setMonsters(data.monsters);
    }
  }, [data]);

  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  return <MonsterListPresentational monsters={monsters} />;
};

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
});
