import { useState, useEffect } from "react";
import useAxiosPrivate from "../../hooks/useAxiosPrivate";
import Movies from "../movies/movies";
import Spinner from "../spinner/Spinner";

const Recommended = () => {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState();
  const axiosPrivate = useAxiosPrivate();

  useEffect(() => {
    const fetchRecommendedMovies = async () => {
      setLoading(true);
      setMessage("Loading recommended movies...");
      try {
        const response = await axiosPrivate.get("/recommendedMovies");
        setMovies(response.data);
      } catch (error) {
        console.error("Error fetching recommended movies:", error);
        setMessage("Failed to load recommended movies.");
      } finally {
        setLoading(false);
      }
    };
    fetchRecommendedMovies();
  }, []);
  return (
    <>{loading ? <Spinner /> : <Movies movies={movies} message={message} />}</>
  );
};
export default Recommended;
