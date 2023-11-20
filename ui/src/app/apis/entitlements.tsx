import { EntitlementsCount } from "./types";

const baseURI = "/api"

export function getEntitlementsCount(): Promise<EntitlementsCount> {
	return fetch(`${baseURI}/azmarketplaceentitlementscount`)
		.then((resp) => {
			if (resp.ok) {
				return resp.json()
			}
			throw new Error(`Unexpected response with status code: ${resp.status}`)
		})
		.then((data) => {
			return new EntitlementsCount(data)
		})
}
