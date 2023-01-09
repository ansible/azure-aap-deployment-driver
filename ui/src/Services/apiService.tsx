
export class ApiService {

    public async getSteps(): Promise<any> {
        const response = await fetch('/api/step');
        return await response.json();
    }

    public async restartStep(id: any) {
        const response = await fetch(`/api/restart`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id })
        })
        return await response.json();
    }

    public async cancelDeployment() {
        const response = await fetch(`/api/cancel`, {
            method: 'POST'
        })
        return await response.json();
    }

}