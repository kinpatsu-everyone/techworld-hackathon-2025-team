import React from 'react';
import { StyleSheet, View, Text, TouchableOpacity } from 'react-native';

export type TrashType = 'all' | 'burnable' | 'non-burnable' | 'plastic' | 'cans-bottles' | 'paper' | 'unknown';

interface FilterOption {
  id: TrashType;
  label: string;
  color: string;
}

const filterOptions: FilterOption[] = [
  { id: 'all', label: 'すべて', color: '#808080' },
  { id: 'burnable', label: '可燃ゴミ', color: '#FF9500' },
  { id: 'non-burnable', label: '不燃ゴミ', color: '#007AFF' },
  { id: 'plastic', label: 'ペットボトル', color: '#34C759' },
  { id: 'cans-bottles', label: '缶・ビン', color: '#FFCC00' },
  { id: 'paper', label: '紙類', color: '#AF52DE' },
  { id: 'unknown', label: '不明', color: '#C7C7CC' },
];

interface FilterMenuProps {
  selectedFilters: TrashType[];
  onFilterChange: (filters: TrashType[]) => void;
}

export const FilterMenu: React.FC<FilterMenuProps> = ({
  selectedFilters,
  onFilterChange,
}) => {
  const handleOptionPress = (optionId: TrashType) => {
    if (optionId === 'all') {
      onFilterChange(['all']);
    } else {
      let newFilters = selectedFilters.filter((f) => f !== 'all');

      if (newFilters.includes(optionId)) {
        newFilters = newFilters.filter((f) => f !== optionId);
        if (newFilters.length === 0) {
          newFilters = ['all'];
        }
      } else {
        newFilters = [...newFilters, optionId];
      }

      onFilterChange(newFilters);
    }
  };

  const isSelected = (optionId: TrashType) => {
    return selectedFilters.includes(optionId);
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>ゴミの種類で絞り込み</Text>
      {filterOptions.map((option) => (
        <TouchableOpacity
          key={option.id}
          style={styles.optionRow}
          onPress={() => handleOptionPress(option.id)}
        >
          <View style={[styles.colorDot, { backgroundColor: option.color }]} />
          <Text style={styles.optionLabel}>{option.label}</Text>
          <Text style={[styles.checkmark, !isSelected(option.id) && styles.checkmarkHidden]}>✓</Text>
        </TouchableOpacity>
      ))}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: 'white',
    borderRadius: 12,
    paddingVertical: 8,
    marginTop: 8,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
    elevation: 5,
    minWidth: 220,
  },
  title: {
    fontSize: 14,
    color: '#666',
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#E5E5E5',
  },
  optionRow: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 14,
  },
  colorDot: {
    width: 20,
    height: 20,
    borderRadius: 10,
    marginRight: 12,
  },
  optionLabel: {
    fontSize: 16,
    color: '#333',
    flex: 1,
  },
  checkmark: {
    fontSize: 18,
    color: '#34C759',
    fontWeight: '600',
    width: 20,
  },
  checkmarkHidden: {
    opacity: 0,
  },
});
