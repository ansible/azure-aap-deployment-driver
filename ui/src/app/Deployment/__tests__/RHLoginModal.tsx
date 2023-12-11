import React from 'react';
import { render, screen, within } from '@testing-library/react';
import { RHLoginModal } from '../RHLoginModal';

describe('RHLoginModal component', ()=>{
	it('renders the modal dialog when flag is true', ()=>{
		render(<RHLoginModal isModalShown={true}/>)
		// using query here because it does not throw an error but returns null instead
		const modal = screen.queryByRole('dialog')
		expect(modal).toBeInTheDocument()
	})

	it('renders the modal dialog with title, content and one button', ()=>{
		render(<RHLoginModal isModalShown={true}/>)
		// using query here because it does not throw an error but returns null instead
		const modal = screen.getByRole('dialog')
		expect(modal).toBeInTheDocument()
		const title = within(modal).getByText('A valid subscription for Ansible Automation Platform in your Red Hat account is required')
		expect(title).toBeVisible()
		let content
		content = within(modal).getByText('Your Ansible Automation Platform deployment is underway.')
		expect(content).toBeVisible()
		content =  within(modal).getByText('To use Ansible Automation Platform on Azure, you MUST have a valid subscription for Ansible Automation Platform in your Red Hat account.')
		expect(content).toBeVisible()
		content = within(modal).getByText('You can set up your Ansible Automation Platform subscription and your Red Hat account by clicking the button below. You will be redirected back to this page upon successful log in or account creation.')
		expect(content).toBeVisible()
		const button = within(modal).getByRole('link')
		expect(button).toBeVisible()
		expect(button).toHaveAttribute('href','/sso')
		expect(button).toHaveAttribute('target','_self')
	})

	it('does not render when flag is false', ()=>{
		render(<RHLoginModal isModalShown={false}/>)
		// using query here because it does not throw an error but returns null instead
		const modal = screen.queryByRole('dialog')
		expect(modal).toBeNull()
	})
})
