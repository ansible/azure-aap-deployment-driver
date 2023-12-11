//import { getSteps } from "../../apis/deployment"
import React from 'react';
import { render, screen } from '@testing-library/react';
import { Deployment } from '../Deployment';
import { DeploymentData } from '../../apis/types';

// TODO This could be moved to __mocks__ folder eventually
// mocking getSteps function on APIs because it is used by Deployment component
jest.doMock('../../apis/deployment', () => {
	const originalDeploymentModule = jest.requireActual('../../apis/deployment')

	return {
    __esModule: true,
    ...originalDeploymentModule,
    getSteps: jest.fn(() => {
			return new Promise((resolve, reject)=>{
				process.nextTick(()=>{
					resolve(new DeploymentData([]))
				})
			})
		})
  }
})

describe('Deployment component', ()=>{
	it('shows Red Hat login modal when flag set to true', ()=>{
		render(<Deployment showLoginDialog={true}/>)

		const modal = screen.getByRole('dialog')
		expect(modal).toBeInTheDocument()
		// no need to test the modal/dialog, that is covered within its own test
	})

	it('does not show Red Hat login modal when flag set to false', ()=>{
		render(<Deployment showLoginDialog={false}/>)

		// using query here because it does not throw an error but returns null instead
		const modal = screen.queryByRole('dialog')
		expect(modal).toBeNull()
	})
})
