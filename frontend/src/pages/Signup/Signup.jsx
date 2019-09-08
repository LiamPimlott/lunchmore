import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';

const useStyles = makeStyles(theme => ({
  formContainer: {
    padding: '25px',
    borderStyle: 'solid',
    borderWidth: '1px',
  },
}));

const Signup = ({ auth, history }) => {
  const classes = useStyles();
  const [values, setValues] = useState({
    org_name: '',
    first_name: '',
    last_name: '',
    email: '',
    password: '',
  });

  const handleChange = name => event => {
    setValues({ ...values, [name]: event.target.value });
  };

  const onSubmit = async (event) => {
    const err = await auth.actions.signup({ ...values });
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
        <h1>Signup</h1>
        <form>
          <Grid item xs={12}>
            <TextField
              label="Organization Name"
              value={values.org_name}
              onChange={handleChange('org_name')}
              margin="normal"
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              label="First Name"
              value={values.first_name}
              onChange={handleChange('first_name')}
              margin="normal"
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              label="Last Name"
              value={values.last_name}
              onChange={handleChange('last_name')}
              margin="normal"
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              label="Email"
              value={values.email}
              onChange={handleChange('email')}
              margin="normal"
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              label="Password"
              type='password'
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

export default Signup;
