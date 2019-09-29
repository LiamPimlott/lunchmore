import React from 'react';
import { Link } from 'react-router-dom';

const Home = ({ auth }) => {
const { id, first_name } = auth;
return (
    <div>
      <h1>Welcome to Lunchmore{first_name ? ` ${first_name}!` : '!'}</h1>
      {id
        ? <Link>Logout</Link>
        : <>
          <div>
            <Link to='/login'>Login</Link>      
          </div>
          <div>
            <Link to='/signup'>Signup</Link>
          </div>
        </>}
    </div>
  );
};

export default Home;
