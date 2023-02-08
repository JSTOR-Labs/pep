<template>
  <v-container>
    <viewer
      :doi="$route.params.doi"
      :id="$route.params.id"
    />
  </v-container>
</template>

<script>
import Viewer from '@/components/Viewer'

export default {
  components: { 
    Viewer,
  },
  name: 'PDFViewer',
}
</script>

<style>
</style>
