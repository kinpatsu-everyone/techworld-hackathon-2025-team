import { HomePresentational } from "./presentational";
import { useLocation } from "./hooks/useLocation";

export const HomeContainer = () => {
  const { location, errorMsg } = useLocation();

  return <HomePresentational location={location} errorMsg={errorMsg} />;
};
