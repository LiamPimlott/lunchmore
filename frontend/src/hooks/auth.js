import { useReducer } from 'react';
import { authReducer, authInitialState } from '../reducers/auth';
import ACTIONS from '../actions/auth';

import axios from 'axios';

const useAuth = () => {
  const [state, dispatch] = useReducer(authReducer, authInitialState);

  const functions = {
    signup: async (form) => {
      try {
        const r = await axios.post('/signup', { ...form });
        dispatch({ type: ACTIONS.SIGNUP, payload: r.data })
      } catch(err) {
        return err.response.data.message
      }
    },
    login: async (email, password) => {
      try {
        const r = await axios.post('/users/login', { email, password });
        dispatch({ type: ACTIONS.LOGIN, payload: r.data })
      } catch(err) {
        return err.response.data.message
      }
    },
    refresh: async () => {
      try {
        const r = await axios.post('/users/refresh',);
        dispatch({ type: ACTIONS.REFRESH, payload: r.data })
      } catch(err) {
        return err.response.data.message
      }
    },
    join: async (code, form) => {
      try {
        const r = await axios.post('/invite/accept', { code, ...form });
        dispatch({ type: ACTIONS.JOIN, payload: r.data })
      } catch(err) {
        return err.response.data.message
      }
    },
  };

  return { state, ...functions };
}

export default useAuth;
