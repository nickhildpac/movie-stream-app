import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useMovies } from '../contexts/MovieContext';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';

const Movies = () => {
  const { movies } = useMovies();
  const [search, setSearch] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const moviesPerPage = 6;

  const filteredMovies = movies.filter(movie =>
    movie.title.toLowerCase().includes(search.toLowerCase()) ||
    movie.description.toLowerCase().includes(search.toLowerCase())
  );

  const totalPages = Math.ceil(filteredMovies.length / moviesPerPage);
  const startIndex = (currentPage - 1) * moviesPerPage;
  const paginatedMovies = filteredMovies.slice(startIndex, startIndex + moviesPerPage);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Movies</h1>
        <Button asChild>
          <Link to="/movies/add">Add Movie</Link>
        </Button>
      </div>

      <div className="mb-6">
        <Input
          placeholder="Search movies..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-md"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-6">
        {paginatedMovies.map((movie) => (
          <Card key={movie.id} className="overflow-hidden">
            <img
              src={movie.posterUrl}
              alt={movie.title}
              className="w-full h-48 object-cover"
            />
            <CardHeader>
              <CardTitle>{movie.title}</CardTitle>
              <CardDescription>{movie.description}</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground mb-2">
                Rating: {movie.rating}/10 | Genres: {movie.genres.join(', ')}
              </p>
              <p className="text-sm text-muted-foreground mb-4">
                Release: {new Date(movie.releaseDate).toLocaleDateString()}
              </p>
              <Button asChild>
                <Link to={`/movies/${movie.id}`}>View Details</Link>
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>

      {totalPages > 1 && (
        <div className="flex justify-center space-x-2">
          <Button
            variant="outline"
            onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
            disabled={currentPage === 1}
          >
            Previous
          </Button>
          <span className="self-center">
            Page {currentPage} of {totalPages}
          </span>
          <Button
            variant="outline"
            onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
            disabled={currentPage === totalPages}
          >
            Next
          </Button>
        </div>
      )}
    </div>
  );
};

export default Movies;