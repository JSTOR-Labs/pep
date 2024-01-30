<template>
   <div style="background-color: white;">

     <v-row class="request-heading">
       <v-col>
         <h2>Your reading requests ({{ reqs.length }} items)</h2>
         <p class="description">Enter your information to print/download/submit your reading requests.</p>

         <v-row>
           <v-col>
             <v-text-field
               v-model="studentName"
               label="Student Name:"
               class="student-name"
               :rules="nameRules"
             />
           </v-col>
           <!--<v-col> taking this out for now because none of the 5 schools needed it
             <v-text-field label="ID No:" class="id-no"/>
           </v-col>-->
         </v-row>
       </v-col>
     </v-row>

     <!--{{ requests }}<br>
    {{ requests.substring(0,310) }}-->

    <v-card-text v-for="(req, idx) in reqs" :key="idx" id="request-list">
      <div style="line-height: 1.7; font-size: 14px">
        <a class="cancel" @click="removeRequest(req)">Cancel</a>
        <span class="content-type"> {{JSON.parse(req).contentType}} </span> ({{JSON.parse(req).pageCount}} pages)<br>
        <span class="article-title">{{(JSON.parse(req)).title}} </span> <br>
        <span style="line-height: 30px" v-if="JSON.parse(req).authors">Author: {{(JSON.parse(req)).authors[0]}}</span><!-- todo format author string? -->
        <p v-html="(JSON.parse(req)).srcHtml"/>
        <hr color="#DCDCDC" class="hr-results" size="1px">
      </div>
    </v-card-text>

     <v-card-text
     >
         <v-textarea
           class="student-note"
           outlined
           name="student-notes"
           v-model="studentNotes"
           label="Notes (optional)"
         />
     </v-card-text>

     <v-card-text class="button-card">
        <v-row class="ma-0">
          <v-btn
            depressed
            dark
            color="primary"
            @click="print()"
          >
            Print this list
          </v-btn>

          <v-btn
            class="ml-auto"
            depressed
            dark
            color="primary"
            @click="downloadTXT()"
          >
            Download list (.txt)
          </v-btn>

          <v-btn
            class="ml-auto"
            depressed
            dark
            color="primary"
            @click="handleSubmit()"
          >
            Submit Requests
          </v-btn>
        </v-row>
     </v-card-text>

  </div>
</template>

<script>
  import manageCart from '@/mixins/manageCart'
  import { mapGetters, mapActions } from 'vuex'
  import { saveAs } from 'file-saver';

  export default {
    name: "StudentRequest",
    mixins: [ manageCart ],
    data: () => ({
      studentName: '',
      studentNotes: '',
      loaded: false,
      nameRules: [
        v => !!v || 'Name is required',
        v => (v && v.length <= 30) || 'Name must be less than 30 characters'
      ],
    }),

    mounted() {
      window.addEventListener('afterprint', this.afterPrint);
      this.loaded = true
    },
    destroyed() {
      window.removeEventListener('afterprint', this.afterPrint);

    },
    computed: {...mapGetters(['reqs']),

    },
    methods: {...mapActions(['setShowCart']),

      async handleSubmit() {
        if (this.studentName !== "") {
          if (this.reqs.length>0) {
            let idArray = []
            for (let i = 0; i < this.reqs.length; i++) {
              idArray.push(JSON.parse(this.reqs[i])['_id'])
            }
            let args = {
              name: this.studentName,
              notes: this.studentNotes,
              articles: idArray
            }
            console.log('request args in StudentRequest to /api/request: ', args)
            let resp = await this.$api.basic.request(args)
            console.log('request status: ', resp)
            this.clearRequests()
            this.setShowCart(false);
            this.$emit('clicked', 'x')
          }
          else {
            alert('No articles to request')
          }
        }
        else {
          alert('Enter student name to submit')
        }
      },

      close() {
        //this.emitGaEvent('request', 'cancel')
      },

      print() { /*had to do this instead of print stylesheet. ask JK for details if needed */

        document.getElementsByClassName('v-overlay')[0].style.visibility = 'hidden'
        document.getElementsByClassName('v-application--wrap')[0].style.visibility= 'hidden'
        document.getElementsByClassName('v-application--wrap')[0].style.height = '100px'
        document.getElementsByClassName('button-card')[0].style.visibility = 'hidden'
        document.getElementsByClassName('cancel')[0].style.visibility = 'hidden'
        document.getElementsByClassName('description')[0].style.visibility = 'hidden'
        document.getElementsByClassName('v-dialog')[0].style.margin = '0px'
        document.getElementsByClassName('v-dialog')[0].style.boxShadow = 'none'
        document.getElementsByClassName('v-dialog__content')[0].style.alignItems = 'flex-start'
        document.getElementsByClassName('student-name')[0].style.marginBottom = '-30px'
        document.getElementsByClassName('hr-results')[0].style.marginTop = '8px'
        document.getElementsByClassName('v-dialog')[0].style.maxHeight = '100%'
        window.print();

      },
      afterPrint()  {

        if (this.loaded) {

          document.getElementsByClassName('v-overlay')[0].style.visibility = 'visible'
          document.getElementsByClassName('v-application--wrap')[0].style.visibility = 'visible'
          document.getElementsByClassName('v-application--wrap')[0].style.height = 'unset'
          document.getElementsByClassName('button-card')[0].style.visibility = 'visible'
          document.getElementsByClassName('cancel')[0].style.visibility = 'visible'
          document.getElementsByClassName('description')[0].style.visibility = 'visible'
          document.getElementsByClassName('v-dialog')[0].style.margin = '24px'
          document.getElementsByClassName('v-dialog')[0].style.boxShadow = '0px 11px 15px -7px rgba(0, 0, 0, 0.2), 0px 24px 38px 3px rgba(0, 0, 0, 0.14), 0px 9px 46px 8px rgba(0, 0, 0, 0.12);'
          document.getElementsByClassName('v-dialog__content')[0].style.alignItems = 'center'
          document.getElementsByClassName('student-name')[0].style.marginBottom = '-16px'
          document.getElementsByClassName('hr-results')[0].style.marginTop = '24px'
        }
      },
      downloadTXT() {
        console.log(this.reqs)
        if (this.reqs.length > 0) {
          let save = ''

          for (let i = 0; i < this.reqs.length; i++) {
            let request = JSON.parse(this.reqs[i])
            console.log(request)
            save = save + request.title + '\n'
            save = save + request.srcHtml + '\n'
            save = save + 'https://www.jstor.org/stable/' + request['_id'] + '\n'
            save = save + '\n'
          }

          let blob = new Blob([save], {
            type: "text/plain;charset=utf-8"
          });
          let date = new Date()
          let name = (date.getMonth() + 1) + '-' + date.getDate() + '-' + date.getFullYear() + ' requests.txt'
          if (this.studentName !== '') {
            name = this.studentName + '\'s ' + name
          }
          saveAs(blob, name);

        } else {
          alert('No requests to save')
        }
      }
    }
  }
</script>

<style scoped>
  * {
    font-family: "Arial", sans-serif;
  }

  .request-heading {
    background-color: #F2F2F2;
    height: 100%;
    padding: 0px 16px 0px 16px;
  }

  .description {
    margin-bottom: 0;
  }

  .student-name {
    margin-bottom: -16px;
    padding: 0;
  }

  .id-no {
    margin-bottom: -16px;
    padding: 0;
  }

  .v-card__text {
    height:auto;
  }

  .cancel {
    float: right;
    color: #BE0101;
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
    line-height: 0;
  }

  .hr-results {
    margin-top: 24px;
  }

  .request-page-buttons {
    padding: 16px;
  }



</style>

<style>

  .student-note .v-text-field__details {
  display: none;
  }



</style>
