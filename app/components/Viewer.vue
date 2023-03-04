<template>
      <v-row class="mt-5">
         <h1>
         {{pdfID}}
         </h1>
         <v-col>
            <v-row v-if="pdfRetrieved === null">
               <v-col>
                  <v-progress-circular
                     
                     indeterminate
                     color="primary"
                  />
               </v-col>
            </v-row>
            <div v-else-if="pdfRetrieved">
               <v-row>
                  <v-col cols="auto">
                     <v-select
                        :items="Array.from({length: numPages}, (v, i) => i + 1)"
                        v-model="page"
                        class="sort mb-3"
                        dense
                        hide-details
                        label="Page"
                        aria-label="Select page"
                     >
                        <template v-slot:selection="{ item }">
                        <span class="font-weight-bold">Page:&nbsp;</span><span class="font-weight-bold text-capitalize">{{ item }}</span>
                        </template>
                     </v-select>
                  </v-col>
                  <v-spacer></v-spacer>
                  <v-col cols="auto">
                     <v-btn
                        class="pdf-button"
                        depressed
                        dark
                        color="primary"
                        @click="downloadPDF()"
                     >
                        Download
                     </v-btn>
                  </v-col>
                  <v-col cols="auto">
                     <v-btn
                        class="pdf-button"
                        depressed
                        dark
                        color="primary"
                        @click="printPDF()"
                     >
                        Print
                     </v-btn>
                  </v-col>
               </v-row>
               <div
                  class="pdfContainer" 
                  scroll="changePage"
               >
                  <pdf 
                     class="content pa-3" 
                     v-for="i in numPages" 
                     :src="pdfData" 
                     :key="`page_${i}`" 
                     :id="`page_${i}`" 
                     :page="i" 
                     :scale="scale"
                     :resize="true"
                  />
               </div>
            </div>

            <v-row v-else>
               <v-col  class="text-center">
                  <p class="text-h4 primary--text">
                     Access Denied
                  </p>
                  <p class="text-subtitle-1">
                     This document may not be available or approved in your area.
                  </p>
               </v-col>
            </v-row>

         </v-col>

      </v-row>
</template>
<script>
import pdfvuer from 'pdfvuer'
import printJS from 'print-js'

export default {
  components: {
    pdf: pdfvuer,
  },
  props: {
   id: {
      type: String,
      required: true,
   },
   doi: {
      type: String,
      required: true,
   }
  },
  computed: {
   pdfId () {
      return this.doi + "/" + this.id
   },
  },
  data () {
    return {
      page: 1,
      pageNum: 1,
      max: 2,
      numPages: 0,
      errors: [],
      scale: 'page-width',
      pageHeight: 'page-height',
      data: undefined,
      pdfdata: undefined,
      pdfRetrieved: null,
      blob: undefined,
    }
  },
  mounted () {
   console.log("Component mounted")
    this.getPdf()
  },
  watch: {
    show: function (s) {
      if(s) {
        this.getPdf();
      }
    },
    key: function (a) {
       this.page = a;
    },
    
    page: function (p) {
      const currentPage = document.getElementById(`page_${p}`)
      const nextPage = document.getElementById(`page_${p+1}`)
      if( window.pageYOffset <= this.findPos(currentPage) || 
      ( nextPage && window.pageYOffset >= this.findPos(nextPage) )) {
         currentPage.parentNode.scrollTo({top: currentPage.offsetTop, behavior: 'smooth'})
      }
    }
    
  },
  methods: {
   downloadPDF () {
      const a = document.createElement("a");
      console.log("downloading", this.pdfId)
      a.style.display = "none";
      document.body.appendChild(a);
      a.href = window.URL.createObjectURL(this.blob);
      // Use download attribute to set set desired file name
      a.setAttribute("download", this.pdfId + ".pdf");
      // Trigger the download by simulating click
      a.click();
      // Cleanup
      window.URL.revokeObjectURL(a.href);
      document.body.removeChild(a);
   },
   printPDF () {
      const href = window.URL.createObjectURL(this.blob)
      printJS(href)
   },
   findPos(obj) {
      return obj.offsetTop;
   },
   changePage () {
      var i = 1, count = Number(pdf.numPages);
      const currentPage = document.getElementById(`page_${i}`)
      const nextPage = document.getElementById(`page_${i+1}`)      
      do {
      if(window.pageYOffset >= self.findPos(currentPage) && 
            window.pageYOffset <= self.findPos(nextPage)) {
         self.page = i
      }
      i++
      } while ( i < count)
      if (window.pageYOffset >= self.findPos(currentPage)) {
      self.page = i
      }
   },
   createBlob (byteString) {
      const arrayBuffer = new ArrayBuffer(byteString.length);
      const int8Array = new Uint8Array(arrayBuffer);
      for (let i = 0; i < byteString.length; i++) {
        int8Array[i] = byteString.charCodeAt(i);
      }
      const blob = new Blob([int8Array], { type: 'application/pdf'});
      return blob;
   },
   blobToBase64(blob) {
      return new Promise((resolve, _) => {
         const reader = new FileReader();
         reader.onloadend = () => resolve(reader.result);
         reader.readAsDataURL(blob);
      });
   },
   async getPdf() {
      let resp
      let finalString
      try {

         resp = await this.$api.basic.pdf.get(this.pdfId)
         this.blob = resp.data

         const base64 = await this.blobToBase64(this.blob)
         const splitString = base64.split(',')[1]
         finalString = window.atob(splitString)

         this.pdfRetrieved = true
      } catch (err) {
         console.log(err)
         this.pdfRetrieved = false
         return
      }
      if (this.pdfRetrieved) {
         var self = this;
         self.pdfData = pdfvuer.createLoadingTask({ data: finalString });
         self.pdfData.then(pdf => {
            self.pdfdata = pdf;
            self.numPages = pdf.numPages;
         })
      }
   },
   prev() {
      if (this.page > 1){
         this.page = this.page-1;
      }
      
   },
   next() {
      if (this.page < this.numPages){
         this.page = this.page + 1;
      }
   },

  }
}
</script>


<style lang="css" scoped>
  .pdfContainer {
    border-style: solid;
    border-width: 2px;
    border-color: black;
    height: 100vh;
    overflow-y: scroll;
    position: relative;
  }
</style>