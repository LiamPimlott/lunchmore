import React from 'react';
import { useMemberState } from '../hooks';
import Members from '../pages/Members/index.js';

const MembersContainer = ({ auth }) => {
  const members = useMemberState();
  return <Members members={members} auth={auth} />;
}

export default MembersContainer;
