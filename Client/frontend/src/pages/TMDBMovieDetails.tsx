import { useEffect, useState } from "react";
import { useParams } from "react-router-dom"; 
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import type { Genre, TMDBMovie } from "../types/movie";
import { Button } from "@/components/ui/button";
import { toast } from "@/hooks/use-toast";

const TMDBMovieDetails = () => {
  const { id } = useParams<{ id: string }>();
  const [movie, setMovie] = useState<TMDBMovie | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMovieDetails = async () => {
      try {
        const response = await fetch(
          `https://api.themoviedb.org/3/movie/${id}`,
          {
            method: "GET",
            headers: {
              Authorization: `Bearer ${import.meta.env.VITE_TMDB_BEARER_TOKEN}`,
              accept: "application/json",
            },
          }
        );
        if (!response.ok) {
          throw new Error("Failed to fetch movie details");
        }
        const data = await response.json();
        setMovie(data);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchMovieDetails();
    }
  }, [id]);

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Loading...</div>;
  }

  if (!movie) {
    return <div className="container mx-auto px-4 py-8">Movie not found</div>;
  }
  async function addToMovie() {
    const genres: Genre[] = []
    movie?.genres.forEach((genre) => {
      genres.push({
        genre_name: genre.name,
        genre_id: genre.id
      })
      fetch(`${import.meta.env.VITE_API_BASE_URL}/genre`, {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          genre_name: genre.name,
          genre_id: genre.id
        })
      });
    })
    const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/addmovie`, {
      credentials: "include",
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        imdb_id: movie?.id.toString(),
        title: movie?.title,
        release_date: movie?.release_date,
        poster_path:`https://image.tmdb.org/t/p/w500${movie?.poster_path}`,
        overview: movie?.overview,
        genre: genres
      })
    });
    if (response.ok) {
        const data = await response.json();
        console.log(data);
        toast({
          title: "Movie added!",
          description: movie?.title +" added to list",
          variant: "success",
        });
      }
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div>
          <img
            src={`https://image.tmdb.org/t/p/w500${movie.poster_path}`}
            alt={movie.title}
            className="w-full rounded-lg shadow-lg"
          />
        </div>
        <div>
          <Card>
            <CardHeader>
              <CardTitle className="text-3xl">{movie.title}</CardTitle>
              <CardDescription className="text-lg">
                {movie.release_date} | {movie.genres.map((g) => g.name).join(", ")}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <p>
                <strong>Overview:</strong> {movie.overview}
              </p>
            </CardContent>
                  <CardFooter className="p-4 bg-secondary/50 flex justify-between items-center">
                    <Button className="w-1/2" onClick={addToMovie}>
                        Add to Movies
                    </Button>
                  </CardFooter>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default TMDBMovieDetails;
