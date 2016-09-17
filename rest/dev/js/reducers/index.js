import { combineReducers } from "redux";

import overview from "js/reducers/overviewReducer";
import log from "js/reducers/logReducer";

export default combineReducers({
  overview,
  log
});
