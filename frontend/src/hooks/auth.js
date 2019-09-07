import { useState } from 'react';
import axios from 'axios';

const useAuthState = () => {
  const [auth, setAuth] = useState({token: ''});

  const actions = {
    login: async (email, password) => {
      try {
        const r = await axios.post('/users/login', { email, password });
        switch (r.status) {
          case 200:
            setAuth({ ...auth, token: r.data.token });
            break;
          case 404:
            return "Email or password is incorrect.";
          default:
            return "An error has occured.";
        }
      } catch(e) {
        // console.error(e)
      }
    },
  };

  return { 
    values: auth,
    actions,
  };
}

export default useAuthState;
