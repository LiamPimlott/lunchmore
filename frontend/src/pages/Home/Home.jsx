import React from 'react';
import { Link } from 'react-router-dom';

const Home = ({ auth }) => {
const { id, first_name } = auth;
return (
    <div>
      <h1>Welcome to Lunchmore{first_name ? ` ${first_name}!` : '!'}</h1>
      {!id
        ? <>
          <div>
            <Link to='/login'>Login</Link>      
          </div>
          <div>
            <Link to='/signup'>Signup</Link>
          </div>
        </>
        : <>
          <div>
            <Link>Logout</Link>
          </div>
          <div>
            <Link to='/members'>Members</Link>
          </div>
          <div>
            <Link to='/schedules'>Schedules</Link>
          </div>
          <div>
            <Link to='/lunches'>Lunches</Link>
          </div>
          <div>
            <Link to='/account'>Account</Link>
          </div>
        </>}
    </div>
  );
};

export default Home;
