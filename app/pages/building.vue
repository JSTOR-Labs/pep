<template>
  <div>
    <div v-if="admin && doneCreating=== 0">
      <h1>Exporting {{ configuration }} </h1>
      <p>It will take minutes or hours to complete this process, depending on the number and size of indexes and content, and the speed and size of the USB drive.</p>

      <v-progress-linear
        height="32px"
        indeterminate
        background-color="#eeeeee"
        color="accent"
        class="progress"
      ><!--{{ Math.ceil(value) }}%-->
      </v-progress-linear>

      <p class="alert"> Do not remove {{ longDriveName }}!</p>

    </div>
    <div v-if="admin && doneCreating === 1">
      <h1>Success!  </h1>
      <p>{{ longDriveName }} has been configured with {{ configuration }} and is ready to use. </p>

      <v-progress-linear
        height="32px"
        value="100"
        background-color="#eeeeee"
        color="accent"
        class="progress"
      ><!--{{ Math.ceil(value) }}%-->
      </v-progress-linear>

      <p>You can <nuxt-link to="/">go back to create another</nuxt-link>.</p>
    </div>
    <div v-if="admin && doneCreating === -1">
      <h1 class="alert">Unfortunately, the drive creation process failed </h1>

      <p>Please <nuxt-link to="/">go back to try again</nuxt-link> or contact <a href="mailto:labs@ithaka.org">JSTOR Labs</a> if the problem persists.</p>
    </div>

  </div>

</template>

<script>
  import { mapGetters } from 'vuex'
    export default {
      name: "building",
      middleware: 'authenticated',
      computed: {...mapGetters(['admin', 'doneCreating', 'longDriveName', 'configuration']),

      },
    }
</script>

<style scoped>

  .progress {
    border-radius: 3px;
  }

  .alert {
    font-size: 21px;
    color: #ff0000;
    margin-top: 16px;
  }

</style>
