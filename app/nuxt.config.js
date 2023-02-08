
export default {
  ssr: false,
  target: 'static',
  generate: {
    fallback: true,
  },
  /*
  ** Headers of the page
  */
  head: {
    //title: process.env.npm_package_name || '',
    title: 'JSTOR',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: process.env.npm_package_description || '' }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
    ]
  },
  /*
  ** Customize the progress-bar color
  */
  loading: { color: '#fff' },
  /*
  ** Global CSS
  */
  css: [
    'material-design-icons/iconfont/material-icons.css',
  ],
  /*
  ** Plugins to load before mounting the App
  */
  plugins: [
    { src: '~/plugins/baseurl.js' }
    //{ src: '~/plugins/vue-json-csv.js' }
  ],
  /*
  ** vuetify module configuration2
  ** https://github.com/nuxt-community/vuetify-module
  */
  vuetify: {
    customVariables: ['~/assets/variables.scss'],
    theme: {
      themes: {
        light: {
          primary: '#990000',
          secondary: '#000000',
          accent: '#29AB53',
          deny: '#BE0101',
          error: '#FF0000',
        },
      },
    },
  },
  /*
  ** Nuxt.js dev-modules
  */
  buildModules: [
  ],
  /*
  ** Nuxt.js modules
  */
  modules: [
    // Doc: https://axios.nuxtjs.org/usage
    '@nuxtjs/axios',
    '@nuxtjs/vuetify',
    '@nuxtjs/markdownit',
    //'vue-json-csv'
  ],
  /*
  ** Axios module configuration
  ** See https://axios.nuxtjs.org/options
  */
  axios: {
    //baseURL: 'http://labs-pep-go.test.cirrostratus.org/'
    // baseURL: 'http://192.168.1.127:1323/',
    credentials: true
  },
  /*
  ** Build configuration
  */
  build: {
    transpile: [
      'pdfjs-dist',
    ],
    /*
    ** You can extend webpack config here
    */
    extend (config, ctx) {
    }
  }
}
