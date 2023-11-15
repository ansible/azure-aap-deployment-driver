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
	const alertContent = screen.getByText(/Your subscription is being entitled and deployed,.+/i)
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
	const alertContent = screen.getByText(/In the meantime, you can manage your subscription.+/i)
	expect(alertTitle).toBeInTheDocument()
	expect(alertTitle).toBeVisible()
	expect(alertContent).toBeInTheDocument()
	expect(alertContent).toBeVisible()
})
