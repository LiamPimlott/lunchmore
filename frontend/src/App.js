import React from 'react';
import { BrowserRouter as Router, Route } from "react-router-dom";

import Home from './pages/Home';
import Login from './pages/Login';
import Signup from './pages/Signup';

const App = ({ globalState }) => {
  return (
    <React.Fragment>
      <Router >
        <div>
          <Route exact path="/" component={Home}/>
          <Route
            path="/login"
            render={props => <Login 
              {...props}
              auth={globalState.auth}
            />}
          />
          <Route
            path="/signup"
            render={props => <Signup 
              {...props}
              auth={globalState.auth}
            />}
          />
        </div>
      </Router>
    </React.Fragment>
  );
}

export default App;
