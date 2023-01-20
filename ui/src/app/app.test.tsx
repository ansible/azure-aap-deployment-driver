import React from 'react';
import { render, screen } from '@testing-library/react';
import App from '.';

test('renders text "Ansible Automation Platform Installer"', () => {
  render(<App />);
  const linkElement = screen.getByText("Ansible Automation Platform Deployment driver web UI not available yet.");
  expect(linkElement).toBeInTheDocument();
});
