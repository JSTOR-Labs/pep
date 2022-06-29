<template>
  <div v-if="admin" id="request-list">

    <v-row>
      <v-col cols="5">
        <h1>View Student Requests</h1>
      </v-col>

      <v-col cols="7" class="request-actions" id="hide-request">
        <v-btn
          outlined
          height="40px"
          color="primary"
          @click="uploadFile()"
        >Import list (.csv)
        </v-btn>

        <v-file-input accept=".csv" id="fileUpload" v-model="chosenFile" @change="importCSV"/>

        <json-c-s-v
          :data="requestsForExport">
          <v-btn
            outlined
            height="40px"
            color="primary"
          >Export list (.csv)
          </v-btn>
        </json-c-s-v>

        <v-btn
          outlined
          height="40px"
          color="primary"
          @click="print()"
        >Print this list
        </v-btn>

        <v-select
          outlined
          dense
          id="selectRequests"
          color="primary"
          value="Pending"
          :items="items"
          @change="getRequests(true)"
        />
      </v-col>
    </v-row>

    <!--debug: sort={{currentSort}}, dir={{currentSortDir}}-->
    <v-simple-table class="disciplines-table">
      <thead>
      <tr>
        <th style="width:30%" @click="sort('title')" class="th-sort">
          Title
          <span class="arrows">
            <span class="up"
                  :class="{active: this.currentSort === 'title' && this.currentSortDir === 'asc'}">&#9650;</span>
            <span class="down"
                  :class="{active: this.currentSort === 'title' && this.currentSortDir === 'desc'}">&#9660;</span>
          </span>

        </th>
        <th style="width:14%" @click="sort('studentName')" class="th-sort">
          Name
          <span class="arrows">
            <span class="up" :class="{active: this.currentSort === 'studentName' && this.currentSortDir === 'asc'}">&#9650;</span>
            <span class="down" :class="{active: this.currentSort === 'studentName' && this.currentSortDir === 'desc'}">&#9660;</span>
          </span>
        </th>
        <th style="width:12%" @click="sort('dateRequested')" class="th-sort">
          Requested
          <span class="arrows">
            <span class="up" :class="{active: this.currentSort === 'dateRequested' && this.currentSortDir === 'asc'}">&#9650;</span>
            <span class="down"
                  :class="{active: this.currentSort === 'dateRequested' && this.currentSortDir === 'desc'}">&#9660;</span>
          </span>
        </th>
        <th style="width:10%" @click="sort('numPages')" class="th-sort">
          Pages
          <span class="arrows">
            <span class="up"
                  :class="{active: this.currentSort === 'numPages' && this.currentSortDir === 'asc'}">&#9650;</span>
            <span class="down" :class="{active: this.currentSort === 'numPages' && this.currentSortDir === 'desc'}">&#9660;</span>
          </span>
        </th>
        <th style="width:15%">
          Notes
        </th>
        <th class="internet-yes hide-action-col" v-if="statusIsPending" style="width:19%"> Action</th>
        <th v-else style="width:19%">Status</th>

      </tr>
      </thead>
      <tbody>
      <tr 
        v-for="(request, i) in sortedRequests"
        :key="`request_${i}`">

        <v-tooltip bottom content-class="tooltip">
          <template v-slot:activator="{ on }" >
              <td v-on="on">
                <a v-if="request.pdf && !online">{{ request.title }} </a> <!-- pdf available, not online -->
                <span v-if="!request.pdf && !online">{{ request.title }}</span> <!-- pdf not available, not online -->
                <a v-if="online" :href="JSTORurl(request.id)" target="_blank">{{ request.title }}</a><!-- online -->

              </td>
          </template>
          <span> {{ formatCitation(request.srcHtml) }} </span>
        </v-tooltip>

        <td> {{ request.studentName }} </td>
        <!-- .toLocaleDateString() + ' ' + d.toTimeString().substring(0, d.toTimeString().indexOf("GMT")) -->
        <td> {{ new Date(request.dateRequested).toLocaleDateString() }} |
          {{ new Date(request.dateRequested).toTimeString().substring(0, new
          Date(request.dateRequested).toTimeString().indexOf("GMT")).substring(0,5) }}
        </td>
        <td> {{ request.numPages }}</td>
        <td> {{ request.notes }} </td>
        <td class="internet-yes hide-action-col" v-if="statusIsPending">
          <span class="request-buttons" >
             <v-btn
               v-if="online || request.pdf"
               depressed
               dark
               x-small
               color="accent"
               @click="update(0, request)"
             >Print
            </v-btn>
            <v-btn
              v-if="!online && !request.pdf"
              disabled
              depressed
              x-small
              title="This article isn't available until you get online"
            >Print
            </v-btn>

            <v-btn
              v-if="online || request.pdf"
              depressed
              dark
              x-small
              color="accent"
              @click="update(1, request)"
            >PDF
            </v-btn>
             <v-btn
               v-if="!online && !request.pdf"
               disabled
               depressed
               x-small
               title="This article isn't available until you get online"
             >PDF
            </v-btn>

            <v-btn
              depressed
              dark
              x-small
              color="deny"
              @click="update(3, request)"
            >Deny
            </v-btn>
          </span>
        </td>

        <td v-else>
          {{ formatStatus(request.status) }}
        </td>

        <!--<td>{{ request }}</td>-->
      </tr>
      <tr v-if="!requests">
        <br>
        No requests at this time
        <br><br><br>
      </tr>

      </tbody>
    </v-simple-table>

  </div>

</template>

<script>

  import { mapGetters } from 'vuex'

  import JsonCSV from '~/components/JsonCSV'
  import PapaParse from 'papaparse'

  export default {

    name: 'requests',
    middleware: 'authenticated',
    computed: {
      ...mapGetters(['admin']),
      sortedRequests: function () {
        if (this.requests) {
          return this.requests.sort((a, b) => {
            let modifier = 1;
            if (this.currentSortDir === 'asc') modifier = -1;
            if (a[this.currentSort] < b[this.currentSort]) return -1 * modifier;
            if (a[this.currentSort] > b[this.currentSort]) return 1 * modifier;
            return 0;
          })
        } else {
          return null
        }
      },
      online: function() {

        if (navigator) {
          return navigator.onLine
        } else {
          return false;
        }
      
      },

    },
    components: {

      JsonCSV
    },
    data: () => ({
      requests: [],
      requestsForExport: [],
      items: ['Pending', 'Completed'],
      statusIsPending: true,
      currentSort: 'dateRequested',
      currentSortDir: 'asc',
      chosenFile: null,
      jstor: 'https://www.jstor.org/stable/',
    }),
    mounted() {
      //this.component('downloadCsv', JsonCSV)
      this.getRequests(false);
    },

    methods: {
      async getRequests(switchView) {
        if (switchView) {
          this.statusIsPending = !this.statusIsPending
        }
        let params = {'pending': this.statusIsPending}
        let resp = await this.$axios.$get("/admin/request", {params: params})
        console.log('studentRequests in requests: ', resp)
        this.requests = resp.data.requests
        this.requestsForExport = _.cloneDeep(resp.data.requests)
        //let dois = []
        for (let r = 0; r < this.requestsForExport.length; r++) { //add URLS in case CSV is exported
          delete this.requestsForExport[r]['pdf']
          if (this.requestsForExport[r]['status'] === 0) {
            this.requestsForExport[r]['status'] = 'printed'
          }
          if (this.requestsForExport[r]['status'] === 1) {
            this.requestsForExport[r]['status'] = 'saved as PDF'
          }
          if (this.requestsForExport[r]['status'] === 2) {
            this.requestsForExport[r]['status'] = 'pending review'
          }
          if (this.requestsForExport[r]['status'] === 3) {
            this.requestsForExport[r]['status'] = 'denied'
          }
          this.requestsForExport[r]['url'] = '=HYPERLINK("' + this.jstor + this.requestsForExport[r].id + '")'

          //the code below replaces the JSTOR URL with Click here
          // this.requestsForExport[r]['url'] = '=HYPERLINK("' + this.jstor + this.requestsForExport[r].id +'"' + ',' + '"' + 'Click here' + '")'
        }

      },
      async update(type, article) { // 0 = Print, 1 = PDF, 2 = Pending, 3 = denied
        let vars = {
          requestID: article.requestID,
          articleID: article.id,
          status: type
        }
        if (type === 0) {//print request
          if (this.online) { //todo or pdf is available
            open("https://www.jstor.org/stable/pdf/" + article.id + ".pdf")
          } else {
            let targetPDF = await this.$axios.get("/admin/pdf/" + article.id, { responseType: "blob" })
            this.showPDF(targetPDF.data, article.id);
          }

        } else if (type === 1) { //pdf request
          if (this.online) { //todo or pdf is available
            open("https://www.jstor.org/stable/pdf/" + article.id + ".pdf")
          } else {

            let targetPDF = await this.$axios.get("/admin/pdf/" + article.id, { responseType: "blob" });
            this.savePDF(targetPDF.data, article.id);
          }
        }


        try {
          let resp = await this.$axios.$patch("/admin/request", vars)
        } catch (err) {
          console.log('logging the error: ', err.response.data);
        }
        this.getRequests(false);

      },
      savePDF(blob, doi) {
        const a = document.createElement("a");
        a.style.display = "none";
        //let decoded = window.atob(blob);
        // Set the HREF to a Blob representation of the data to be downloaded
        a.href = window.URL.createObjectURL(blob);

        // Use download attribute to set set desired file name
        a.setAttribute("download", doi + ".pdf");
        document.body.appendChild(a);
        // Trigger the download by simulating click
        a.click();

        // Cleanup
        window.URL.revokeObjectURL(a.href);
        document.body.removeChild(a);
      },
      showPDF(blob, doi) {
        let fileURL = window.URL.createObjectURL(blob);

        window.open(fileURL).print();
      },
      formatStatus(status) {
        if (status === 0) {
          return 'Printed'
        } else if (status === 1) {
          return 'Downloaded'
        } else if (status === 3) {
          return 'Denied'
        } else {
          return 'Unknown'
        }
      },
      formatCitation(cite) {

        if (cite.indexOf('</cite>') > -1) {
          cite = cite.substring(6);
          let first = cite.substring(0, cite.indexOf('</cite>'))
          let second = cite.substring(cite.indexOf('</cite>') + 7)
          return "from " + first + second
        } else {
          return cite
        }

      },
      print() {
        window.print();
      },
      sort(s) {
        //if s == current sort, reverse
        if (s === this.currentSort) {
          this.currentSortDir = this.currentSortDir === 'asc' ? 'desc' : 'asc';
        }
        this.currentSort = s;
      },
      uploadFile() {
        document.getElementById('fileUpload').click()
      },
      importCSV() {
        let json = []
        if (this.chosenFile) {
          let reader = new FileReader();
          reader.readAsText(this.chosenFile);
          reader.onload = () => {
            json = PapaParse.parse(reader.result, {header: true, dynamicTyping: true})
            console.log('json.data ', json.data)

            let imported = [];

            for (let item in json.data) {
              imported.push(json.data[item]);
            }

            for (let i = 0; i < this.requests.length; i++) {
              let found = false;
              for (let j = 0; j < imported.length; j++) {

                if (imported[j].id == this.requests[i].id &&
                  imported[j].studentName == this.requests[i].studentName &&
                  imported[j].dateRequested === this.requests[i].dateRequested) {
                  found = true;
                  //console.log('found true for ', imported[j].id)
                  //imported[j].nodes = imported[j].nodes.concat(this.requests[i].nodes);
                  break;
                }
              }
              if (!found) {
                imported.push(this.requests[i]);
              }
            }
             //console.log(JSON.stringify(imported))
            this.requests = imported
          }

        }

      },
      JSTORurl: function(id) {
        return 'https://www.jstor.org/stable/' + id
      }
    }
  }

</script>

<style scoped>

  .disciplines-table {
    margin-bottom: 16px;
  }

  .disciplines-table th {
    background-color: #F5F5F5;
    color: black !important;
    font-size: 16px;
    height: 56px;
    text-align: left;
    padding: 6px;
  }

  td {
    padding: 8px;
    line-height: 1.6;
  }

  td:not(:last-child), th:not(:last-child) {
    border-right: 1px solid #dcdcdc;
  }

  .table-small-font {
    font-size: 12px;
    line-height: 1.2;
  }

  .v-tex-field__details {
    display: none !important;
  }

  .internet-no .v-btn {
    color: #727272;
    background-color: #dddddd !important;
  }

  .request-actions {
    display: inline-block;
  }

  .request-actions .v-select {
    display: inline-block;
    width: 142px;
  }

  .arrows {
    display: inline-block;
    float: right;
    visibility: hidden;
    color: #cccccc;
  }

  .th-sort:hover {
    cursor: pointer;
  }

  th:hover > .arrows {
    cursor: pointer;
    visibility: visible;
  }

  .up {
    position: relative;
    bottom: 6px;
    left: 20px;
  }

  .down {
    position: relative;
    bottom: -10px;
  }

  .active {
    visibility: visible;
    color: #727272;
  }

  .request-buttons .v-btn {
    font-size: 12px !important;
  }

</style>

<style>
  .request-actions .v-input__control .v-input__slot fieldset {
    border-color: #006179;
  }

  .request-actions .v-select__slot label, .request-actions .v-select__slot .v-icon {
    color: #006179;
  }

  #request-list .v-file-input {
    display: none;
  }

  @media print {
    #header {
      display: none;
    }

    #footer {
      display: none;
    }

    #hide-request {
      display: none;
    }

    .arrows {
      display: none !important;
    }

    .hide-action-col {
      display: none !important;
    }
  }

  .tooltip {
    opacity: 1 !important;
    background: #00152b !important;
    max-width: 400px;
    line-height: 1.5 !important;
    cursor: pointer;
  }

</style>
