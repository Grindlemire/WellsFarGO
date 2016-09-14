import 'materialize-css/bin/materialize.css';
import 'materialize-css/bin/materialize.js';

import React from 'react';
import {ReactDOM, render} from 'react-dom';
import { Router, Route, Link, hashHistory, Redirect, IndexRedirect} from 'react-router';


import { Provider } from "react-redux";
import store from "./store"

import Overview from 'js/components/overview';
import Home from 'js/components/home';

const appEl = document.getElementById('app')

render(
    <Provider store={store}>
        <Router history={hashHistory}>
            <Route path="/" component={Home}>
                <IndexRedirect to="/overview"/>
                <Route path="overview" component={Overview}/>
            </Route>
        </Router>
    </Provider>
    ,
    appEl
)
