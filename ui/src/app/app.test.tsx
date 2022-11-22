import React from 'react';
import { render, screen } from '@testing-library/react';
import App from '.';
import { Provider } from 'react-redux';
import store from './redux/store';

test('renders text "Ansible Automation Platform Installer"', () => {
  render(<Provider store={store}><App /></Provider>);
  const linkElement = screen.getByText("Ansible Automation Platform Installer");
  expect(linkElement).toBeInTheDocument();
});
