import React, { useState, useEffect, useContext } from 'react';
import { useHistory, useRouteMatch } from "react-router-dom";
import { useJoinOrganizationState } from '../../hooks';
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

const Join = () => {
  const classes = useStyles();

  const invite = useJoinOrganizationState();
  
  const [values, setValues] = useState({
    first_name: '',
    last_name: '',
    password: '',
  });

  const history = useHistory();
  
  const match = useRouteMatch('/invite/:code');

  const { auth } = useContext(StoreContext);

  const { params: { code } } = match;
  
  useEffect(() => {
    invite.actions.getOrgName(code);
  }, []);

  const handleChange = name => event => {
    setValues({ ...values, [name]: event.target.value });
  };

  const onSubmit = async (event) => {
    const err = await auth.join(code, { ...values });
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
        <h1>{`Join${invite.orgName ? ` ${invite.orgName}` : ''}`}</h1>
        <form>
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
              Join Organization
            </Button>
          </Grid>
        </form>
      </Grid>
    </Container>
  );
};

export default Join;
