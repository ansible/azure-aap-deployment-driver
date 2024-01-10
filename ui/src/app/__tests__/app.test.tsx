import React from 'react';
import { RouterProvider, createMemoryRouter } from "react-router-dom";
import { render, screen } from '@testing-library/react';
import App from '..';
import appRoutes from '../app.routes';

describe('Main app', ()=>{

  it('route /login renders login form', ()=>{
    navigateTo('/login')
    verifyInitialLoginForm()
  })

  it('route /rhlogin renders deployment component with SSO login dialog over it', ()=>{
    navigateTo('/rhlogin')
    const dialog = screen.getByRole('dialog')
    expect(dialog).toBeVisible()
    const login = screen.getByRole('link', {name:'Log in with Red Hat account'})
    expect(login).toBeVisible()
  })

  it('route / renders deployment component without the SSO login dialog', ()=>{
    //const user = userEvent.setup()
    navigateTo('/')
    const dialog = screen.queryByRole('dialog')
    expect(dialog).toBeNull()
  })
})

// helper verification functions

const navigateTo = (path:string) => {
  const router = createMemoryRouter(appRoutes, {
    initialEntries: [path],
  });
  // what's passed to the render() is the same what's done index.tsx file
  render(<RouterProvider router={router} />);
}

// helper function because login form needs to be verified twice in this test
const verifyInitialLoginForm = () => {
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
}
