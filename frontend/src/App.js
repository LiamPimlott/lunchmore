import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import { StoreProvider } from './contexts/StoreContext'
import Home from './pages/Home';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Join from './pages/Join';
import MembersContainer from './containers/MembersContainer';

const App = () => {
  return (
    <StoreProvider>
      <Router>
        <Switch>
          <Route path='/login'>
            <Login />
          </Route>
          <Route path='/signup'>
            <Signup />
          </Route>
          <Route path='/members'>
            <MembersContainer />
          </Route>
          <Route path='/invite/:code'>
            <Join />
          </Route>
          <Route exact path='/'>
            <Home />
          </Route>
        </Switch>
      </Router>
    </StoreProvider>
  );
}

export default App;
