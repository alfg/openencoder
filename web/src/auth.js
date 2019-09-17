import jwtDecode from 'jwt-decode';
import cookie from './cookie';
import store from './store';

const LOGIN_URL = '/api/login';
const REGISTER_URL = '/api/register';

export default {

  user: {
    username: null,
    role: null,
    authenticated: false,
  },

  login(context, creds, redirect) {
    context.$http.post(LOGIN_URL, creds).then((data) => {
      cookie.set('token', data.body.token);
      store.setTokenAction(data.body.token);

      this.user.authenticated = true;
      this.user.username = jwtDecode(data.body.token).id;

      if (redirect) {
        context.$router.push({ name: redirect });
      }
    }, (err) => {
      console.log(err);
    });
  },

  register(context, creds, redirect) {
    context.$http.post(REGISTER_URL, creds).then((data) => {
      // TODO: Authenticate after registration?
      // cookie.set('token', data.body.token);
      // store.setTokenAction(data.body.token);

      // this.user.authenticated = true;
      // this.user.username = jwtDecode(data.body.token).id;

      if (redirect) {
        context.$router.push({ name: redirect });
      }
    }, (err) => {
      console.log(err);
    });
  },

  logout(context) {
    cookie.remove('token');
    this.user.authenticated = false;
    context.$router.push({ name: 'login' });
  },

  checkAuth(context) {
    const jwt = cookie.get('token');

    if (context.$route.name === 'register') {
      return;
    }

    // Check if token exists from cookie and set the store.
    // If not, then redirect to the login page to get a new token.
    if (jwt && !this.isExpired(jwt)) {
      store.setTokenAction(jwt);
      this.user.authenticated = true;
      this.user.username = jwtDecode(jwt).id;
      this.user.role = jwtDecode(jwt).role;
    } else {
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
