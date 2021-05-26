export const state = () => ({
  offset: 0,
  limit: 20,
  sort: 'relevance',
  searchTerms: '',
  contentTypes: [],
  contentTypeList: [],
  disciplines: [],
  disciplineList: [],
  drive: null,
  doneCreating: 0, //0 = working on it (or not started), -1 = failed creation, 1 = successful creation
  longDriveName: '',
  partitionName: '',
  partitionSize: 0,
  configuration: '',

  pubYearStart: 1700,
  pubYearEnd: 2020,

  newSearch: true,
  newSearchCounts: false,
  pageNo: 1,
  searchReq: {},
  searchResp: {},

  admin: false,
  reqs: JSON.parse(localStorage.getItem('requests')) || [],
  token: '',

  showCart: false,
  helpTab: 0
})

export const mutations = {
  setOffset(state, offset) { state.offset = offset },
  setLimit(state, limit) { state.limit = limit },
  setSort(state, sort) { state.sort = sort },
  setSearchTerms(state, searchTerms) { state.searchTerms = searchTerms},
  setContentTypes(state, contentTypes) { state.contentTypes = contentTypes},
  setContentTypeList(state, contentTypeList) { state.contentTypeList = contentTypeList } ,
  setDisciplines(state, disciplines) { state.disciplines = disciplines },
  setDisciplineList(state, disciplineList) { state.disciplineList = disciplineList },
  setDrive(state, drive) { state.drive = drive},
  setDoneCreating(state, doneCreating) { state.doneCreating = doneCreating },
  setLongDriveName(state, longDriveName) { state.longDriveName = longDriveName },
  setPartitionName(state, partitionName) {state.partitionName = partitionName },
  setPartitionSize(state,partitionSize) {state.partitionSize = partitionSize },
  setConfiguration(state, configuration) { state.configuration = configuration },
  setPubYearStart(state, pubYearStart) { state.pubYearStart = pubYearStart },
  setPubYearEnd(state, pubYearEnd) { state.pubYearEnd = pubYearEnd },

  setNewSearch(state, newSearch) { state.newSearch = newSearch },
  setNewSearchCounts(state, newSearchCounts) { state.newSearchCounts = newSearchCounts },
  setPageNo(state, pageNo) { state.pageNo = pageNo },
  setSearchReq(state, searchReq) { state.searchReq = searchReq },
  setSearchResp(state, searchResp) { state.searchResp = searchResp },

  setAdmin(state, loggedin) { state.admin = loggedin },
  setReqs(state, reqs) {
    state.reqs = reqs
    console.log('saving requests')
    localStorage.setItem('requests', JSON.stringify(reqs));
  },
  setToken(state, token) { state.token = token},

  setShowCart(state, showCart) { state.showCart = showCart },
  setHelpTab(state, helpTab) { state.helpTab = helpTab },
   reset (state) {
    state.admin = false,
    state.newSearch = true,
      state.newSearchCounts = false,
      state.pageNo = 1,
    state.token = '',
    state.sort = 'relevance',
    state.searchTerms = '',
    state.showCart = false,
    state.drive = '',
      state.doneCreating = false,
      state.partitionSize = 0,
      state.longDriveName = '',
      state.partitionName = '',
      state.configuration = '',
   state.helpTab = 0,
     state.contentTypes = [],
     state.disciplines = [],
       state.pubYearStart = 1700,
       state.pubYearEnd = 2020
  }
}

export const actions = {
  setOffset: ({ commit }, offset) => commit('setOffset', offset),
  setLimit: ({ commit }, limit) => commit('setLimit', limit),
  setSort: ({ commit }, sort) => commit('setSort', sort),
  setSearchTerms: ({ commit }, searchTerms) => commit('setSearchTerms', searchTerms),
  setContentTypes: ({ commit }, contentTypes) => commit('setContentTypes', contentTypes),
  setContentTypeList: ({ commit }, contentTypeList) => commit('setContentTypeList', contentTypeList),
  setDisciplines: ({ commit }, disciplines) => commit('setDisciplines', disciplines),
  setDisciplineList: ({ commit }, disciplineList) => commit('setDisciplineList', disciplineList),
  setDrive: ({ commit }, drive) => commit('setDrive', drive),
  setDoneCreating: ({ commit }, doneCreating) => commit('setDoneCreating', doneCreating),
  setLongDriveName: ({ commit }, longDriveName) => commit('setLongDriveName', longDriveName),
  setPartitionName: ({ commit }, partitionName) => commit('setPartitionName', partitionName),
  setPartitionSize: ({ commit }, partitionSize) => commit('setPartitionSize', partitionSize),
  setConfiguration: ({ commit }, configuration) => commit('setConfiguration', configuration),
  setPubYearStart: ({ commit }, pubYearStart) => commit('setPubYearStart', pubYearStart),
  setPubYearEnd: ({ commit }, pubYearEnd) => commit('setPubYearEnd', pubYearEnd),

  setNewSearch: ({ commit }, newSearch) => commit('setNewSearch', newSearch),
  setNewSearchCounts: ({ commit }, newSearchCounts) => commit('setNewSearchCounts', newSearchCounts),
  setPageNo: ({ commit } , pageNo) => commit('setPageNo', pageNo),
  setSearchReq: ({ commit }, searchReq) => commit('setSearchReq', searchReq),
  setSearchResp: ({ commit }, searchResp) => commit('setSearchResp', searchResp),

  setAdmin: ({ commit }, loggedin) => commit('setAdmin', loggedin),
  setReqs: ({ commit }, reqs) => commit('setReqs', reqs),
  setToken: ({ commit }, token) => commit('setToken', token),

  setShowCart: ({ commit }, showCart) => commit('setShowCart', showCart),
  setHelpTab: ({ commit }, helpTab) => commit('setHelpTab', helpTab),
}

export const getters = {
  offset: state => state.offset,
  limit: state => state.limit,
  sort: state => state.sort,
  searchTerms: state => state.searchTerms,
  contentTypes: state => state.contentTypes,
  contentTypeList: state => state.contentTypeList,
  disciplines: state => state.disciplines,
  disciplineList: state => state.disciplineList,
  drive: state => state.drive,
  doneCreating: state => state.doneCreating,
  longDriveName: state => state.longDriveName,
  partitionName: state => state.partitionName,
  partitionSize: state => state.partitionSize,
  configuration: state => state.configuration,
  pubYearStart: state => state.pubYearStart,
  pubYearEnd: state => state.pubYearEnd,

  newSearch: state => state.newSearch,
  newSearchCounts: state => state.newSearchCounts,
  pageNo: state => state.pageNo,
  searchReq: state => state.searchReq,
  searchResp: state => state.searchResp,

  admin: state => state.admin,
  reqs: state => state.reqs,
  token: state => state.token,

  showCart: state => state.showCart,
  helpTab: state => state.helpTab,
}
