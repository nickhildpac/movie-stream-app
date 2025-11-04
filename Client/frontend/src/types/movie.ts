export interface Genre {
  genre_id: number;
  genre_name: string;
}

export interface Ranking {
  ranking_value: number;
  ranking_name: string;
}

export interface Movie {
  _id: string;
  imdb_id: string;
  title: string;
  poster_path: string;
  youtube_id: string;
  genre: Genre[];
  admin_review: string;
  ranking: Ranking;
}

export interface CreateMovieInput {
  imdb_id: string;
  title: string;
  poster_path: string;
  youtube_id: string;
  genre: Genre[];
  admin_review: string;
  ranking: Ranking;
}



export interface TMDBMovie {
  id: number;
  title: string;
  poster_path: string;
  overview: string;
  release_date: string;
  genres: { id: number; name: string }[];
}







export interface TMDBMovieSearchResult {
  id: number;
  title: string;
  poster_path: string;
  release_date: string;
  genre_ids: number[];
  overview: string;
  backdrop_path: string;
}




