import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { MovieProvider } from './contexts/MovieContext';
import Navbar from './components/Navbar';
import Homepage from './pages/Homepage';
import Login from './pages/Login';
import Register from './pages/Register';
import Movies from './pages/Movies';
import AddMovie from './pages/AddMovie';
import MovieDetails from './pages/MovieDetails';

function App() {
  return (
    <AuthProvider>
      <MovieProvider>
        <Router>
          <Navbar />
          <main className="min-h-screen bg-background">
            <Routes>
              <Route path="/" element={<Homepage />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/movies" element={<Movies />} />
              <Route path="/movies/add" element={<AddMovie />} />
              <Route path="/movies/:id" element={<MovieDetails />} />
            </Routes>
          </main>
        </Router>
      </MovieProvider>
    </AuthProvider>
  );
}

export default App;
