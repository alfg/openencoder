import cookie from './cookie';
import store from './store';

const LOGIN_URL = '/api/login';
const REGISTER_URL = '/api/register';

export default {

  user: {
    authenticated: false,
  },

  login(context, creds, redirect) {
    context.$http.post(LOGIN_URL, creds).then((data) => {
      cookie.set('token', data.body.token);
      store.setTokenAction(data.body.token);

      this.user.authenticated = true;

      if (redirect) {
        context.$router.push({ name: redirect });
      }
    }, (err) => {
      console.log(err);
    });
  },

  register(context, creds, redirect) {
    context.$http.post(REGISTER_URL, creds).then((data) => {
      cookie.set('token', data.body.token);
      store.setTokenAction(data.body.token);

      this.user.authenticated = true;

      if (redirect) {
        context.$router.push({ name: redirect });
      }
    }, (err) => {
      console.log(err);
    });
  },

  logout() {
    cookie.remove('token');
    this.user.authenticated = false;
  },

  checkAuth() {
    const jwt = cookie.get('token');
    if (jwt) {
      this.user.authenticated = true;
    } else {
      this.user.authenticated = false;
    }
  },

  getAuthHeader() {
    return {
      Authorization: `Bearer ${cookie.get('token')}`,
    };
  },
};
