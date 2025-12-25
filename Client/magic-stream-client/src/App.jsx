import Home from "./components/home/Home";
import "./App.css";
import Header from "./components/header/Header";
import { Route, Routes, useNavigate } from "react-router-dom";
import Register from "./components/register/Register";
import Login from "./components/login/Login";
import axiosClient from "./api/axiosConfig";
import useAuth from "./hooks/useAuth";
import Layout from "./components/Layout";
import RequiredAuth from "./components/RequiredAuth";
import Recommended from "./components/recommended/Recommended";
import Review from "./components/review/Review";
import StreamMovie from "./components/stream/StreamMovie";

function App() {
  const navigate = useNavigate();
  const { auth, setAuth } = useAuth();

  const updateMovieReview = (imdb_id) => {
    navigate(`/review/${imdb_id}`);
  };
  const handleLogout = async () => {
    try {
      await axiosClient.post("/user/logout", {
        user_id: auth.user_id,
      });
      setAuth(null);
      navigate("/login", { replace: true });
    } catch (err) {
      console.error("Error logging out:", err);
    }
  };

  return (
    <>
      <Header handleLogout={handleLogout} />
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route
            path="/"
            element={<Home updateMovieReview={updateMovieReview} />}
          />
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />
          <Route element={<RequiredAuth />}>
            <Route path="/recommended" element={<Recommended />} />
            <Route path="/review/:imdb_id" element={<Review />} />
            <Route path="/stream/:yt_id" element={<StreamMovie />} />
          </Route>
        </Route>
      </Routes>
    </>
  );
}

export default App;
