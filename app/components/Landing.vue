<template>
  <div>
    <div class="hero">
      <v-container class="landing-hero-container" >
        <v-row>
          <nuxt-link to="/">
            <img src="../static/J_horizontal_white.svg" width="150" class="landing-logo" >
          </nuxt-link>
          <v-col align="right">
            <nuxt-link
              to="#help"

              class="admin-login-landing"
            >Student Help</nuxt-link>

            <v-menu
              id="login-prompt"
              offset-y
              :close-on-content-click="false"
            >
              <template v-slot:activator="{ on }">
                <a
                  class="admin-login-landing"
                  v-on="on"
                > Admin Login
                </a>
              </template>
              <div class="enter-admin-pwd">
                <h3 style="margin-top: 0; margin-bottom: 8px"> Enter admin password</h3>
                <v-text-field
                  id="password"
                  type="password"
                  v-model="password"
                  outlined
                  dense
                  @keydown.enter="checkPassword"
                > </v-text-field>
                <v-btn
                  @click="checkPassword"
                  depressed
                  dark
                  color="primary"
                > Login
                </v-btn>
            </div>
          </v-menu>
          </v-col>
        </v-row>
            <div style="position:relative;height: 336px"  >
              <div style="position: absolute; top: 50%; transform: translateY(-50%);width: 100%; ">

                  <h1 style="text-align: center; color: white; padding-bottom: 8px; font-size: 48px"> Welcome to JSTOR!</h1>

                  <v-col >
                    <div class="landing-search-container">
                      <v-text-field
                        class="search-box-landing"
                        v-model="searchterms"
                        outlined
                        background-color="white"
                        placeholder="Search JSTOR"
                        @keydown.enter="search()"
                      />
                      <v-btn
                        depressed
                        dark
                        class="search-button-landing"
                        color="primary"
                        @click="search()"
                      ><v-icon>search</v-icon>
                      </v-btn>
                    </div>
                  </v-col>
                </div>

            </div>


      </v-container>
    </div>

    <h1 class="landing-h1" >How to use this tool</h1>

    <v-row>
      <v-col cols="4">
        <img  class="landing-circle-image" src="../static/search1.png">
        <div>
          <h2 class="landing-h2">Search</h2>
          <p class="landing-description">Start your search by entering the term or terms that you’re researching.</p>
          <a @click="navTo('keyword')">Tips for searching using keywords</a>
        </div>
      </v-col>

      <v-col cols="4">
        <img class="landing-circle-image" src="../static/filter2.png">
        <div>
          <h2 class="landing-h2">Filter</h2>
          <p class="landing-description">Narrow down your results by date, subject & content type.</p>
            <a @click="navTo('filter')">Tips for filtering</a>
        </div>
      </v-col>

      <v-col cols="4">
        <img class="landing-circle-image" src="../static/request3.png">
        <div>
          <h2 class="landing-h2">Request</h2>
          <p class="landing-description">When you find articles or chapters you’re interested in, request them.</p>
            <a @click="navTo('request')">Tips for making and managing your requests</a>
        </div>
      </v-col>
    </v-row>
    <br>
    <hr color="#DCDCDC" class="hr-results" size="1px">
    <br>

    <h1 class="landing-h1" id="help">Help for students</h1>
    <v-row>
      <v-col cols="12" class="help-viewer">
        <student-help/>
      </v-col>
    </v-row>

  </div>
</template>

<script>

  import { mapGetters, mapActions } from 'vuex'
  import searchConstructor from '@/mixins/searchConstructor'
  import StudentHelp from '@/components/StudentHelp'
  //import _ from 'lodash'

  export default {
    components: {
      StudentHelp
    },
    data: () => ({
      password: ''
    }),
    computed: {...mapGetters(['admin']),
      searchterms: {
        get: function () { return this.$store.state.searchTerms },
        set: function (val) { this.$store.commit('setSearchTerms', val)}
      },
    },
    mixins: [ searchConstructor ],
    created () {
      //window.addEventListener('scroll', this.handleScroll);
    },
    mounted() {
      document.getElementById('header').classList.add('noheader');

    },
    beforeDestroy() {
      document.getElementById('header').classList.remove('noheader');
      //window.removeEventListener('scroll', this.handleScroll);
    },
    methods: {...mapActions([ 'setAdmin', 'setSearchTerms', 'setToken', 'setHelpTab']),
      async checkPassword() {
        try {
          let resp = await this.$axios.$post("/login", { password: this.password })

          console.log(resp)
          console.log('token ', resp.data.token)
          this.setToken(resp.data.token);
          this.$axios.setToken(resp.data.token, 'Bearer')
          this.setAdmin(true);
        }
        catch(err) {
          console.log('Error in api.js ', err)
          alert('Login failed. Wrong password, probably.')
        }
      },
      search: function() {
        this.setSearchTerms(this.searchterms)
        this.doSearch(true)
        this.$router.push({path: '/search'})
      },
      navTo: function(location) {
       // console.log('navto ', location)
        if (location === 'filter') {
          this.setHelpTab(2)
        } else {
          this.setHelpTab(1)
        }
        setTimeout(function(){
          document.getElementById(location).scrollIntoView({ behavior: 'smooth'})
        }, 500);
      }
    }
  }
</script>

<style>

  .hero {
    background-image:url("../static/PEP_hero.jpg");
    height: 400px;
    background-size: cover;
    width: 100vw;
    margin-left: -50vw;
    left: 50%;
    position: relative;
    margin-bottom: 24px;
    margin-top: -12px;
  }

  .landing-hero-container {
    height: 400px;
  }

  .admin-login-landing {
    color: white !important;
    margin-bottom: 8px;
    margin-left: 40px;
  }

  .v-menu__content {
    margin-top: 4px;
    background-color: white;
  }

  .enter-admin-pwd {
    padding: 16px;
  }

  .landing-search-container {
    margin: auto;
    width: 660px;
    border-radius: 4px;
    box-shadow: 0 3px 2px 0 rgba(0,0,0,0.40);
    height: 56px
  }

  .search-box-landing {
    width: 598px;
    margin: auto!important;
    display: inline-block;
  }

  .search-button-landing {
    height: 56px !important;
    margin-top: -2px;
    margin-left: -6px;
    border-radius: 0px 4px 4px 0px;
    display: inline-block;
  }

  .landing-circle-image {
    width:320px;
    height: 320px;
    border-radius: 50%;
    border: 1px solid #7c7c7c;
  }

  .landing-h1 {
    margin-bottom: 24px;
    font-family: Georgia,serif;
    font-weight: normal;
  }

  .landing-h2 {
    font-size: 24px;
    font-weight: normal;
    margin-top: 8px;
    margin-bottom: 8px;
  }

  .landing-description {
    line-height: 1.6;
    font-size: 16px;
    margin-bottom: 4px !important;
  }

  .help-viewer {
    min-height: 600px;
  }

</style>
