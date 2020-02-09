import AUTH from '../actions/auth';

const initialState = {
  id: '',
  email: '',
  org_id: '',
  first_name: '',
  last_name: '',
  token: '',
};

const reducer = (state, action) => {
  switch (action.type) {
    case AUTH.SIGNUP:
      return { ...state, ...action.payload }  
    case AUTH.LOGIN:
      return { ...state, ...action.payload }
    case AUTH.REFRESH:
      return { ...state, ...action.payload }
    case AUTH.JOIN:
      return { ...state, ...action.payload }  
    default:
      return state;
  }
};

export {
  initialState as authInitialState,
  reducer as authReducer,
};
