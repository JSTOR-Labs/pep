<template>
  <div>
    <div v-if="!admin">
      <landing></landing>
    </div>
    <div v-if="admin">
      <admin-home></admin-home>
    </div>
  </div>


</template>

<script>
  import Landing from '~/components/Landing'

  import AdminHome from '~/components/AdminHome'
  import { mapGetters } from 'vuex'

export default {
  components: {
    Landing,
    AdminHome
  },
  computed: {...mapGetters(['admin'])}

}
</script>

<style>
</style>
