export default {
  debug: true,
  state: {
    token: null,
  },

  setTokenAction(newValue) {
    if (this.debug) console.log('setTokenAction triggered with', newValue);
    this.state.token = newValue;
  },
};
