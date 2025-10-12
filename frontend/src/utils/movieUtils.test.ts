import { describe, it, expect } from 'vitest';
import { formatRating, getGenresString, isValidRating } from './movieUtils';

describe('movieUtils', () => {
  describe('formatRating', () => {
    it('should format rating correctly', () => {
      expect(formatRating(8.5)).toBe('8.5/10');
      expect(formatRating(9)).toBe('9.0/10');
    });
  });

  describe('getGenresString', () => {
    it('should join genres with comma', () => {
      expect(getGenresString(['Action', 'Drama'])).toBe('Action, Drama');
      expect(getGenresString(['Sci-Fi'])).toBe('Sci-Fi');
    });
  });

  describe('isValidRating', () => {
    it('should validate rating range', () => {
      expect(isValidRating(5)).toBe(true);
      expect(isValidRating(0)).toBe(true);
      expect(isValidRating(10)).toBe(true);
      expect(isValidRating(-1)).toBe(false);
      expect(isValidRating(11)).toBe(false);
    });
  });
});