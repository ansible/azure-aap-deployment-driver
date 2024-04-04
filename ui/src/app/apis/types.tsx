export enum StepStatuses {
	PENDING,
	STARTED,
	SUCCEEDED,
	FAILED,
	RESTART_PENDING,
	CANCELED,
	PERM_FAILED
}

export class DeploymentStepStatusData {
	// instance fields
	status: StepStatuses
	duration: string = ""
	message: string = ""
	error: string = ""
	errorDetails: string = ""
	code: string = ""
	correlationId: string = ""
	attempts: number = 0
	lastExecutionId: number = -1

	constructor(executions) {
		if (!executions || executions.length === 0) {
			this.status = StepStatuses.PENDING
		} else {
			this.attempts = executions.length
			const stepExecution = executions[executions.length -1]
			this.lastExecutionId = stepExecution.ID
			// map fields from the API payload into this class fields
			switch (stepExecution.status) {
				case "Started":
				case "Restarted":
					this.status = StepStatuses.STARTED
					break
				case "Restart":  // step marked restart was previously failed
					this.status = StepStatuses.RESTART_PENDING
					break
				case "PermanentlyFailed":
				case "RestartTimedOut":
					this.status = StepStatuses.PERM_FAILED
					break
				case "Failed":
					this.status = StepStatuses.FAILED
					break
				case "Succeeded":
					this.status = StepStatuses.SUCCEEDED
					break
				case "Canceled":
					this.status = StepStatuses.CANCELED
					break
				default:
					this.status = StepStatuses.PENDING
					break
			}
			this.duration = stepExecution.duration
			this.message = stepExecution.message
			this.error = stepExecution.error
			this.errorDetails = stepExecution.errorDetails
			this.code = stepExecution.code
			this.correlationId = stepExecution.correlationId
			/* this.#lastExecutionId= stepExecution.ID */
		}
	}
}

export class DeploymentStepData {
	id: number = -1
	name: string = ""
	status: DeploymentStepStatusData|null = null
	constructor(stepData) {
		if (stepData) {
			this.id = stepData.ID
			this.name = stepData.name.replaceAll("__", " ")
			this.status = new DeploymentStepStatusData(stepData.executions)
		}
	}
}

export class DeploymentProgressData {
	progress:number = 0
	failedStepIds:number[] = []
	failedStepNames: string[] = []
	failedExId:number =-1
	isComplete:boolean = false
	isCanceled:boolean = false
	isPermanentlyFailed:boolean = false
	constructor(steps?:DeploymentStepData[]) {
		if (Array.isArray(steps) && steps.length > 0) {
			const succeeded = steps.reduce(
				(succeededCount, currentStep) => {
					return succeededCount += (currentStep.status && currentStep.status.status === StepStatuses.SUCCEEDED ? 1 : 0 )
				},
				0
			)
			this.progress = (steps.length === 0) ? 0 : Math.floor((succeeded/steps.length)*100)
			const failedSteps:DeploymentStepData[] = steps.filter((step)=> step.status && step.status.status === StepStatuses.FAILED)
			this.failedStepIds = failedSteps.map((failedStep)=>failedStep.id)
			this.failedStepNames = failedSteps.map((failedStep)=>failedStep.name)
			this.isComplete = this.progress === 100
			this.isPermanentlyFailed = steps.some((currentStep) => currentStep.status?.status === StepStatuses.PERM_FAILED)
			this.isCanceled = steps.some((currentStep) => currentStep.status?.status === StepStatuses.CANCELED)
			// TODO Fix this up and probably move somewhere else (not exactly progress related)
			if(failedSteps[0] != null && failedSteps[0].status != null)
			{
				this.failedExId = failedSteps[0].status.lastExecutionId
			}
		}
	}
}

export class DeploymentData {
	steps:DeploymentStepData[] = []
	progress:DeploymentProgressData
	constructor(deploymentData) {
		if (Array.isArray(deploymentData) && deploymentData.length>0) {
			// assuming we got the data already ordered so not ordering the array
			// just mapping deployment data to array of our classes
			this.steps = deploymentData.map((stepData) => new DeploymentStepData(stepData))
			this.progress = new DeploymentProgressData(this.steps)
		} else {
			this.progress = new DeploymentProgressData()
		}
	}
}

export class EntitlementsCount {
	count:number = 0
	error:string = ""
	constructor(entitlementCount) {
		if (entitlementCount) {
			if ('count' in entitlementCount && typeof entitlementCount.count === "number") {
				this.count = entitlementCount.count
			}
			if ('error' in entitlementCount && typeof entitlementCount.error === "string") {
				this.error = entitlementCount.error
				if (this.error.length > 0 ){
					this.count = 0
				}
			}
		}
	}
}

export class EngineConfiguration {
	stepRestartTimeout:number = 0
	overallTimeout:number = 0
	engineExitDelay:number = 0
	autoRetryDelay:number = 0
	stepDeploymentTimeout:number = 0
	stepMaxRetries:number = 0
	constructor(config) {
		if (config) {
			if ('stepRestartTimeoutSec' in config && typeof config.stepRestartTimeoutSec === 'number') {
				this.stepRestartTimeout = config.stepRestartTimeoutSec
			}
			if ('overallTimeoutSec' in config && typeof config.overallTimeoutSec === 'number') {
				this.overallTimeout = config.overallTimeoutSec
			}
			if ('engineExitDelaySec' in config && typeof config.engineExitDelaySec === 'number') {
				this.engineExitDelay = config.engineExitDelaySec
			}
			if ('autoRetryDelaySec' in config && typeof config.autoRetryDelaySec === 'number') {
				this.autoRetryDelay = config.autoRetryDelaySec
			}
			if ('stepDeploymentTimeoutSec' in config && typeof config.stepDeploymentTimeoutSec === 'number') {
				this.stepDeploymentTimeout = config.stepDeploymentTimeoutSec
			}
			if ('stepMaxRetries' in config && typeof config.stepMaxRetries === 'number') {
				this.stepMaxRetries = config.stepMaxRetries
			}
		}
	}
}
