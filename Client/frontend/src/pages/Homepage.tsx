import { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { Link } from "react-router-dom";
import { Button } from "../components/ui/button";
import { useMovies } from "../contexts/MovieContext";
import { Input } from "@/components/ui/input";

const Homepage = () => {
  const { recommendedMovies, fetchRecommendedMovies } = useMovies();
  const [search, setSearch] = useState("");
  const [searchResults, setSearchResults] = useState<any[]>([]);

  useEffect(() => {
    fetchRecommendedMovies();
  }, [fetchRecommendedMovies]);
  const TMDB_API_URL = "https://api.themoviedb.org/3/search/movie";

  function getBearerToken(): string | undefined {
    return import.meta.env.VITE_TMDB_BEARER_TOKEN;
  }

  const handleSearch = async () => {
    if (!search.trim()) return;

    try {
      const url = new URL(TMDB_API_URL);
      url.searchParams.set("query", search);
      url.searchParams.set("include_adult", "false");
      url.searchParams.set("language", "en-US");
      url.searchParams.set("page", "1");

      const response = await fetch(url.toString(), {
        method: "GET",
        headers: {
          Authorization: `Bearer ${getBearerToken()}`,
          accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch search results");
      }

      const data = await response.json();
      setSearchResults(data.results || []);
    } catch (error) {
      console.error("Error searching movies:", error);
      setSearchResults([]);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <header className="text-center mb-8">
        <h1 className="text-4xl font-bold">Welcome to MovieApp</h1>
        <p className="text-lg text-muted-foreground">
          Your ultimate destination for movie streaming and reviews.
        </p>
      </header>

      <section className="mb-8">
        <div className="w-full flex gap-2">
          <Input
            type="text"
            placeholder="Search movies..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSearch()}
          />
          <Button onClick={handleSearch}>Search Movie</Button>
        </div>
        {searchResults.length > 0 && (
          <div className="mt-4">
            <h3 className="text-xl font-semibold mb-4">Search Results</h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
              {searchResults.map((movie) => (
                <Card key={movie.id} className="overflow-hidden">
                  <CardHeader className="p-0">
                    <img
                      src={
                        movie.poster_path
                          ? `https://image.tmdb.org/t/p/w500${movie.poster_path}`
                          : "/placeholder-movie.jpg"
                      }
                      alt={movie.title}
                      className="w-full h-48 object-cover"
                    />
                  </CardHeader>
                  <CardContent className="p-4">
                    <CardTitle className="text-center">{movie.title}</CardTitle>
                    <CardDescription className="text-center">
                      {movie.release_date}
                    </CardDescription>
                  </CardContent>
                  <CardFooter className="p-4 bg-secondary/50 flex justify-between items-center">
                    <Button className="w-full" asChild>
                      <Link to={`/tmdb-movies/${movie.id}`}>
                        View Details
                      </Link>
                    </Button>
                  </CardFooter>
                </Card>
              ))}
            </div>
          </div>
        )}
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
          {recommendedMovies.map((movie) => (
            <Card key={movie.imdb_id} className="overflow-hidden">
              <CardHeader className="p-0">
                <img
                  src={movie.poster_path}
                  alt={movie.title}
                  className="w-full h-48 object-cover"
                />
              </CardHeader>
              <CardContent className="p-4">
                <CardTitle className="text-center">{movie.title}</CardTitle>
                <CardDescription className="text-center">
                  {movie.genre.map((g) => g.genre_name).join(", ")}
                </CardDescription>
              </CardContent>
              <CardFooter className="p-4 bg-secondary/50 flex justify-between items-center">
                <Button className="w-full">
                  <Link to={`/movies/${movie.imdb_id}`}>
                    {movie.ranking.ranking_name}
                  </Link>
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
