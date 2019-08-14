import Cookies from 'js-cookie';

export default {
  set(key, value) {
    const options = {
      expires: 7,
    };
    Cookies.set(key, value, options);
  },

  get(key) {
    return Cookies.get(key);
  },

  remove(key) {
    Cookies.remove(key);
  },
};
