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

const initialMovies: Movie[] = [
  {
    id: '1',
    title: 'Inception',
    description: 'A mind-bending thriller about dream invasion',
    releaseDate: '2010-07-16',
    posterUrl: 'https://via.placeholder.com/300x400?text=Inception',
    genres: ['Action', 'Sci-Fi'],
    rating: 8.8,
  },
  {
    id: '2',
    title: 'The Dark Knight',
    description: 'Batman faces the Joker in Gotham',
    releaseDate: '2008-07-18',
    posterUrl: 'https://via.placeholder.com/300x400?text=Dark+Knight',
    genres: ['Action', 'Crime'],
    rating: 9.0,
  },
  {
    id: '3',
    title: 'Interstellar',
    description: 'A journey through space and time',
    releaseDate: '2014-11-07',
    posterUrl: 'https://via.placeholder.com/300x400?text=Interstellar',
    genres: ['Adventure', 'Drama'],
    rating: 8.6,
  },
  {
    id: '4',
    title: 'The Matrix',
    description: 'A hacker discovers the truth about reality',
    releaseDate: '1999-03-31',
    posterUrl: 'https://via.placeholder.com/300x400?text=Matrix',
    genres: ['Action', 'Sci-Fi'],
    rating: 8.7,
  },
  {
    id: '5',
    title: 'Pulp Fiction',
    description: 'Non-linear interconnected tales of crime',
    releaseDate: '1994-10-14',
    posterUrl: 'https://via.placeholder.com/300x400?text=Pulp+Fiction',
    genres: ['Crime', 'Drama'],
    rating: 8.9,
  },
];

export const MovieProvider: React.FC<MovieProviderProps> = ({ children }) => {
  const [movies, setMovies] = useState<Movie[]>([]);

  useEffect(() => {
    // Load from localStorage or use initial
    const stored = localStorage.getItem('movies');
    if (stored) {
      setMovies(JSON.parse(stored));
    } else {
      setMovies(initialMovies);
      localStorage.setItem('movies', JSON.stringify(initialMovies));
    }
  }, []);

  const saveMovies = (newMovies: Movie[]) => {
    setMovies(newMovies);
    localStorage.setItem('movies', JSON.stringify(newMovies));
  };

  const addMovie = (movieInput: CreateMovieInput) => {
    const newMovie: Movie = {
      ...movieInput,
      id: Date.now().toString(),
    };
    saveMovies([...movies, newMovie]);
  };

  const updateMovie = (id: string, updates: Partial<Movie>) => {
    const newMovies = movies.map(movie =>
      movie.id === id ? { ...movie, ...updates } : movie
    );
    saveMovies(newMovies);
  };

  const deleteMovie = (id: string) => {
    const newMovies = movies.filter(movie => movie.id !== id);
    saveMovies(newMovies);
  };

  const getMovie = (id: string) => {
    return movies.find(movie => movie.id === id);
  };

  return (
    <MovieContext.Provider value={{ movies, addMovie, updateMovie, deleteMovie, getMovie }}>
      {children}
    </MovieContext.Provider>
  );
};