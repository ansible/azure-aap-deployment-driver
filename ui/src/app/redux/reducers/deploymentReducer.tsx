import { ActionTypes } from "../contants/action-types";

const initialState = {
    products: [],
    err: null
}
export const deploymentReducer = (state = initialState, {type, payload, error}) => {
    switch(type) {
        case ActionTypes.DEPLOYMENT_STEPS:
                return {...state, products: payload, err: error};
        default:
            return state; 
    }
}
