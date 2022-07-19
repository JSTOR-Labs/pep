import throttle from 'lodash/throttle'
import { mapGetters, mapActions } from 'vuex'


export default {
  computed: {
    ...mapGetters([
      'searchResp',
      'searchReq',
      'offset',
      'sort',
      'searchTerms',
      'contentTypes',
      'contentTypeList',
      'disciplines',
      'disciplineList',
      'pubYearStart',
      'pubYearEnd',
      'newSearch',
      'newSearchCounts'
    ]),
    data: () => ({
    }),
    mounted() {

    },
    theFilters() {
      const filters = []

      if (this.disciplines != '' && this.disciplines != undefined) {
        console.log(this.disciplines)
        filters.push('disciplines:' + this.disciplines);
      }

      if (this.contentTypes != '' && this.contentTypes != undefined) {
        console.log(this.contentTypes)
        filters.push('contentType:' + this.contentTypes);
      }
      filters.push('year:[' + this.pubYearStart + ' TO ' + this.pubYearEnd + ']')

      return filters
    },

  },
  methods: {
    ...mapActions(['setSearching', 'clearSearchResults', 'setSearchResp', 'setDisciplineList', 'setContentTypeList', 'setNewSearch', 'setPageNo', 'setNewSearchCounts']),

    async doSearch(resetPagination) {
      this.setSearching(true)
      //console.log('resetPagination? ', resetPagination)
      let args =
          {offset: resetPagination ? 0 : this.offset,
        limit: 20,
        fields: ['title', 'authors', 'srcHtml', 'contentType', 'pageCount', 'semanticTerms', 'ocr', 'firstPage', 'lastPage', 'publisher', 'book_description', 'year', 'abstract', 'series'],
        facets: ['contentType', 'disciplines'],
        filters: this.theFilters,/*  filters: ['disciplineCode:(D1 OR d2)'] this.theFilters,*/
        query: this.searchTerms}
        console.log('search args: ', args)
        this.setSearchResp({
          ...this.searchResp,
          docs: []
        })
      let resp = await this.$axios.$post("/search", args)
          console.log('search response in searchConstructor: ', resp)
          this.setSearchResp(resp)
          //console.log('new search? ', this.newSearch)
          //.log('new search counts? ', this.newSearchCounts)
          if (this.newSearch) { //the original search,
            this.setDisciplineList(resp.facets.disciplines)
            this.setContentTypeList(resp.facets.contentType)
            this.setNewSearch(false)
            //} else if (this.newSearchCounts) { // there's a new search term or new year range. A true filter, not a facet
          } else {
            this.updateDisciplineCounts(resp.facets.disciplines)
            this.updateContentTypeCounts(resp.facets.contentType)
            this.setNewSearchCounts(false)
          }
          if (resetPagination) {
            this.setPageNo(1)
          }
          this.setSearching(false)
    },
    updateDisciplineCounts(resp) {
      let existing = _.cloneDeep(this.disciplineList)
      for (let i = 0; i < existing.length; i++){
        let foundThisOne = false
        for (let j = 0; j < resp.length; j++) {
          if (existing[i].value === resp[j].value) {
            existing[i].count = resp[j].count
            foundThisOne = true;
          }
        }
        if (!foundThisOne) {
            existing[i].count = 0
        }
      }
      this.setDisciplineList(existing)
    },
    updateContentTypeCounts(resp) {
      let existing = _.cloneDeep(this.contentTypeList)
      for (let i = 0; i < existing.length; i++){
        let foundThisOne = false
        for (let j = 0; j < resp.length; j++) {
          console.log('existing ', existing[i].value)
          console.log('resp ', resp[j].value)
          if (existing[i].value === resp[j].value) {
            existing[i].count = resp[j].count
            foundThisOne = true;
          }
        }
        if (!foundThisOne) {
          existing[i].count = 0
        }
      }
     // console.log('existing content types ', existing)
      this.setContentTypeList(existing)
    }
  }
}
