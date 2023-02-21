export default function ({ store, redirect, route }) {
  // If the user is not authenticated
  if (!store.state.admin && route.path!=="/") {
    return redirect('/')
  }
}
