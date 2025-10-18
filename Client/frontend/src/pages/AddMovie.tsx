import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMovies } from '../contexts/MovieContext';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { Label } from '../components/ui/label';
import { Textarea } from '../components/ui/textarea';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';

const AddMovie = () => {
  const { addMovie } = useMovies();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [releaseDate, setReleaseDate] = useState('');
  const [posterUrl, setPosterUrl] = useState('');
  const [genres, setGenres] = useState('');
  const [rating, setRating] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const navigate = useNavigate();

  const validateForm = () => {
    const newErrors: Record<string, string> = {};
    if (!title.trim()) newErrors.title = 'Title is required';
    if (!description.trim()) newErrors.description = 'Description is required';
    if (!releaseDate) newErrors.releaseDate = 'Release date is required';
    if (!posterUrl.trim()) newErrors.posterUrl = 'Poster URL is required';
    if (!genres.trim()) newErrors.genres = 'Genres are required';
    const ratingNum = parseFloat(rating);
    if (isNaN(ratingNum) || ratingNum < 0 || ratingNum > 10) {
      newErrors.rating = 'Rating must be a number between 0 and 10';
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;

    addMovie({
      title,
      description,
      releaseDate,
      posterUrl,
      genres: genres.split(',').map(g => g.trim()),
      rating: parseFloat(rating),
    });

    // Reset form
    setTitle('');
    setDescription('');
    setReleaseDate('');
    setPosterUrl('');
    setGenres('');
    setRating('');
    setErrors({});

    navigate('/movies');
  };

  return (
    <div className="container mx-auto px-4 py-8 flex justify-center">
      <Card className="w-full max-w-2xl">
        <CardHeader>
          <CardTitle>Add New Movie</CardTitle>
          <CardDescription>Fill in the details to add a movie to the collection</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <Label htmlFor="title">Title</Label>
              <Input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
              />
              {errors.title && <p className="text-red-500 text-sm">{errors.title}</p>}
            </div>
            <div>
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                required
              />
              {errors.description && <p className="text-red-500 text-sm">{errors.description}</p>}
            </div>
            <div>
              <Label htmlFor="releaseDate">Release Date</Label>
              <Input
                id="releaseDate"
                type="date"
                value={releaseDate}
                onChange={(e) => setReleaseDate(e.target.value)}
                required
              />
              {errors.releaseDate && <p className="text-red-500 text-sm">{errors.releaseDate}</p>}
            </div>
            <div>
              <Label htmlFor="posterUrl">Poster URL</Label>
              <Input
                id="posterUrl"
                type="url"
                value={posterUrl}
                onChange={(e) => setPosterUrl(e.target.value)}
                required
              />
              {errors.posterUrl && <p className="text-red-500 text-sm">{errors.posterUrl}</p>}
            </div>
            <div>
              <Label htmlFor="genres">Genres (comma separated)</Label>
              <Input
                id="genres"
                value={genres}
                onChange={(e) => setGenres(e.target.value)}
                placeholder="Action, Drama, Sci-Fi"
                required
              />
              {errors.genres && <p className="text-red-500 text-sm">{errors.genres}</p>}
            </div>
            <div>
              <Label htmlFor="rating">Rating (0-10)</Label>
              <Input
                id="rating"
                type="number"
                min="0"
                max="10"
                step="0.1"
                value={rating}
                onChange={(e) => setRating(e.target.value)}
                required
              />
              {errors.rating && <p className="text-red-500 text-sm">{errors.rating}</p>}
            </div>
            <Button type="submit" className="w-full">
              Add Movie
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default AddMovie;