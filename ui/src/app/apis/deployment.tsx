import { DeploymentData } from "./types"

const baseURI = "/api"

export function getSteps(): Promise<DeploymentData> {
	return fetch(`${baseURI}/step`)
		.then((resp) => {
			if (resp.ok) {
				return resp.json()
			}
			throw new Error(`Unexpected response with status code: ${resp.status}`)
		})
		.then((data) => {
			return new DeploymentData(data)
		});
}

export function restartStep(id: number): Promise<boolean> {
	return fetch(`${baseURI}/execution/${id}/restart`,
		{
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ id })
		})
		.then((resp)=>resp.ok);
}

export function cancelDeployment(): Promise<boolean> {
	return fetch(`${baseURI}/cancelAllSteps`,
		{
			method: 'POST',
			headers: { 'Content-Type': 'application/json' }
		})
		.then((resp) => resp.ok);
}
