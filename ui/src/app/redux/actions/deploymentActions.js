import { ActionTypes } from "../contants/action-types";
export const setDeploymentSteps = (products, error) => {
    if (products.hasOwnProperty('data')) {
    return  {
        type: ActionTypes.DEPLOYMENT_STEPS,
        payload: products.data,
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
