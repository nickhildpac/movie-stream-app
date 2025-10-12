import { Link } from 'react-router-dom';
import { useMovies } from '../contexts/MovieContext';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';
import { Button } from '../components/ui/button';

const Homepage = () => {
  const { movies } = useMovies();

  const featuredMovies = movies.slice(0, 3); // First 3 movies as featured

  const totalMovies = movies.length;
  const averageRating = movies.length > 0
    ? (movies.reduce((sum, movie) => sum + movie.rating, 0) / movies.length).toFixed(1)
    : '0';
  const uniqueGenres = new Set(movies.flatMap(movie => movie.genres)).size;

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Hero Section */}
      <section className="text-center py-16">
        <h1 className="text-4xl font-bold mb-4">Welcome to MovieApp</h1>
        <p className="text-xl text-muted-foreground mb-8">
          Discover, manage, and share your favorite movies
        </p>
        <Button asChild size="lg">
          <Link to="/movies">Browse Movies</Link>
        </Button>
      </section>

      {/* Stats Section */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <Card>
          <CardHeader>
            <CardTitle>Total Movies</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{totalMovies}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Average Rating</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{averageRating}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Genres</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{uniqueGenres}</p>
          </CardContent>
        </Card>
      </section>

      {/* Featured Movies */}
      <section>
        <h2 className="text-2xl font-bold mb-6">Featured Movies</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {featuredMovies.map((movie) => (
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
                <p className="text-sm text-muted-foreground">Rating: {movie.rating}/10</p>
                <Button asChild className="mt-2">
                  <Link to={`/movies/${movie.id}`}>View Details</Link>
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>
    </div>
  );
};

export default Homepage;