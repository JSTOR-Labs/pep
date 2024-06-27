import { NuxtAxiosInstance } from '@nuxtjs/axios'

export default ($axios) => ({
  admin: {
    request: {
      get: (params) => $axios.get("/admin/request", params),
      update: (vars) => $axios.patch("/admin/request", vars),
    },
    pdf: {
      get: (id) => $axios.get("/admin/pdf/"+ id, { responseType: "blob" })
    }
  },
  basic: {
    version: () => $axios.$get("/api/version"),
    login: (password) => $axios.$post("/api/login", { password }),
    request: (args) => $axios.$post("/api/request", args),
    pdf: {
      get: (id) => $axios.get("/api/pdf/"+ id, { responseType: "blob" })
    },
    search: (args) => $axios.post("/api/search", args),
  },

})
