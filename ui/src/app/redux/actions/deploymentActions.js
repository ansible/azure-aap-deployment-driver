import { ActionTypes } from "../contants/action-types";
export const setDeploymentSteps = (products, error) => {
    if (products !== undefined) {
    return  {
        type: ActionTypes.DEPLOYMENT_STEPS,
        payload: products,
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
