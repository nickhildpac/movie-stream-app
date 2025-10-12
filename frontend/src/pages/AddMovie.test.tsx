import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MovieProvider } from '../contexts/MovieContext';
import AddMovie from './AddMovie';

// Mock useNavigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

const renderWithProviders = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      <MovieProvider>
        {component}
      </MovieProvider>
    </BrowserRouter>
  );
};

describe('AddMovie', () => {
  it('shows validation errors for empty form', async () => {
    renderWithProviders(<AddMovie />);

    const submitButton = screen.getByRole('button', { name: /add movie/i });
    fireEvent.click(submitButton);

    expect(screen.getByText('Title is required')).toBeInTheDocument();
    expect(screen.getByText('Description is required')).toBeInTheDocument();
    expect(screen.getByText('Release date is required')).toBeInTheDocument();
    expect(screen.getByText('Poster URL is required')).toBeInTheDocument();
    expect(screen.getByText('Genres are required')).toBeInTheDocument();
  });

  it('shows validation error for invalid rating', async () => {
    renderWithProviders(<AddMovie />);

    const ratingInput = screen.getByLabelText(/rating/i);
    fireEvent.change(ratingInput, { target: { value: '15' } });

    const submitButton = screen.getByRole('button', { name: /add movie/i });
    fireEvent.click(submitButton);

    expect(screen.getByText('Rating must be a number between 0 and 10')).toBeInTheDocument();
  });
});