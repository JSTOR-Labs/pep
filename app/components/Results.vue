<template>
  <div id="results">

    <v-row class="results-bar">
      <v-col md="12">
        <h2 style="margin-top: -8px">Showing {{(pageNo - 1 ) * limit + 1}} - {{(pageNo ) * limit}} of  {{new Intl.NumberFormat('en-US').format(searchResp.total) || 0}} search results </h2>
        <!--<p>{{ resultDescription }}</p>-->
      </v-col>

      <!--<v-col align="right" md="3" class="pa-1">
          <v-select
            outlined
            class="sort"
            label="Sort by: Relevance"
            dense
          />
      </v-col>-->
      </v-row>

    <div v-for="doc in searchResp.docs || []" :key="doc['_id']">
      <v-row>
        <v-col md="9">
          <div style="line-height: 1.7">
            <span class="content-type"> {{ doc.contentType == 'book' ? 'Book Chapter' : doc.contentType}} </span>
            <span v-if="doc.pageCount">({{ doc.pageCount }} pages) </span><br>
            <span class="article-title"> {{doc.title}} </span> <br>
            <span style="line-height: 30px" v-if="doc.authors">Author: {{doc.authors[0]}}</span><!-- todo format author string? -->
            <p v-html="doc.srcHtml"/>
            <p v-if="doc.semanticTerms">
              <b>Topics: </b>
              <span v-for="(topic, key) in doc.semanticTerms" :key = 'key'>
                <a @click="searchFor(topic)">{{topic}}</a><span v-if="key+1 != doc.semanticTerms.length">, </span>
              </span>
            </p>
            <p v-if="doc.ocr">
              {{ doc.ocr }}
            </p>

          </div>
        </v-col>

        <v-col md="3" align="right" class="results-button">
          <v-btn
            depressed
            dark
            full-width:true
            color="primary"
            v-if="reqs.includes(JSON.stringify(doc))"
            title="Click again to remove"
            @click="removeRequest(doc)"
          >
            Added to cart
          </v-btn>

          <v-btn
            outlined
            dark
            color="primary"
            v-else
            title="Click to request"
            @click="addRequest(doc)">
            Request this
          </v-btn>
        </v-col>
      </v-row>
      <hr color="#DCDCDC" class="hr-results" size="1px">
    </div>
    <v-pagination
      class="pagination"
      depressed
      v-model="pageNo"
      @input="onPageChange"
      :total-visible="6"
      :length="numPages || 0"
    />
  </div>
</template>

<script>
  import { mapGetters, mapActions } from 'vuex'
  import manageCart from '@/mixins/manageCart'
  import searchConstructor from '@/mixins/searchConstructor'

  export default {
    name: "Results",
    computed: {...mapGetters(['searchResp', 'reqs', 'limit', 'searchReq']),
      numPages() {
          return Math.ceil(this.searchResp.total / this.limit)
      },
      resultDescription() {
        console.log('result!! ', this.searchReq)
        return this.searchReq.toString()
      }
    },
    data: () => ({
      pageNo: 1,
      success: 'Request submitted',
    }),
    mixins: [ manageCart, searchConstructor ],
    mounted() {
      this.doSearch(true)
    },

    methods: {...mapActions(['setSearchResp', 'setAdmin', 'setReqs', 'setLimit', 'setOffset', 'setSearchTerms']),
      onPageChange() {
        this.setOffset((this.pageNo - 1) * this.limit)
        this.doSearch(false)
      },
      searchFor(topic) {
        this.setSearchTerms(topic)
        this.doSearch(true);
      },

    }
  }
</script>

<style scoped>

  .results-bar {
    border-bottom: 1px solid #c9c9c9;
    margin-bottom: 16px;
  }

  .pub-date .theme--light.v-label {
    color: #8d8d8d;
  }

  .hr-results {
    margin-bottom: 24px;
  }

  .pagination {
    margin-top: 16px;
  }

  .content-type {
    text-transform: uppercase;
    font-size: 12px;
    font-weight: bold;
    line-height: 0;
  }

  .article-title {
    font-size: 16px;
    color: #006179;
    font-weight: bold;
  }

  .results-button {
    padding-right: 4px;
  }

</style>

<style>

  .sort .v-input__control .v-input__slot fieldset {
    border-color: #006179;
    color: #006179;
  }

  .sort .v-input__control .v-input__slot:before {
    border: none;
  }

  .sort .v-input__control .v-input__slot:after {
    border: none;
  }

  .sort .v-text-field__details {
    display: none;
  }

  .sort .v-select__slot label, .sort .v-select__slot .v-icon {
    color: #006179;
  }

</style>
