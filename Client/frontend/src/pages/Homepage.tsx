import { useEffect } from 'react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../components/ui/card';
import { Link } from 'react-router-dom';
import { Button } from '../components/ui/button';
import { useMovies } from '../contexts/MovieContext';

const Homepage = () => {
  const { recommendedMovies, fetchRecommendedMovies } = useMovies();

  useEffect(() => {
    fetchRecommendedMovies();
  }, []);

  return (
    <div className="container mx-auto px-4 py-8">
      <header className="text-center mb-8">
        <h1 className="text-4xl font-bold">Welcome to MovieApp</h1>
        <p className="text-lg text-muted-foreground">Your ultimate destination for movie streaming and reviews.</p>
      </header>

      <section className="mb-8">
        <h2 className="text-2xl font-semibold mb-4">At a Glance</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-center">
          <Card>
            <CardHeader>
              <CardTitle>Total Movies</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{recommendedMovies.length}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Average Rating</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{5}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Genres Available</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{10}</p>
            </CardContent>
          </Card>
        </div>
      </section>
      <section>
        <h2 className="text-2xl font-semibold mb-4">Recommended Movies</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {recommendedMovies.map(movie => (
            <Card key={movie._id} className="overflow-hidden">
              <CardHeader className="p-0">
                <img src={movie.poster_path} alt={movie.title} className="w-full h-48 object-cover" />
              </CardHeader>
              <CardContent className="p-4">
                <CardTitle>{movie.title}</CardTitle>
                <CardDescription>{movie.admin_review}</CardDescription>
              </CardContent>
              <CardFooter className="p-4 bg-secondary/50 flex justify-between items-center">
                <p className="text-sm text-muted-foreground">Rating: {movie.ranking.ranking_value}/10</p>
                <Button asChild size="sm">
                  <Link to={`/movies/${movie.imdb_id}`}>View Details</Link>
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      </section>
    </div>
  );
};

export default Homepage;
