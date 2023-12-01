import React from 'react';
import { render, screen } from '@testing-library/react';
import App from '.';

test('The path / renders login form', () => {
  render(<App />);
  let linkElement;
  linkElement = screen.getByText("Deployment Engine");
  expect(linkElement).toBeInTheDocument();

  // the description is also used in another hidden element so we need to get them all but only verify
  linkElement = screen.getAllByText("Please use the administrative credentials for Red Hat Ansible Automation Platform on Microsoft Azure.");
  expect(linkElement[0]).toBeInTheDocument();

  // user name should be hard-coded to admin
  linkElement = screen.getByText("Username");
  expect(linkElement).toBeInTheDocument();
  linkElement = screen.getByDisplayValue("admin");
  expect(linkElement).toBeInTheDocument();
  linkElement = screen.getByText("Password");
  expect(linkElement).toBeInTheDocument();
  // there should be a login button
  linkElement = screen.getByText("Log in");
  expect(linkElement).toBeInTheDocument();
});
