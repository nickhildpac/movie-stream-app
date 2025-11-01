import { useState, useEffect } from "react";
import { Label } from "../components/ui/label";
import type { CSSObjectWithLabel } from "react-select";
import Select from "react-select";
import type { Genre } from "../types";

interface GenreSelectProps {
  value: Genre[];
  onChange: (genres: Genre[]) => void;
  className?: string;
}

const GenreSelect = ({ value, onChange, className }: GenreSelectProps) => {
  const [genres, setGenres] = useState<{ value: Genre; label: string }[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchGenres = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/genres`);
        const data = await response.json();
        setGenres(
          data.map((genre: Genre) => ({
            value: genre,
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
    <div className={className}>
      <Label htmlFor="genres">Favourite Genres</Label>
      <Select
        id="genres"
        isMulti
        options={genres}
        styles={customStyles}
        value={genres.filter(option => value.some(favGenre => favGenre.genre_id === option.value.genre_id))}
        onChange={(selectedOptions) =>
          onChange(
            selectedOptions.map((option) => option.value),
          )
        }
      />
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
};

export default GenreSelect;