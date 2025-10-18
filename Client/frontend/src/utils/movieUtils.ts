export const formatRating = (rating: number): string => {
  return `${rating.toFixed(1)}/10`;
};

export const getGenresString = (genres: string[]): string => {
  return genres.join(', ');
};

export const isValidRating = (rating: number): boolean => {
  return rating >= 0 && rating <= 10;
};