import { applyMiddleware, createStore } from "redux";

import thunk from "redux-thunk";
import reducer from "./reducers";
import logger from "redux-logger";

// const middleware = applyMiddleware(promiseMiddleware(), thunk);

// logger() will log every state change
const middleware = applyMiddleware(thunk, logger());
export default createStore(reducer, middleware);
