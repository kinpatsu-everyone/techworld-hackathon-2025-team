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
          image:
            'https://images.unsplash.com/photo-1532996122724-e3c354a0b15b?w=400&h=300&fit=crop&crop=center',
        },
        {
          id: 2,
          latitude: latitude - 0.002,
          longitude: longitude + 0.003,
          title: 'ゴミ箱 #2',
          image:
            'https://images.unsplash.com/photo-1532996122724-e3c354a0b15b?w=400&h=300&fit=crop&crop=center',
        },
        {
          id: 3,
          latitude: latitude + 0.003,
          longitude: longitude - 0.002,
          title: 'ゴミ箱 #3',
          image:
            'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400&h=300&fit=crop&crop=center',
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
