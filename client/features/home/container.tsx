import { useEffect, useState } from 'react';
import { ActivityIndicator, View, StyleSheet } from 'react-native';
import { HomePresentational } from './presentational';
import { useLocation } from './hooks/useLocation';
import { useApi } from '@/hooks/use-api';
import { MonsterItem } from '@/lib/client';

export const HomeContainer = () => {
  const { location, errorMsg } = useLocation();
  const [trashBins, setTrashBins] = useState<MonsterItem[]>([]);

  const { data, isLoading } = useApi('/monster/v1/GetMonsters', {});

  useEffect(() => {
    if (data) {
      setTrashBins(data.monsters);
    }
  }, [data]);

  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  return (
    <HomePresentational
      location={location}
      errorMsg={errorMsg}
      trashBins={trashBins}
    />
  );
};

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
});
