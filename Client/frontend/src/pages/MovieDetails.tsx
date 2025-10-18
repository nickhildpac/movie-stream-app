import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { useMovies } from '../contexts/MovieContext';
import type { Movie } from '../types';
import { Button } from '../components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '../components/ui/dialog';
import { Input } from '../components/ui/input';
import { Label } from '../components/ui/label';
import { Textarea } from '../components/ui/textarea';

const MovieDetails = () => {
  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();
  const { getMovie, updateMovie, deleteMovie } = useMovies();
  const navigate = useNavigate();
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isDeleteConfirmOpen, setIsDeleteConfirmOpen] = useState(false);

  const movie = id ? getMovie(id) : undefined;

  if (!movie) {
    return <div className="container mx-auto px-4 py-8">Movie not found</div>;
  }

  const handleDelete = () => {
    if (movie) {
      deleteMovie(movie._id);
      setIsDeleteConfirmOpen(false);
      navigate('/movies');
    }
  };

  const handleEdit = (updatedMovie: Partial<Movie>) => {
    if (movie) {
      updateMovie(movie._id, updatedMovie);
      setIsEditOpen(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div>
          <img
            src={movie.poster_path}
            alt={movie.title}
            className="w-full rounded-lg shadow-lg"
          />
        </div>
        <div>
          <Card>
            <CardHeader>
              <CardTitle className="text-3xl">{movie.title}</CardTitle>
              <CardDescription className="text-lg">
                Ranking: {movie.ranking.ranking_value} | {movie.ranking.ranking_name} | Genres: {movie.genre.map(g => g.genre_name).join(', ')}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <p><strong>Admin Review:</strong> {movie.admin_review}</p>
              {user && (
                <div className="flex space-x-2 pt-4">
                  <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
                    <DialogTrigger asChild>
                      <Button variant="outline">Edit</Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Edit Movie</DialogTitle>
                        <DialogDescription>Update movie details</DialogDescription>
                      </DialogHeader>
                      <EditMovieForm movie={movie} onSave={handleEdit} />
                    </DialogContent>
                  </Dialog>
                  <Dialog open={isDeleteConfirmOpen} onOpenChange={setIsDeleteConfirmOpen}>
                    <DialogTrigger asChild>
                      <Button variant="destructive">Delete</Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Confirm Delete</DialogTitle>
                        <DialogDescription>
                          Are you sure you want to delete "{movie.title}"? This action cannot be undone.
                        </DialogDescription>
                      </DialogHeader>
                      <div className="flex justify-end space-x-2">
                        <Button variant="outline" onClick={() => setIsDeleteConfirmOpen(false)}>
                          Cancel
                        </Button>
                        <Button variant="destructive" onClick={handleDelete}>
                          Delete
                        </Button>
                      </div>
                    </DialogContent>
                  </Dialog>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

const EditMovieForm = ({ movie, onSave }: { movie: Movie; onSave: (movie: Partial<Movie>) => void }) => {
  const [title, setTitle] = useState(movie.title);
  const [description, setDescription] = useState(movie.admin_review);
  const [posterUrl, setPosterUrl] = useState(movie.poster_path);
  const [genres, setGenres] = useState(movie.genre.map((genre: { genre_name: string,genre_id: number}) => genre.genre_name).join(', '));
  const [rating, setRating] = useState(movie.ranking.ranking_value.toString());

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave({
      ...movie,
      title,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <Label htmlFor="edit-title">Title</Label>
        <Input id="edit-title" value={title} onChange={(e) => setTitle(e.target.value)} />
      </div>
      <div>
        <Label htmlFor="edit-description">Review</Label>
        <Textarea id="edit-description" value={description} onChange={(e) => setDescription(e.target.value)} />
      </div>
      <div>
        <Label htmlFor="edit-posterUrl">Poster URL</Label>
        <Input id="edit-posterUrl" type="url" value={posterUrl} onChange={(e) => setPosterUrl(e.target.value)} />
      </div>
      <div>
        <Label htmlFor="edit-genres">Genres</Label>
        <Input id="edit-genres" value={genres} onChange={(e) => setGenres(e.target.value)} />
      </div>
      <div>
        <Label htmlFor="edit-rating">Ranking</Label>
        <Input id="edit-rating" type="number" min="0" max="10" step="0.1" value={rating} onChange={(e) => setRating(e.target.value)} />
      </div>
      <Button type="submit">Save Changes</Button>
    </form>
  );
};

export default MovieDetails;