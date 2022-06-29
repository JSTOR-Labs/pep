<template>
  <div id="header" >
    <!--<div v-if="online" class="internet-signal-yes">
      <v-icon class="internet-icon-yes" small color="green">check_circle</v-icon> Connected to the internet
    </div>

    <div v-else class="internet-signal-no">
      <v-icon class="internet-icon-no" small color="black">error</v-icon> Not connected to the internet
    </div>-->

    <div class="menu-bar">
      <nuxt-link to="/">
        <img src="../static/J_horizontal.svg" width="150" style="margin-bottom: -10px;">
      </nuxt-link>

      <p v-if="admin" class="admin-true"> Search appliance - admin</p>

      <div v-if="admin" id="logout">

        <div v-if="online" class="internet-signal-yes">
          <v-icon class="internet-icon-yes" small color="green">check_circle</v-icon> Connected to the internet
        </div>

        <div v-else class="internet-signal-no">
          <v-icon class="internet-icon-no" small color="black">error</v-icon> Not connected to the internet
        </div>

        <div class="admin-menu">

          <nuxt-link to="/">Admin Home</nuxt-link> <v-icon class="icon" color="primary" small>home</v-icon> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          <!-- Admin configuration <a>Admin</a> <v-icon class="icon" color="primary"  small>settings</v-icon> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-->
          <a @click="logout">Logout</a> <v-icon class="icon" color="primary" small>arrow_forward</v-icon>
        </div>
      </div>

      <div v-if="admin" style="margin-bottom: 16px"> </div>
    </div>

    <v-row class="header-row">
      <v-col md="5">
        <div v-if="!admin">
          <v-text-field
            class="search-box"
            v-model="searchterms"
            dense
            outlined
            placeholder="Search JSTOR"
            @keydown.enter="search()"
          />
          <v-btn
            depressed
            dark
            class="search-button"
            color="primary"
            @click="search()"
          ><v-icon>search</v-icon>
          </v-btn>
        </div>
      </v-col>

      <v-col align="right" md="7">
        <div v-if="!admin" id="login">
          <nuxt-link
            :to="{ path: '/',hash:'#help'}"
            class="admin-login-search"
          >Student Help</nuxt-link>

          <v-dialog
            v-model="showCart"
            width="750"
          >
            <template v-slot:activator="{ on }">
              <v-btn
                class="cart-button"
                depressed
                dark
                color="primary"
                v-on="on"
              > <v-icon class="pr-2">mdi-format-list-bulleted</v-icon>   | &nbsp; Your requests ({{ reqs.length }} items)
              </v-btn>
            </template>
            <student-request @clicked="showSnackbar"/>
          </v-dialog>
          <v-snackbar
            v-model="snackbar"
            :timeout=2000
          >
            Your requests were submitted successfully
            <v-btn
              color="blue"
              text
              @click="snackbar = false"
            >
              Close
            </v-btn>
          </v-snackbar>
        </div>
      </v-col>
    </v-row>



  </div>
</template>

<script>
  import { mapGetters, mapActions } from 'vuex'
  import manageCart from '@/mixins/manageCart'
  import searchConstructor from '@/mixins/searchConstructor'
  import StudentRequest from '~/components/StudentRequest'

  //import { search } from '~/api/search'

  export default {
    name: 'AppHeader',
    components: { StudentRequest },
    computed: {...mapGetters(['searchResp', 'admin', 'reqs', 'showCart']),
         searchterms: {
             get: function () { return this.$store.state.searchTerms },
             set: function (val) { this.$store.commit('setSearchTerms', val)}
         },
        showCart: {
            get: function () { return this.$store.state.showCart },
            set: function(val) { this.$store.commit('setShowCart', val)}
        }
    },
    mixins: [ manageCart, searchConstructor ],
    mounted() {
      this.online=navigator.onLine
    },
    data: () => ({
      //
      online: false,
      snackbar: false
    }),
    methods: {...mapActions(['setSearchResp', 'setAdmin', 'setSearchTerms', 'setToken', 'setShowCart', 'setNewSearchCounts']),
      logout() {
        this.setAdmin(false);
        this.setToken('');
        this.$axios.setToken('', 'Bearer')
        this.$router.push('/')
      },
      search: function() {
        this.setSearchTerms(this.searchterms)
        this.setNewSearchCounts(true)
        this.doSearch(true)
      },
      showSnackbar() {
        this.snackbar = true
      }
    }
  }
</script>

<style>
  #login .v-input__slot {
    width: 150px;
  }

  #login .v-input {
    float:left;
  }

  #login button {
    margin-top: 0px;
  }

  .admin-true {
    font-style: italic;
    font-size: 18px;
    display: inline;
    position: relative;
    bottom: -5px;
  }

  .header-row {
    margin-bottom:-36px;
  }

  .menu-bar {
    margin-top: 8px;
  }

  .search-box {
    width: calc(100% - 70px);
    display: inline-block;
  }

  .search-button {
    height: 40px !important;
    margin-top: -2px;
    margin-left: -6px;
    border-radius: 0px 4px 4px 0px;
  }

  .admin-login-search{
    margin-bottom: 8px;
    margin-right: 56px;
  }

  .cart-button {
    height: 40px !important;
  }

  #logout {
    display: inline-block;
    float: right;
    margin-top: 16px;
  }

  .admin-login {
    margin-left: 72px;
    margin-right: 24px;
  }

  .icon {
    text-decoration: none !important;
    padding-left: 2px;
  }

  .internet-signal-yes {
    position: absolute;
    top: 0;
    left: calc(50vw - 80px);
    background-color: #D3F0DF;
    border-radius: 0 0 6px 6px;
    padding: 0 16px;
    font-size: 10px;
  }

  .internet-icon-yes {
    font-size: 8px !important;
    margin-top: -3px;
    margin-right: 3px;
  }

  .internet-signal-no {
    position: absolute;
    top: 0;
    left: calc(50vw - 85px);
    background-color: #D9D9D9;
    border-radius: 0 0 6px 6px;
    padding: 0 16px;
    font-size: 10px;
  }

  .internet-icon-no {
    font-size: 8px !important;
    margin-top: -3px;
    margin-right: 3px;
  }

  .noheader {
    display: none;
  }

  .admin-menu {
    padding-top: 0;
  }



</style>
