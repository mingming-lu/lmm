import * as https from 'https'

export default function({ $axios, isDev }) {
  $axios.onRequest(config => {
    if (isDev) {
      config.httpsAgent = new https.Agent({ rejectUnauthorized: false })

      if (process.server) {
        config.headers.common['Host'] = process.env.API_HOST
      }
    }
  })
}
