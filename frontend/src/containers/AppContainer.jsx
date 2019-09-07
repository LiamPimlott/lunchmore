import React from 'react';
import App from '../App.js';
import { useAuthState } from '../hooks';

const AppContainer = () => {
  const auth = useAuthState();

  return <App globalState={{
    auth
  }}/>;
}

export default AppContainer;
