import { useState } from 'react';
import axios from 'axios';

const useAuthState = () => {
  const [auth, setAuth] = useState({token: ''});

  const actions = {
    login: async (email, password) => {
      try {
        const r = await axios.post('/users/login', { email, password });
        setAuth({ ...auth, token: r.data.token });
      } catch(err) {
        return err.response.data.message
      }
    },
  };

  return { 
    values: auth,
    actions,
  };
}

export default useAuthState;
