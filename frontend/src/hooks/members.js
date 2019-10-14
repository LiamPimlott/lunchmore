import { useState } from 'react';
import axios from 'axios';
import { http } from '../helpers';

const useMemberState = () => {
  const [members, setMembers] = useState({
    members: [],
  });

  const actions = {
    invite: async (email, token) => {
      try {
        await axios({
          method: 'post',
          url: '/invite',
          headers: http.getHeaders(token),
          data: { email },
        });
      } catch(err) {
        return err.response.data.message
      }
    },
    fetch: async (token) => {
      try {
        const r = await axios({
          method: 'get',
          url: '/organization/members',
          headers: http.getHeaders(token),
        });
        setMembers({
          ...members,
          members: r.data.members,
        });
      } catch(err) {
        return err.response.data.message
      }
    },
  };

  return { 
    ...members,
    actions,
  };
}

export default useMemberState;
