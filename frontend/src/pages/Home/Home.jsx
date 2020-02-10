import React, { useContext } from 'react';
import { Link } from 'react-router-dom';
import { StoreContext } from '../../contexts/StoreContext';

const Home = () => {
const { auth: { state: { id, first_name }, logout } } = useContext(StoreContext);

const handleLogout = async () => {
  const err = await logout();
  if (err) {
    alert(err);
  }
}

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
            <Link to='/' onClick={handleLogout}>Logout</Link>
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
