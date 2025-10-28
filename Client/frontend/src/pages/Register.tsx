import { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Label } from "../components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import type { CSSObjectWithLabel } from "react-select";
import Select from "react-select";

const Register = () => {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [favouriteGenres, setFavouriteGenres] = useState<string[]>([]);
  const [genres, setGenres] = useState<{ value: string; label: string }[]>([]);
  const [error, setError] = useState("");
  const { register } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchGenres = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/genres`);
        const data = await response.json();
        setGenres(
          data.map((genre: { genre_name: number; genre_id: string }) => ({
            value: { genre_id: genre.genre_id, genre_name: genre.genre_name },
            label: genre.genre_name,
          })),
        );
      } catch (error) {
        console.error("Failed to fetch genres:", error);
        setError("Failed to load genres. Please try again later.");
      }
    };
    fetchGenres();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }
    try {
      await register({
        first_name: firstName,
        last_name: lastName,
        email,
        password,
        confirmPassword,
        favourite_genres: favouriteGenres,
      });
      navigate("/");
    } catch {
      setError("Registration failed");
    }
  };

  const customStyles = {
    control: (provided: CSSObjectWithLabel) => ({
      ...provided,
      backgroundColor: "#2d3748",
      borderColor: "#4a5568",
      color: "white",
    }),
    option: (provided: CSSObjectWithLabel, state: { isFocused: boolean }) => ({
      ...provided,
      backgroundColor: state.isFocused ? "#4a5568" : "#2d3748",
      color: "white",
    }),
    multiValue: (provided: CSSObjectWithLabel) => ({
      ...provided,
      backgroundColor: "#4a5568",
    }),
    multiValueLabel: (provided: CSSObjectWithLabel) => ({
      ...provided,
      color: "white",
    }),
    multiValueRemove: (provided: CSSObjectWithLabel) => ({
      ...provided,
      color: "white",
      ":hover": {
        backgroundColor: "#718096",
        color: "white",
      },
    }),
  };

  return (
    <div className="container mx-auto px-4 py-8 flex justify-center">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Register</CardTitle>
          <CardDescription>Create a new account</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="flex space-x-4">
              <div className="w-1-2">
                <Label htmlFor="firstName">First Name</Label>
                <Input
                  id="firstName"
                  type="text"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  required
                />
              </div>
              <div className="w-1/2">
                <Label htmlFor="lastName">Last Name</Label>
                <Input
                  id="lastName"
                  type="text"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  required
                />
              </div>
            </div>
            <div>
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div>
              <Label htmlFor="genres">Favourite Genres</Label>
              <Select
                id="genres"
                isMulti
                options={genres}
                styles={customStyles}
                onChange={(selectedOptions) =>
                  setFavouriteGenres(
                    selectedOptions.map((option) => option.value),
                  )
                }
              />
            </div>
            <div>
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <div>
              <Label htmlFor="confirmPassword">Confirm Password</Label>
              <Input
                id="confirmPassword"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
            </div>
            {error && <p className="text-red-500 text-sm">{error}</p>}
            <Button type="submit" className="w-full">
              Register
            </Button>
          </form>
          <p className="text-center mt-4">
            Already have an account?{" "}
            <Link to="/login" className="text-primary hover:underline">
              Login
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
};

export default Register;
