import { combineReducers } from "redux";
import { deploymentReducer } from "./deploymentReducer";

const reducers = combineReducers({
    deployment: deploymentReducer
})
export type RootState = ReturnType<typeof reducers>;
export default reducers;
