<template>
  <div class="container">
    <v-row >
      <v-col md="3" id="filters-col">
        <filters></filters>
      </v-col>

      <v-col md="9" id="results-col">
        <results></results>
      </v-col>
    </v-row>

  <!--  <div v-for="(facetData, facet) in searchResp.facets || []" :key="facet">
      {{facet}}
      <ul>
        <li v-for="fval in facetData || []" :key="fval.value">
          {{fval.value}} {{new Intl.NumberFormat('en-US').format(fval.count)}}
        </li>
      </ul>
    </div>-->

  </div>
</template>

<script>

  import Filters from '@/components/Filters'
  import Results from '@/components/Results'

  export default {
    data: () => ({
      start: 1700,
      end: 2019,

    }),
    components: { Filters, Results },

  }
</script>

<style>

  #filters-col {
    background-color: #F8F8F8;
    padding: 16px;
    height: 100%;
  }

  #results-col {
    padding-left: 36px;
  }

</style>
