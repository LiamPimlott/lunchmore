import { useState } from 'react';
import axios from 'axios';

const useAuthState = () => {
  const [auth, setAuth] = useState({
    id: '',
    email: '',
    org_id: '',
    first_name: '',
    last_name: '',
    token: '',
  });

  const actions = {
    login: async (email, password) => {
      try {
        const r = await axios.post('/users/login', { email, password });
        setAuth({
          ...auth,
          id: r.data.id,
          email: r.data.email,
          org_id: r.data.org_id,
          first_name: r.data.first_name,
          last_name: r.data.last_name,
          token: r.data.token,
        });
      } catch(err) {
        return err.response.data.message
      }
    },
    signup: async (form) => {
      try {
        const r = await axios.post('/signup', { ...form });
        setAuth({
          ...auth,
          id: r.data.id,
          email: r.data.email,
          org_id: r.data.org_id,
          first_name: r.data.first_name,
          last_name: r.data.last_name,
          token: r.data.token,
        });
      } catch(err) {
        return err.response.data.message
      }
    },
  };

  return { 
    ...auth,
    actions,
  };
}

export default useAuthState;
