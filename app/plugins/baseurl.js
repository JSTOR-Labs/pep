export default function({ $axios }) {
  if (process.client) {
    const origin = window.location.origin
    $axios.defaults.baseURL = origin
  }
}
