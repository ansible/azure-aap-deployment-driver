import { ActionTypes } from "../contants/action-types";
export const setDeploymentSteps = (products) => {
    return  {
        type: ActionTypes.DEPLOYMENT_STEPS,
        payload: products,
    };
};