import { Pressable, View, StyleSheet } from 'react-native';
import { SvgXml } from 'react-native-svg';

interface TrashPlotProps {
  id: string | number;
  x: number;
  y: number;
  isSelected?: boolean;
  onPress: (e: any) => void;
}

const trashIconSvg = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#ffffff"><path d="M280-120q-33 0-56.5-23.5T200-200v-520h-40v-80h200v-40h240v40h200v80h-40v520q0 33-23.5 56.5T680-120H280Zm400-600H280v520h400v-520ZM360-280h80v-360h-80v360Zm160 0h80v-360h-80v360ZM280-720v520-520Z"/></svg>`;

export function TrashPlot({ id, x, y, isSelected, onPress }: TrashPlotProps) {
  return (
    <Pressable
      onPress={onPress}
      className="absolute z-30"
      style={[
        styles.container,
        {
          left: `${x}%`,
          top: `${y}%`,
        },
      ]}
    >
      {({ pressed }) => (
        <View
          className="bg-green-500 rounded-full p-3 shadow-lg border-white relative"
          style={[
            styles.button,
            pressed && styles.buttonPressed,
          ]}
        >
          <SvgXml xml={trashIconSvg} width={24} height={24} />
          {isSelected && (
            <View className="absolute -bottom-1 -right-1 w-3 h-3 bg-blue-500 rounded-full border-2 border-white" />
          )}
        </View>
      )}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: {
    transform: [{ translateX: '-50%' }, { translateY: '-50%' }],
  },
  button: {
    borderWidth: 3,
    borderColor: 'white',
  },
  buttonPressed: {
    transform: [{ scale: 1.1 }],
  },
});
