import React, { createContext } from 'react';
import useAuth from '../hooks/auth'

const StoreContext = createContext();

const StoreProvider = ({ children }) => {
  // Auth state
  const auth = useAuth(); 

  return (
    <StoreContext.Provider value={{ auth }}>
      {children}
    </StoreContext.Provider>
  );
};

export { StoreContext, StoreProvider };
