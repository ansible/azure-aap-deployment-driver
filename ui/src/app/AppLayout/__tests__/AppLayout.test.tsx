import React from 'react';
import { RouterProvider, createMemoryRouter } from 'react-router-dom';
import { render, screen } from '@testing-library/react';
import { AppLayout } from '../AppLayout';

describe('AppLayout', ()=>{
	it('renders Page with proper AAP logo', ()=> {
		const routes = [
			{
				element: <AppLayout navigation={[]} />,
				children: [
					{
						path: "/",
						element: <></>
					}
				]
			}
		]

		const router = createMemoryRouter(routes, {
			initialEntries: ["/"],
		});
		// what's passed to the render() is the same what's done index.tsx file
		render(<RouterProvider router={router} />);

		const logo = screen.getByAltText("Red Hat Ansible Automation Platform Logo")
		expect(logo).toBeInTheDocument()
		expect(logo).toBeVisible()
		expect(logo).toHaveAttribute("src","Technology_icon-Red_Hat-Ansible_Automation_Platform-Standard-RGB.svg")
	})
})
