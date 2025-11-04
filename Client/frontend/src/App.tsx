import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./contexts/AuthContext";
import { MovieProvider } from "./contexts/MovieContext";
import { Toaster } from "./components/ui/toaster";
import Navbar from "./components/Navbar";
import Homepage from "./pages/Homepage";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Movies from "./pages/Movies";
import AddMovie from "./pages/AddMovie";
import MovieDetails from "./pages/MovieDetails";
import { PublicLayout, ProtectedLayout } from "./components/Layout";
import Profile from "./pages/Profile";
import ForgotPassword from "./pages/ForgotPassword";
import ResetPassword from "./pages/ResetPassword";
import TMDBMovieDetails from "./pages/TMDBMovieDetails";

function App() {
  return (
    <AuthProvider>
      <MovieProvider>
        <BrowserRouter>
          <Navbar />
          <main className="min-h-screen bg-background">
            <Routes>
              <Route path="/" element={<Homepage />} />
              <Route element={<PublicLayout />}>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/forgot-password" element={<ForgotPassword />} />
                <Route path="/reset-password" element={<ResetPassword />} />
                <Route path="/movies" element={<Movies />} />
                <Route path="/movies/:id" element={<MovieDetails />} />
                <Route path="/tmdb-movies/:id" element={<TMDBMovieDetails />} />
              </Route>
              <Route element={<ProtectedLayout />}>
                <Route path="/movies/add" element={<AddMovie />} />
                <Route path="/profile" element={<Profile />} />
              </Route>
            </Routes>
          </main>
          <Toaster />
        </BrowserRouter>
      </MovieProvider>
    </AuthProvider>
  );
}

export default App;
