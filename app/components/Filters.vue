<template>
  <div id="filters">
    <h2>Refine Results</h2>
    <h3>Content Type</h3>

     <div v-for="(type, index) in contentTypeList" :key="index">
        <v-checkbox
          dense
          color="primary"
          off-icon="check_box_outline_blank"
          on-icon="check_box"
          class="discipline-checkbox"
          v-model="contentTypesChecked"
          :label="contentLabel(index)"
          :value="type.value"
          @change="filterSearch()"
        />

    </div>

    <h3>Publication Date</h3>
    <v-row class="pub-date">
      <v-col md="5">
        <v-text-field
          outlined
          background-color="white"
          dense
          v-model="start"
          placeholder="yyyy"
          :clearable="false"
          :validate-on-blur="true"
          @keydown="typeYear"
        />
      </v-col>

      <v-col md="2" style="text-align: center"><span style="display: inline-block; padding: 5px 0 0 0px; ">&nbsp;to</span></v-col>

      <v-col md="5">
        <v-text-field
          outlined
          background-color="white"
          dense
          v-model="end"
          placeholder="yyyy"
          :clearable="false"
          :validate-on-blur="true"
          @keydown="typeYear"
        />
      </v-col>
    </v-row>

    <h3 style="margin-top: -8px">Discipline</h3>
      <div v-for="(disc, index) in disciplineList">
        <v-checkbox
          dense
          off-icon="check_box_outline_blank"
          on-icon="check_box"
          class="discipline-checkbox"
          color="primary"
          v-model="disciplinesChecked"
          :label="discLabel(index)"
          :value="disc.value"
          @change="filterSearch"
        />
      </div>
  </div>
</template>
<script>

  import searchConstructor from '@/mixins/searchConstructor'
  import { debounce } from 'lodash';
  import { mapGetters, mapActions } from 'vuex'

  export default {
    name: "Filters",
    mixins: [ searchConstructor ],
    computed: {...mapGetters(['searchResp', 'disciplineList', 'contentTypeList']),
    },

    data: () => ({
      start: null,
      end: null,
      disciplinesChecked: [],
      contentTypesChecked: [],

    }),
    mounted() {
    },
    methods: {...mapActions(['setSearchResp', 'setDisciplines', 'setContentTypes', 'setPubYearStart', 'setPubYearEnd', 'setNewSearchCounts']),

      typeYear: debounce(function() {
        this.checks()
      }, 500),
      checks: function() {
        /* todo ideally wouldn't be reproducing these checks here but i wasn't able to successfully call the ones above -JK */
        if (this.start <= this.end) {
          this.setPubYearStart(this.start)
          this.setPubYearEnd(this.end)
          this.setNewSearchCounts(true)
          this.doSearch(true)
        }
      },
      discLabel(idx) {
        return this.disciplineList[idx].value + ' (' + new Intl.NumberFormat('en-US').format(this.disciplineList[idx].count) + ')'
      },
      contentLabel(index) {
        let label = this.contentTypeList[index].value.replace(/(^\w{1})|(\s{1}\w{1})/g, match => match.toUpperCase())
        if (label == 'Book') {
          label = 'Book Chapter'
        }
        return label + ' (' + new Intl.NumberFormat('en-US').format(this.contentTypeList[index].count) + ')'
      },

      filterSearch() {
        let disciplineString = ''
        if (this.disciplinesChecked) {
          console.log('disciplinesChecked ', this.disciplinesChecked)
          if (this.disciplinesChecked.length === 1) {
            disciplineString = '\"' + this.disciplinesChecked[0] + '\"'
          } else if (this.disciplinesChecked.length > 1) {
            disciplineString = '('
            for (let i = 0; i < this.disciplinesChecked.length; i++) {
              disciplineString += '\"' + this.disciplinesChecked[i] + '\"'
              if (i + 1 < this.disciplinesChecked.length) {
                disciplineString += ' OR '
              } else {
                disciplineString += ')'
              }
            }
          }
        }
        this.setDisciplines(disciplineString)

        let contentTypeString = ''
        if (this.contentTypesChecked) {
          if (this.contentTypesChecked.length === 1) {
            contentTypeString = '\"' + this.contentTypesChecked[0] + '\"'
          } else if(this.contentTypesChecked.length > 1) {
            contentTypeString = '('
            for (let i = 0; i < this.contentTypesChecked.length; i++) {
              contentTypeString += '\"' + this.contentTypesChecked[i] + '\"'
              if (i + 1 < this.contentTypesChecked.length) {
                contentTypeString += ' OR '
              } else {
                contentTypeString += ')'
              }
            }
          }
        }
        this.setContentTypes(contentTypeString)
        this.doSearch(true)
      }
    }
  }
</script>

<style scoped>

  .pub-date {
    margin-bottom: -24px;
  }

  .discipline-checkbox {
    height: 28px;
  }

  .v-input__slot {
    margin-bottom: 8px;
  }

  .pub-date .theme--light.v-label {
    color: #8d8d8d;
  }

  .v-text-field__details {
    display: none;
  }

</style>
