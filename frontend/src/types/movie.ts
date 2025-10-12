export interface Movie {
  id: string;
  title: string;
  description: string;
  releaseDate: string;
  posterUrl: string;
  genres: string[];
  rating: number;
}

export interface CreateMovieInput {
  title: string;
  description: string;
  releaseDate: string;
  posterUrl: string;
  genres: string[];
  rating: number;
}