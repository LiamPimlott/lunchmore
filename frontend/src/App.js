import React from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';

import Home from './pages/Home';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Join from './pages/Join';
import MembersContainer from './containers/MembersContainer';

const App = ({ globalState }) => {
  return (
    <React.Fragment>
      <Router >
        <div>
          <Route
            path='/login'
            render={props => <Login 
              {...props}
              auth={globalState.auth}
            />}
          />
          <Route
            path='/signup'
            render={props => <Signup 
              {...props}
              auth={globalState.auth}
            />}
          />
          <Route
            path='/members'
            render={props => <MembersContainer 
              {...props}
              auth={globalState.auth}
            />}
          />
          <Route
            path='/invite/:code'
            render={props => <Join 
              {...props}
              auth={globalState.auth}
            />}
          />
          <Route
            exact
            path='/'
            render={props => <Home 
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
