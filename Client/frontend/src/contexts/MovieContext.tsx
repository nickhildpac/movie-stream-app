import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import type { Movie, CreateMovieInput } from '../types';

interface MovieContextType {
  movies: Movie[];
  addMovie: (movie: CreateMovieInput) => void;
  updateMovie: (id: string, movie: Partial<Movie>) => void;
  deleteMovie: (id: string) => void;
  getMovie: (id: string) => Movie | undefined;
}

const MovieContext = createContext<MovieContextType | undefined>(undefined);

export const useMovies = () => {
  const context = useContext(MovieContext);
  if (!context) {
    throw new Error('useMovies must be used within a MovieProvider');
  }
  return context;
};

interface MovieProviderProps {
  children: ReactNode;
}

export const MovieProvider: React.FC<MovieProviderProps> = ({ children }) => {
  const [movies, setMovies] = useState<Movie[]>([]);

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        const response = await fetch('http://localhost:8080/v1/movies');
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data = await response.json();
        setMovies(data);
      } catch (error) {
        console.error("Failed to fetch movies:", error);
        // Optionally, set some error state here
      }
    };

    fetchMovies();
  }, []);

  const addMovie = (movieInput: CreateMovieInput) => {
    // This will be replaced with an API call
    const newMovie: Movie = {
      ...movieInput,
      _id: Date.now().toString(),
    };
    setMovies([...movies, newMovie]);
  };

  const updateMovie = (id: string, updates: Partial<Movie>) => {
    const newMovies = movies.map(movie =>
      movie._id === id ? { ...movie, ...updates } : movie
    );
    setMovies(newMovies);
  };

  const deleteMovie = (id: string) => {
    const newMovies = movies.filter(movie => movie._id !== id);
    setMovies(newMovies);
  };

  const getMovie = (id: string) => {
    return movies.find(movie => movie.imdb_id === id);
  };

  return (
    <MovieContext.Provider value={{ movies, addMovie, updateMovie, deleteMovie, getMovie }}>
      {children}
    </MovieContext.Provider>
  );
};