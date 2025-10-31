import { useState } from 'react';
import { useMovies } from '../contexts/MovieContext';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../components/ui/card';
import { Input } from '../components/ui/input';
import { Link } from 'react-router-dom';
import { Button } from '../components/ui/button';

const Movies = () => {
  const { movies } = useMovies();
  const [search, setSearch] = useState('');

  const filteredMovies = movies.filter(movie =>
    movie.title.toLowerCase().includes(search.toLowerCase()) ||
    movie.admin_review.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Movies</h1>
        <div className="w-1/3">
          <Input
            type="text"
            placeholder="Search movies..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {filteredMovies.map(movie => (
          <Card key={movie._id} className="overflow-hidden">
            <CardHeader className="p-0">
              <img src={movie.poster_path} alt={movie.title} className="w-full h-48 object-cover" />
            </CardHeader>
            <CardContent className="p-4 text-center">
              <CardTitle>{movie.title}</CardTitle>
              <CardDescription>{movie.genre.map(g => g.genre_name).join(', ')}</CardDescription>
            </CardContent>
            <CardFooter className="p-4 bg-secondary/50 flex justify-between items-center">
              <Button className='w-full'>
                <Link to={`/movies/${movie.imdb_id}`}>{movie.ranking.ranking_name}</Link>
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default Movies;