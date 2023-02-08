import React from 'react';
import { render, screen } from '@testing-library/react';
import App from '.';

test('renders text "Ansible Automation Platform Deployment Engine"', () => {
  render(<App />);
  const linkElement = screen.getByText("Documentation");
  expect(linkElement).toBeInTheDocument();
});
