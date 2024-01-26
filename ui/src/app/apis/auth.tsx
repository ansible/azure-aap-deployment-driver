const baseURI = "/api"

export function login(loginData: ILoginData) :Promise<IAuthResponse>  {
	return fetch(`${baseURI}/login`,{
		method: "POST",
		body: JSON.stringify(loginData),
		headers: {
			"Content-Type": "application/json"
		}
	}).then((resp) => {
		// both 200 and 401 should come back with JSON
		if (resp.ok) {
			return resp.json()
		}
		else if (resp.status === 401){
			return ({status:"", error:`Incorrect Password`} as IAuthResponse);
		}
		else if (resp.status == 502){
			return ({status:"unavailable", error:`Deployment Driver server unavailable`} as IAuthResponse);
		}
		return ({status:"", error:`Unexpected response with status code: ${resp.status}`} as IAuthResponse);
	}).catch((err) =>{
		return ({status:"", error: err.message} as IAuthResponse);
	})
}

export function logout(): Promise<IAuthResponse> {
	return fetch(`${baseURI}/logout`,{
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		}
	}).then((resp) => {
		if (resp.ok) {
			return resp.json()
		}
		return ({status:"",error:`Unexpected failure with status code: ${resp.status}`} as IAuthResponse)
	}).catch(err =>{
		return ({status:"", error: err.message});
	})
}

export interface ILoginData {
	uid: string
	pwd: string
}

export interface IAuthResponse {
	status: string
	error?: string
}
