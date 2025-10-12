# MovieApp Frontend

A modern movie management application built with Vite, React, TypeScript, and Tailwind CSS.

## Features

- **Movie Management**: Add, view, edit, and delete movies
- **User Authentication**: Login and register functionality
- **Responsive Design**: Works on desktop and mobile devices
- **Dark Mode**: Toggle between light and dark themes
- **Search & Pagination**: Find movies easily with search and pagination

## Tech Stack

- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui with Radix UI
- **Routing**: React Router
- **State Management**: React Context
- **Testing**: Vitest + React Testing Library

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

4. Open [http://localhost:5173](http://localhost:5173) in your browser.

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint
- `npm run test` - Run tests with Vitest

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── ui/             # shadcn/ui components
│   └── Navbar.tsx      # Main navigation
├── contexts/           # React contexts for state management
│   ├── AuthContext.tsx # Authentication state
│   └── MovieContext.tsx # Movie data state
├── pages/              # Page components
│   ├── Homepage.tsx
│   ├── Login.tsx
│   ├── Register.tsx
│   ├── Movies.tsx
│   ├── AddMovie.tsx
│   └── MovieDetails.tsx
├── types/              # TypeScript type definitions
├── utils/              # Utility functions
├── test/               # Test setup
└── App.tsx             # Main app component
```

## Authentication

The app uses client-side authentication with localStorage persistence. Users can:

- Register with name, email, and password
- Login with email and password
- Logout to clear session

Protected routes (Add Movie, Edit/Delete Movie Details) require authentication.

## Movie Management

Movies are stored in localStorage with the following fields:

- Title
- Description
- Release Date
- Poster URL
- Genres (comma-separated)
- Rating (0-10)

## Extending to Real Backend

To connect to a real backend API:

1. Replace mock functions in `AuthContext.tsx` and `MovieContext.tsx` with actual API calls
2. Update API endpoints to match your backend routes
3. Add error handling for network requests
4. Implement proper authentication tokens

Example API integration:

```typescript
// In AuthContext.tsx
const login = async (input: LoginInput) => {
  const response = await fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  });
  const data = await response.json();
  // Handle response...
};
```

## Testing

Run tests with:
```bash
npm run test
```

Tests include:
- Utility function tests
- Component validation tests
- Using Vitest and React Testing Library

## Deployment

Build for production:
```bash
npm run build
```

The build artifacts will be stored in the `dist/` directory.

## Contributing

1. Follow the existing code style
2. Add tests for new features
3. Update documentation as needed
4. Ensure all tests pass before submitting

## License

This project is part of the MovieApp application.

You can also install [eslint-plugin-react-x](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```
