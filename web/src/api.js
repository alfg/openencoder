import auth from './auth';

const root = '/api';

const Endpoints = {
  Version: `${root}/`,
  Stats: `${root}/stats`,
  Health: `${root}/health`,
  Pricing: `${root}/machines/pricing`,

  JobsList: page => `${root}/jobs?page=${page}`,
  Jobs: `${root}/jobs`,
  JobsCancel: id => `${root}/jobs/${id}/cancel`,
  JobsRestart: id => `${root}/jobs/${id}/restart`,

  Machines: `${root}/machines`,
  MachinesRegions: `${root}/machines/regions`,
  MachinesSizes: `${root}/machines/sizes`,
  MachinesVPCs: `${root}/machines/vpc`,
  MachinesId: id => `${root}/machines/${id}`,

  PresetsList: page => `${root}/presets?page=${page}`,
  Presets: `${root}/presets`,
  PresetsId: id => `${root}/presets/${id}`,

  FileList: prefix => `${root}/storage/list?prefix=${prefix}`,

  WorkerQueue: `${root}/worker/queue`,
  WorkerPools: `${root}/worker/pools`,

  Users: `${root}/users`,
  UsersId: id => `${root}/users/${id}`,

  Settings: `${root}/settings`,

  CurrentUser: `${root}/me`,
};

function get(context, url, callback) {
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
}

function post(context, url, data, callback) {
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
}

function update(context, url, data, callback) {
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
}

function del(context, url, callback) {
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
}

export default {
  getVersion(context, callback) {
    return get(context, Endpoints.Version, callback);
  },

  getStats(context, callback) {
    return get(context, Endpoints.Stats, callback);
  },

  getHealth(context, callback) {
    return get(context, Endpoints.Health, callback);
  },

  getPricing(context, callback) {
    return get(context, Endpoints.Pricing, callback);
  },

  getJobs(context, page, callback) {
    return get(context, Endpoints.JobsList(page), callback);
  },

  createJob(context, data, callback) {
    return post(context, Endpoints.Jobs, data, callback);
  },

  cancelJob(context, id, callback) {
    return post(context, Endpoints.JobsCancel(id), {}, callback);
  },

  restartJob(context, id, callback) {
    return post(context, Endpoints.JobsRestart(id), {}, callback);
  },

  getMachines(context, callback) {
    return get(context, Endpoints.Machines, callback);
  },

  getMachineRegions(context, callback) {
    return get(context, Endpoints.MachinesRegions, callback);
  },

  getMachineSizes(context, callback) {
    return get(context, Endpoints.MachinesSizes, callback);
  },

  getMachineVPCs(context, callback) {
    return get(context, Endpoints.MachinesVPCs, callback);
  },

  createMachine(context, data, callback) {
    return post(context, Endpoints.Machines, data, callback);
  },

  deleteMachine(context, id, callback) {
    return del(context, Endpoints.MachinesId(id), callback);
  },

  deleteAllMachines(context, callback) {
    return del(context, Endpoints.Machines, callback);
  },

  getPresets(context, page, callback) {
    return get(context, Endpoints.PresetsList(page), callback);
  },

  createPreset(context, data, callback) {
    return post(context, Endpoints.Presets, data, callback);
  },

  updatePreset(context, data, callback) {
    return update(context, Endpoints.PresetsId(data.id), data, callback);
  },

  getFileList(context, prefix, callback) {
    return get(context, Endpoints.FileList(prefix), callback);
  },

  getWorkerQueue(context, callback) {
    return get(context, Endpoints.WorkerQueue, callback);
  },

  getWorkerPools(context, callback) {
    return get(context, Endpoints.WorkerPools, callback);
  },

  getUsers(context, callback) {
    return get(context, Endpoints.Users, callback);
  },

  updateUser(context, data, callback) {
    return update(context, Endpoints.UsersId(data.id), data, callback);
  },

  getSettings(context, callback) {
    return get(context, Endpoints.Settings, callback);
  },

  updateSettings(context, data, callback) {
    return update(context, Endpoints.Settings, data, callback);
  },

  getCurrentUser(context, callback) {
    return get(context, Endpoints.CurrentUser, callback);
  },

  updateCurrentUser(context, data, callback) {
    return update(context, Endpoints.CurrentUser, data, callback);
  },
};
