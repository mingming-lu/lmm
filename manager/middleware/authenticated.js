export default function({ isServer, route, req, redirect }) {
  if (isServer && !req) {
    return
  }
  if (route.path === '/login') {
    return
  }
  if (!window.localStorage.getItem('accessToken')) {
    if (route.path === '/logout') {
      redirect('/login')
    }
    redirect(`/login?redirect=${route.path}`)
  }
}
