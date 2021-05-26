<template>

    <v-tabs vertical color="#006179" v-model="selectedHelpTab" class="help-tabs">
      <v-tab index='0' key="0" class="help-tab">
        What is JSTOR?
      </v-tab>

      <v-tab index='1' key='1'>
        How to use JSTOR?
      </v-tab>

      <v-tab index='2' key='2'>
        How to research?
      </v-tab>


      <v-tab-item>
        <div v-html="what" class="student-help"/>
      </v-tab-item>

      <v-tab-item>
        <div v-html="how" class="student-help"/>
      </v-tab-item>

      <v-tab-item>
        <div v-html="research" class="student-help"/>
      </v-tab-item>

    </v-tabs>

</template>

<script>
  import what from '@/static/markdown/what-is-jstor.md'
  import how from '@/static/markdown/how-to-use-jstor.md'
  import research from '@/static/markdown/research-basics.md'
  import { mapGetters, mapActions } from 'vuex'

  export default {
    data: () => ({
      markdownSource: '',
    }),

    computed: { ...mapGetters(['helpTab']),
      selectedHelpTab: {
        get: function() { return this.helpTab },
        set (tab) { this.setHelpTab(tab)}
      },
      what () { return what },
      how () { return how},
      research () { return research}
    },
    methods: {...mapActions([ 'setHelpTab']),


    }
  }
</script>

<style>


  .help-tabs .v-tab {
    font-size: 16px !important;
    font-weight: bold !important;
    text-transform: none !important;
    justify-content: left !important;
    padding: 24px !important;
  }

  .v-tabs-items {
    padding-left: 16px;
    min-height: 600px;
    overflow: scroll;
  }

  .v-slide-group__wrapper {
    background-color: #F7F7F7;
  }

  .v-tabs-slider-wrapper {
    width: 5px !important;
  }

  .student-help h1 {
    font-size: 36px;
  }

  .student-help img {
    max-width: 750px !important;
  }

  img[alt=logo] { width: 300px; }

</style>
