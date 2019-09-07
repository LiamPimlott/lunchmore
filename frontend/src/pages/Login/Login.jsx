import React, { useState } from 'react';
import styled from 'styled-components'
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';


const FormContainer = styled.div`
  padding: 25px;
  border-style: solid;
  border-width: 1px; 
`

const Login = ({ auth, history }) => {
  const [values, setValues] = useState({
    email: '',
    password: '',
  });

  const handleChange = name => event => {
    setValues({ ...values, [name]: event.target.value });
  };

  const onSubmit = event => {
    const err = auth.actions.login(values.email, values.password);
    // if (err) {
    //   alert(err);
    // } else {
    //   history.push('/');
    // }
  }

  return (
    <Container>
      <FormContainer>
        <h1>Login</h1>
        <form>
          <Grid item xs={12}>
            <TextField
              id="standard-name"
              label="Email"
              // className={classes.textField}
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
              // className={classes.textField}
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
      </FormContainer>  
    </Container>
  );
};

export default Login;
