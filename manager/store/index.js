export const state = () => ({
  accessToken: ''
})

export const mutations = {
  setAccessToken(state, accessToken) {
    state.accessToken = accessToken
  }
}
