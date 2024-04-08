import { EngineConfiguration } from "./types";

const baseURI = "/api"

export function getEngineConfiguration(): Promise<EngineConfiguration> {
	return fetch(`${baseURI}/engineconfiguration`)
	.then((resp) => {
		if (resp.ok) {
			return resp.json()
		}
		throw new Error(`Unexpected response with status code: ${resp.status}`)
	})
	.then((data) => {
		return new EngineConfiguration(data)
	});
}
