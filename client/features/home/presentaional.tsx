import { StyleSheet } from "react-native";
import MapView from "react-native-maps";

export const HomePresentational = () => {
  return <MapView style={styles.map} />;
};

const styles = StyleSheet.create({
  titleContainer: {
    flexDirection: "row",
    alignItems: "center",
    gap: 8,
  },
  stepContainer: {
    gap: 8,
    marginBottom: 8,
  },
  reactLogo: {
    height: 178,
    width: 290,
    bottom: 0,
    left: 0,
    position: "absolute",
  },
  map: {
    width: "100%",
    height: "100%",
    marginBottom: 16,
  },
});
