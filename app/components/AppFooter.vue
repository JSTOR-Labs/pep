<template>
  <div style="background-color: rgb(232,232,232);" id="footer">
    <v-container>
    <v-row>
      <v-col md="3">
        Brought to you by: <br>
        <img src="../static/Labs_logo.svg" width="150" >
        <br>
      </v-col>

      <v-col md="3">
        With funding from: <br>
        <img src="../static/mellon.svg" width="120" style="margin-bottom: -10px;">
      </v-col>

      <v-col md="6">
        <p class="small-para">
          JSTOR is part of ITHAKA, a not-for-profit organization helping the academic community use digital technologies to preserve the scholarly record and to advance research and teaching in sustainable ways.
          <br> <br>
          ©2000-2020 ITHAKA. All Rights Reserved. JSTOR®, the JSTOR logo, JPASS®, Artstor®, and ITHAKA® are registered trademarks of ITHAKA.
        </p>
      </v-col>
    </v-row>
    </v-container>

    <v-container>
      <hr class="disclaimer">
      <div>
        <b>Your access to JSTOR is provided courtesy of your educational institution and your use is subject to the terms of agreement between JSTOR and that institution.</b>
      </div>
      <div>
        <p class="version">Version: {{versionNo}}<span v-if="hasContent">&nbsp;(PDF Version)</span>, Last updated: {{lastUpdate}}</p>
      </div>


    </v-container>


  </div>
</template>

<script>
export default {
  name: 'AppFooter',
  data: () => ({
    versionNo: 'loading',
    lastUpdate: 'loading',
    hasContent: false,
  }),
  mounted() {
    this.getVersion()

  },
  methods: {
    async getVersion() {

      //console.log('resetPagination? ', resetPagination)
      let version = await this.$api.basic.version()
      if (!(version || {}).data) {
        this.versionNo="Unknown"
        this.lastUpdate="Unknown"
        return
      }
      console.log('version in appFooter: ', version)
      this.versionNo = version.data.version;
      this.hasContent = version.data.has_content;
      this.lastUpdate = (new Date(version.data.lastUpdate)).toLocaleDateString()
    },
  }

}
</script>

<style>
  .small-para {
    font-size: 12px;
    line-height: 1.5;
  }

  .disclaimer {
    padding: 4px;
    position: relative;
    width: 100vw;
    left: 50%;
    margin-left: -50vw;
    margin-bottom: 0;
    text-align: center;
    height: auto;
  }

  .version {
    font-size: 12px;
    margin-bottom: 0 !important;
  }

</style>
