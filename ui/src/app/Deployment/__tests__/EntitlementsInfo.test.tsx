import React from 'react';
import { render, screen } from '@testing-library/react';
import { EntitlementsCount } from '../../apis/types';
import { EntitlementsInfo } from '../EntitlementsInfo';

test("Subscriptions pending", () => {
	const noEntitlements:EntitlementsCount = {
		count: 0,
		error: ""
	}
	render(<EntitlementsInfo entitlementsCount={noEntitlements}></EntitlementsInfo>)
	const alertTitle = screen.getByText("Your Ansible Automation Platform subscription is pending")
	const alertContent = screen.getByText(/You do not have an Ansible Automation Platform subscription or your subscription is being entitled\..+/i)
	expect(alertTitle).toBeInTheDocument()
	expect(alertTitle).toBeVisible()
	expect(alertContent).toBeInTheDocument()
	expect(alertContent).toBeVisible()
})

test("Several subscriptions", () => {
	const noEntitlements:EntitlementsCount = {
		count: 3,
		error: ""
	}
	render(<EntitlementsInfo entitlementsCount={noEntitlements}></EntitlementsInfo>)
	const alertTitle = screen.getByText("You currently have a subscription to Ansible Automation Platform")
	const alertContent = screen.getByText(/To manage or setup new subscription, .+/i)
	expect(alertTitle).toBeInTheDocument()
	expect(alertTitle).toBeVisible()
	expect(alertContent).toBeInTheDocument()
	expect(alertContent).toBeVisible()
})

test("Coudln't fetch subscriptions", () => {
	const noEntitlements:EntitlementsCount = {
		count: 0,
		error: "Something odd happened"
	}
	render(<EntitlementsInfo entitlementsCount={noEntitlements}></EntitlementsInfo>)
	const alertTitle = screen.getByText("We're temporarily unable to fetch your subscription information")
	const alertContent = screen.getByText(/Click on the following link to enable your Ansible Automation Platform subscription and to access Red Hat support\..+/i)
	expect(alertTitle).toBeInTheDocument()
	expect(alertTitle).toBeVisible()
	expect(alertContent).toBeInTheDocument()
	expect(alertContent).toBeVisible()
})
