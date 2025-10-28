import { type Genre } from "./movie";

export interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  first_name?: string;
  last_name?: string;
  favourite_genres?: Genre[];
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
  confirmPassword: string;
  favourite_genres: string[];
}

export interface UpdateUserInput {
  first_name?: string;
  last_name?: string;
  email?: string;
  password?: string;
  favourite_genres?: Genre[];
}
