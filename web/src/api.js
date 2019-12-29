import auth from './auth';


export default {
  getStats(context, callback) {
    const url = '/api/stats';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        console.log(err);
        callback(err);
      });
  },

  getHealth(context, callback) {
    const url = '/api/health';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        console.log(err);
        callback(err);
      });
  },

  getPricing(context, callback) {
    const url = '/api/machines/pricing';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        console.log(err);
        callback(err);
      });
  },

  getJobs(context, page, callback) {
    const url = `/api/jobs?page=${page}`;

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        console.log(err);
        callback(err);
      });
  },

  createJob(context, data, callback) {
    const url = '/api/jobs';

    context.$http.post(url, data, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  cancelJob(context, id, callback) {
    const url = `/api/jobs/${id}/cancel`;

    context.$http.post(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  restartJob(context, id, callback) {
    const url = `/api/jobs/${id}/restart`;

    context.$http.post(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getMachines(context, callback) {
    const url = '/api/machines';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getMachineRegions(context, callback) {
    const url = '/api/machines/regions';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getMachineSizes(context, callback) {
    const url = '/api/machines/sizes';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  createMachine(context, data, callback) {
    const url = '/api/machines';

    context.$http.post(url, data, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  deleteMachine(context, id, callback) {
    const url = `/api/machines/${id}`;

    context.$http.delete(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  deleteAllMachines(context, callback) {
    const url = '/api/machines';

    context.$http.delete(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getPresets(context, page, callback) {
    const url = `/api/presets?page=${page}`;

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  createPreset(context, data, callback) {
    const url = '/api/presets';

    context.$http.post(url, data, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  updatePreset(context, data, callback) {
    const url = `/api/presets/${data.id}`;

    context.$http.put(url, data, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getS3List(context, prefix, callback) {
    const url = `/api/s3/list?prefix=${prefix}`;

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getWorkerQueue(context, callback) {
    const url = '/api/worker/queue';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getWorkerPools(context, callback) {
    const url = '/api/worker/pools';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getSettings(context, callback) {
    const url = '/api/settings';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  updateSettings(context, data, callback) {
    const url = '/api/settings';

    context.$http.put(url, data, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getUsers(context, callback) {
    const url = '/api/users';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  updateUser(context, data, callback) {
    const url = `/api/users/${data.id}`;

    context.$http.put(url, data, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  getCurrentUser(context, callback) {
    const url = '/api/me';

    context.$http.get(url, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

  updateCurrentUser(context, data, callback) {
    const url = '/api/me';

    context.$http.put(url, data, {
      headers: auth.getAuthHeader(),
    })
      .then(response => (
        response.json()
      ))
      .then((json) => {
        callback(null, json);
      })
      .catch((err) => {
        callback(err);
      });
  },

};
