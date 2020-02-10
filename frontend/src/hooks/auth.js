import { useReducer, useEffect } from 'react';
import axios from 'axios';
import { authReducer, authInitialState } from '../reducers/auth';
import ACTIONS from '../actions/auth';

axios.defaults.baseURL = process.env.REACT_APP_API_HOST;

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
    logout: async () => {
      try {
        await axios.get('/users/logout');
        dispatch({ type: ACTIONS.LOGOUT })
      } catch(err) {
        return err.response.data.message
      }
    },
    refresh: async () => {
      try {
        const r = await axios.get('/users/refresh');
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

  // try to refresh once on load
  useEffect(() => {
    functions.refresh();
  }, []);

  return { state, ...functions };
}

export default useAuth;
