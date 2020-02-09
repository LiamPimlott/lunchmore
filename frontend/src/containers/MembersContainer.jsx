import React from 'react';
import { useMemberState } from '../hooks';
import Members from '../pages/Members/index.js';

const MembersContainer = () => {
  const members = useMemberState();
  return <Members members={members} />;
}

export default MembersContainer;
