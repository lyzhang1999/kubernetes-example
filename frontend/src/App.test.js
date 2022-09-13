import { render, screen } from '@testing-library/react';
import App from './App';

test('renders full stack', () => {
  render(<App />);
  const textElement = screen.getByText(/full-stack/i);
  expect(textElement).toBeInTheDocument();
});
