import { ActionTypes } from "../constants/action-types";
export const setDeploymentSteps = (deploymentSteps, error) => {
    if (deploymentSteps !== undefined) {
    return  {
        type: ActionTypes.DEPLOYMENT_STEPS,
        payload: deploymentSteps,
        error: null
    };
}
else {
    return  {
        type: ActionTypes.DEPLOYMENT_STEPS,
        payload: [],
        error: error.message
    };
}
};
