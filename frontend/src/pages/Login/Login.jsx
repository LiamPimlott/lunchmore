import React, { useState, useContext } from 'react';
import { useHistory } from "react-router-dom";
import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import { StoreContext } from '../../contexts/StoreContext'

const useStyles = makeStyles(theme => ({
  formContainer: {
    padding: '25px',
    borderStyle: 'solid',
    borderWidth: '1px',
  },
}));

const Login = () => {
  const classes = useStyles();
  const [values, setValues] = useState({
    email: '',
    password: '',
  });

  const history = useHistory();

  const { auth } = useContext(StoreContext);

  const handleChange = name => e => {
    setValues({ ...values, [name]: e.target.value });
  };

  const onSubmit = async () => {
    const err = await auth.login(values.email, values.password);
    if (err) {
      alert(err);
    } else {
      history.push('/');
    }
  }

  return (
    <Container maxWidth='xs' className={classes.formContainer}>
      <Grid
        container
        direction='column'
        alignItems='center'
      >
        <h1>Login</h1>
        <form>
          <Grid item xs={12}>
            <TextField
              id="standard-name"
              label="Email"
              value={values.email}
              onChange={handleChange('email')}
              margin="normal"
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              id="standard-name"
              label="Password"
              value={values.password}
              onChange={handleChange('password')}
              margin="normal"
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12}>
            <Button
              variant="contained"
              color="primary"
              onClick={onSubmit}
            >
              Log In
            </Button>
          </Grid>
        </form>
      </Grid>  
    </Container>
  );
};

export default Login;
