<template>
  <v-app>
    <div>
      <v-container style="min-height: calc(100vh - 350px);">
      <v-row 
        align="center"
        justify="center"
        no-gutters
      >
        <v-col cols="auto">
          <h1>{{ error.message || "An error occurred"}}</h1>
          <p>If this error persists, please contact an administrator.</p>
          <NuxtLink to="/">Home page</NuxtLink>
        </v-col>
      </v-row>
      </v-container>
    </div>
  </v-app>
</template>

<script>
  export default {
    props: ['error'],

  }
</script>

<style>

  html, body {
    font-size: 14px;
    box-sizing: border-box;
  }

  .container {
    max-width: 1120px;
  }

  *,
  *:before,
  *:after {
    box-sizing: border-box;
    margin: 0;
  }

  a {
    text-decoration: underline;
    font-weight: bold;
  }

  h1 {
    font-size: 36px;
    padding: 0;
    margin: 0;
    line-height: 1;
  }

  h3 {
    margin-top:24px;
  }

  .theme--light.v-application {
    background-color: #ffffff  !important;
    color: black;
  }

  .v-application {
    font-family: "Arial", sans-serif;
    line-height: 2;
    color: black;
  }

  .theme--light.v-label {
    color: black;
  }

  fieldset {
     border-style: none;
  }

  .v-btn {
    letter-spacing: normal;
    font-weight: normal;
    text-transform: none;
    height: 40px;
    font-size: 14px !important;
  }

  .v-input {
    font-size: 14px;
    max-width: unset;
  }

  /*.v-input__slot {
    box-shadow: inset 1px 0 10px -3px #ccc;
  }*/

  .v-label {
    font-size: 14px;
  }

  .v-pagination__navigation {
    box-shadow: none !important;
  }

  .v-pagination__item {
    box-shadow: none !important;
  }

  .v-input--selection-controls__input {
    align-self: flex-start;
  }

  .header-line {
    width: 100vw;
    margin-left: -50vw;
    left: 50%;
    position: relative;
    margin-top: 16px;
    border-bottom: 1px solid #dcdcdc;
  }

</style>
