import jwtDecode from 'jwt-decode';
import cookie from './cookie';
import store from './store';

const LOGIN_URL = '/api/login';
const REGISTER_URL = '/api/register';
const UPDATE_PASSWORD_URL = '/api/update-password';

export default {

  user: {
    username: null,
    role: null,
    authenticated: false,
  },

  login(context, creds, redirect, callback) {
    context.$http.post(LOGIN_URL, creds).then((data) => {
      cookie.set('token', data.body.token);
      store.setTokenAction(data.body.token);

      this.user.authenticated = true;
      this.user.username = jwtDecode(data.body.token).id;

      if (redirect) {
        context.$router.push({ name: redirect });
        context.$router.go();
      }
    }, (err) => {
      // If password needs to be updated.
      if (err && err.body.message === 'require password reset') {
        context.$router.push({ name: 'update-password' });
        context.$router.go();
      }
      callback(err);
    });
  },

  register(context, creds, redirect, callback) {
    context.$http.post(REGISTER_URL, creds).then(() => {
      // TODO: Authenticate after registration?
      // cookie.set('token', data.body.token);
      // store.setTokenAction(data.body.token);

      // this.user.authenticated = true;
      // this.user.username = jwtDecode(data.body.token).id;

      if (redirect) {
        context.$router.push({ name: redirect });
      }
    }, (err) => {
      callback(err);
    });
  },

  updatePassword(context, creds, redirect, callback) {
    context.$http.post(UPDATE_PASSWORD_URL, creds).then(() => {
      if (redirect) {
        context.$router.push({ name: redirect });
      }
    }, (err) => {
      callback(err);
    });
  },


  logout(context) {
    cookie.remove('token');
    this.user.authenticated = false;
    context.$router.push({ name: 'login' });
    context.$router.go();
  },

  checkAuth(context) {
    const jwt = cookie.get('token');

    if (context.$route.name === 'register' || context.$route.name === 'update-password') {
      return;
    }

    // Check if token exists from cookie and set the store.
    // If not, then redirect to the login page to get a new token.
    if (jwt && !this.isExpired(jwt)) {
      store.setTokenAction(jwt);
      this.user.authenticated = true;
      this.user.username = jwtDecode(jwt).id;
      this.user.role = jwtDecode(jwt).role;
    } else if (context.$route.name !== 'login') {
      context.$router.push({ name: 'login' });
    }
  },

  isExpired(jwt) {
    return Date.now() >= jwtDecode(jwt).exp * 1000;
  },

  getAuthHeader() {
    return {
      Authorization: `Bearer ${store.state.token}`,
    };
  },
};
