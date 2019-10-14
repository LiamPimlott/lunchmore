import { useState } from 'react';
import axios from 'axios';

const useJoinOrganizationState = () => {
  const [state, setState] = useState({
    orgName: "",
  });

  const actions = {
    getOrgName: async (code) => {
      try {
        const res = await axios({
          method: 'get',
          url: '/invite',
          params: {
            code: code
          }
        });
        setState((prevState) => ({
          ...prevState,
          orgName: res.data.org_name
        }));
      } catch(err) {
        return err.response.data.message
      }
    },
  };

  return { 
    ...state,
    actions,
  };
}

export default useJoinOrganizationState;
