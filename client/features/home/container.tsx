import { useEffect, useState } from 'react';
import { HomePresentational } from './presentational';
import { useLocation } from './hooks/useLocation';
import { TrashBin } from '@/types/model';

export const HomeContainer = () => {
  const { location, errorMsg } = useLocation();
  const [trashBins, setTrashBins] = useState<TrashBin[]>([]);

  // 初回位置取得時のみゴミ箱位置を設定（後でfetchに変更予定）
  useEffect(() => {
    if (location && trashBins.length === 0) {
      const { latitude, longitude } = location.coords;
      setTrashBins([
        {
          id: 1,
          latitude: latitude + 0.001,
          longitude: longitude + 0.001,
          title: 'ゴミ箱 #1',
        },
        {
          id: 2,
          latitude: latitude - 0.002,
          longitude: longitude + 0.003,
          title: 'ゴミ箱 #2',
        },
        {
          id: 3,
          latitude: latitude + 0.003,
          longitude: longitude - 0.002,
          title: 'ゴミ箱 #3',
        },
      ]);
    }
  }, [location, trashBins.length]);

  return (
    <HomePresentational
      location={location}
      errorMsg={errorMsg}
      trashBins={trashBins}
    />
  );
};
