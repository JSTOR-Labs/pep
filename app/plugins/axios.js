import api from '../config/api'

const axiosPlugin = (ctx, inject) => {

  // Create a custom axios instance
  const axiosAPI = ctx.$axios.create({
    headers: {
      common: {
        Accept: 'text/plain, application/json'
      }
    }
  })

  // Set baseURL to something different
  axiosAPI.setBaseURL(process.env.baseURL)

  const routes = api(ctx.$axios)
  // Inject to context as $api
  inject('api', routes)

  ctx.$axios.onError((error) => {
    ctx.error({ statusCode: error.response?.status, message: 'An error occurred.' })
  })

}

export default axiosPlugin