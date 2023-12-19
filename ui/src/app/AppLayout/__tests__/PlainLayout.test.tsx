import React from 'react';
import { render, screen } from '@testing-library/react';
import { PlainLayout } from '../PlainLayout';

describe('AppLayout', ()=>{
	it('renders Page with proper AAP logo', ()=> {
		render(<PlainLayout />)

		const logo = screen.getByAltText("Red Hat Ansible Automation Platform Logo")
		expect(logo).toBeInTheDocument()
		expect(logo).toBeVisible()
		expect(logo).toHaveAttribute("src","Technology_icon-Red_Hat-Ansible_Automation_Platform-Standard-RGB.svg")
	})
})
