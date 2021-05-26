<template>
  <div>
    <div v-if="admin" id="configuration">

      <h1>Configure Drive</h1>
      <p>Select discipline indexes you want to include on your drive. You can also choose to include available offline content.
        If offline content is selected, the corresponding discipline must also be selected.</p>


      <v-row>
        <v-col md="8" class="pa-0">
          <!--{{ selectedOptions }}-->
          <v-simple-table class="disciplines-table">
            <thead>
            <tr>
              <th style="width:55%"> Discipline name</th>
              <th style="width:20%"> Search index</th>
              <th style="width:25%"> + Offline content</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="(disc, idx) in contents">
              <td class="disc-name">{{ disc.name }} ({{ discTotal(disc) }} GB total) </td>

              <td class="disc-search-index">
                <v-checkbox
                  :value="disc.name + '|index|' + disc.indexSize"
                  :ref="disc.indexSize"
                  color="primary"
                  v-model="selectedOptions"


                @change="calculateSizes(`${disc.value}index`)"
                :label="`${disc.indexSize.toLocaleString()} GB`">


                </v-checkbox>
              </td>
              <td class="disc-offline-content">
                <v-checkbox
                  :value="disc.name + '|content|' + disc.contentSize "
                  color="primary"
                  @change="calculateSizes(disc.contentSize)"
                  :label="`${disc.contentSize.toLocaleString()} GB`"
                  v-model="selectedOptions"
                  :disabled="disc.contentSize == 0"
                  >

                </v-checkbox>
              </td>
            </tr>

            </tbody>
          </v-simple-table>
        </v-col>

        <v-col md="4" class="pa-0">
          <div></div>
        </v-col>
      </v-row>

      <div class="container" style=" width: 400px; height: 500px">
        <div class="floater">
          <h3 style="margin-top: 0px; margin-bottom: 8px">Select disciplines to begin</h3>

          <p style="margin-top: -16px">Disk space remaining in {{this.longDriveName}}</p>
          <h1 style="margin-top: 0px" :class="{ tooBigRed: tooBig }">{{ (partitionSizeNum - totalSize).toLocaleString() }} GB</h1>
          <!-- "GB" is hardcoded here, and if a drive is actually XMB, this will misrepresent it.
          But god help them for trying with something so small - shouldn't be a realistic use case-->
          <p>Your selected contents ({{totalSize.toLocaleString() }} GB)</p>
          <v-progress-linear height="8px" :value="percentUsed" rounded background-color="white" :color="progressColor" class="file-size-bar"/>
          <v-simple-table class="content-table" dense>
            <tbody>
            <tr v-for="option in optionsForTable">
              <td style="width: 70%" v-html="option.title"></td>
              <td style="width: 30%">{{ option.size }}</td>
            </tr>
            </tbody>
          </v-simple-table>

          <nuxt-link to="building" v-if="!tooBig">
            <v-btn
              class="export-button"
              depressed
              dark
              color="primary"
              dense
              @click="saveToUSB"
            >Export</v-btn>
          </nuxt-link>
          <v-btn
          class="export-button"
          depressed
          color="primary"
          dense
          disabled
          v-if="tooBig"
          >Export</v-btn>
        </div>
      </div>
    </div>


  </div>
</template>

<script>

  import { mapGetters, mapActions } from 'vuex'
  import searchConstructor from '@/mixins/searchConstructor'


  export default {

    name: 'configure',
    mixins: [ searchConstructor ],
    middleware: 'authenticated',
    computed: {...mapGetters(['admin', 'drive', 'partitionName', 'partitionSize', 'longDriveName']),
      percentUsed() {
        console.log('this.totalSize ', this. totalSize)

        return this.totalSize/this.partitionSizeNum * 100
      },
      tooBig() {
        return this.percentUsed > 100
      },
      progressColor() {
        if (this.tooBig) {
          return 'red'
        } else {
          return '#006179'
        }
      }
    },
    data: () => ({
      systemDataSize: .05, //in GB, approx system files max size
      selectedOptions: [],
      contents: [],
      snapshot: '',
      totalSize: 0,
      optionsForTable: [],
      optionsForApi: {},
      partitionSizeNum: 0,
      timer: null

    }),
    mounted() {
      this.getContents()
      this.partitionSizeNum = this.partitionSize.substring(0, this.partitionSize.length - 3 )
      this.totalSize = this.systemDataSize
      this.resetLists()
    },
    methods:  {...mapActions(['setDriveName', 'setConfiguration', 'setDoneCreating']),

      calculateSizes() { //and format list for multiple outputs (the table, the API . . )

        this.totalSize = this.systemDataSize
        this.checkIndexes()
        this.resetLists()
        this.setOptionsForApi()
        this.setOptionsForList();

        //discipline|type|size

      },
     async checkProgress() {
        console.log('checking progress')
       try {
         let resp = await this.$axios.$get('/admin/snapshot?snapshot=' + this.snapshot);
         console.log('is it ready? ', resp)
         if (resp.data.failed === true) {
           clearInterval(this.timer);
           this.setDoneCreating(-1)
         } else if (resp.data.done === true) {
           clearInterval(this.timer);
           this.setDoneCreating(1);
         }
       }catch (e) {
          console.log('the error ', e)
       }

      },
      async getContents() {

        let resp = await this.$axios.$get("/admin/indices")
        this.contents = resp.data
        console.log('contents ', this.contents)

      },
      async saveToUSB() {
        console.log(this.partitionName)
        this.setConfiguration(this.getDescription())
        try {
          console.log("/admin/usb/" + this.partitionName, {'indices': this.optionsForApi})
          let resp = await this.$axios.$post("/admin/usb/" + this.partitionName, {'indices': this.optionsForApi})
          console.log('USB save response', resp)
          this.snapshot = resp.data;
          this.timer = setInterval(this.checkProgress, 10000);
        } catch (err) {
            console.log(err.response.data)
          }
      },
      getDescription() {
        let config = ''
        let title = ''
          for (let i = 0; i < this.optionsForApi.length; i++) {
            title = this.optionsForApi[i]
            if (i === 0) { //first
              config = title
            } else if (i + 1 < this.optionsForApi.length) { //not first or last last
              config = config + ', ' + title
            } else if (i+ 1 === this.optionsForApi.length) { //last
              config = config + ' and ' + title
            }
          }
        return config
      },
      discTotal(disc) {
        return (disc.indexSize + disc.contentSize).toLocaleString()
      },

      checkIndexes() {
        let newCheckboxes = []
        for (let i = 0; i < this.selectedOptions.length; i++) {
          if (this.selectedOptions[i].indexOf('content') > -1) { //if it's a content checkbox
            let lookingFor = this.selectedOptions[i].split("|")[0] + '|index';
            let indexIsChecked = false
            for (let j = 0; j < this.selectedOptions.length; j++) { //look through everything that's checked
              if (this.selectedOptions[j].indexOf(lookingFor) > -1) { //to see if it's index is checked too
                indexIsChecked = true;
                break;
              }
            }
            if (!indexIsChecked) { //if it's index ISN'T checked . . .
              lookingFor = this.selectedOptions[i].split("|")[0]
              for (let k = 0; k < this.contents.length; k++) { //look through the full list
                if (this.contents[k].name === lookingFor) {
                  newCheckboxes.push(this.contents[k].name + '|index|' +this.contents[k].indexSize) //and check it's index
                }
              }
            }
          }
        }
        for (let l = 0; l < newCheckboxes.length; l++) {
          this.selectedOptions.push(newCheckboxes[l])
        }
      },
      resetLists() {
        this.optionsForApi = {}
        this.optionsForTable = []
        let option = []
        option.title = "<b>System Files</b>"
        option.label = ""
        option.size = this.systemDataSize.toLocaleString() + ' GB'
        this.optionsForTable.push(option)

      },
      isIn(option, list) {
        return list.hasOwnProperty(option);
      },

      setOptionsForApi() {

        //when API includes content
        /*
        {"indices": {"biology": {"includeContent": true}}}
        */
        for (let m = 0; m < this.selectedOptions.length; m++) {
          this.totalSize = this.totalSize + parseInt(this.selectedOptions[m].split("|")[2])
          if (this.selectedOptions[m].indexOf('index') > -1) {//if it's an index checkbox
            if (!this.isIn(this.selectedOptions[m].split("|")[0], this.optionsForApi)) {//and not already added
              this.optionsForApi[this.selectedOptions[m].split("|")[0]] = {'includeContent': false}
            }
          } else { //it's a content checkbox
            this.optionsForApi[this.selectedOptions[m].split("|")[0]] = {'includeContent': true}
          }
        }
        console.log('full options for API ', this.optionsForApi)
      },
      setOptionsForList() {
        for (let i = 0; i < this.selectedOptions.length; i++) {
          if (this.selectedOptions[i].indexOf('index') > -1) {//if it's an index checkbox
            let correspondingcontent = this.selectedOptions[i].split("|")[0] + "|content"
            let pair = null
            for (let j = 0; j < this.selectedOptions.length; j++) {
              if (this.selectedOptions[j].indexOf(correspondingcontent) > -1) {
                pair = this.selectedOptions[j]
                break
              }
            }
            let option = []

            if (pair) {
              option.title = '<b>' + this.selectedOptions[i].split("|")[0] + '</b> <i>(index+content)</i>'
              option.size = (parseFloat(this.selectedOptions[i].split("|")[2]) + parseFloat(pair.split("|")[2])).toLocaleString() + ' GB'

            } else {

              option.title = '<b>' + this.selectedOptions[i].split("|")[0] + '</b> <i>(index)</i>'
              option.size = parseFloat(this.selectedOptions[i].split("|")[2]).toLocaleString() + ' GB'

            }
            this.optionsForTable.push(option)

          }
        }
      }

    }
  }
</script>

<style scoped>

  .floater p {
    margin: 0px;
  }

  .disciplines-table {
    background-color: #F5F5F5;
  }

  thead th {
    color: black !important;
    font-size: 18px;
    height: 72px;
    text-align: left;
  }

  td {
    height: 48px;
  }

  .floater {
    padding: 8px;
    position: absolute;
    left: calc(50vw + 230px);
    top: 217px;
    width: 315px;
    background-color: #f5f5f5;
    border-radius: 6px;
    box-shadow: 0px 2px 4px 0px rgba(0,0,0,0.5);
    z-index: 20
  }

  .file-size-bar {
    border: 1px solid #DADADA;
    margin-top: 0px;
  }

  .v-text-field__details {
    display: none;
  }

  .content-table {
    background-color: #f5f5f5;
    margin-top: 8px;
  }

  .content-table td {
    height: auto;
    font-size: 12px;
    padding: 0px;
    line-height: 1.5;
  }

  .content-table tr td {
    border-bottom: none !important;
  }

  .export-button {
    display: block;
    width: 130px;
    align-items: center;
    margin: 32px auto
  }

  .disciplines-table .disc-search-index .v-input , .disciplines-table .disc-offline-content .v-input  {
    margin-top: 8px;
  }

  .tooBigRed {
    color: red
  }

  .tooBigProgress .v-progress-linear__determinate.primary {
    background-color: red!important;
    border-color: red!important
  }

  #configuration p {
    margin-bottom: 1em;
  }
</style>

<style>
  #configuration .v-messages {
    display:none;
  }
</style>
