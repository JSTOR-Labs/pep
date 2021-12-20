<template>
  <div>
    <h1>Admin Home</h1><br>
      <!--<h2>Your device configurations</h2>-->

    <v-row>
      <!--<v-col id="completed" class="device-tile" md="3"  @click="skipToConfigure">
        <p class="device-name"> Spring 2020- Cornell, Auburn</p>
        <p>18.4 GB</p>
        <p> Includes 16 disciplines</p>
        <p> Created - 1/20/2019</p>
      </v-col>

      <v-col id="in-progress" class="device-tile" md="3" >
        <p class="device-name"> Fall 2020- Cornell, Auburn</p>
        <p>31.1 GB</p>
        <p>Includes 21 disciplines</p>
        <p>45% export complete (Do not remove H:/)</p>
        <v-progress-linear height="6px" value="45" background-color="#eeeeee" color="accent" class="build-in-progress"/>
      </v-col>-->

      <v-col md="6"><!-- class="create-device" -->
        <p class="new-drive-heading">Select a drive to configure a new device. Its contents will be overwritten.</p>

        <v-row>
          <v-col md="12">
            <v-select
              outlined
              dense
              background-color="white"
              label="Select Drive"
              v-model="selectedDrive"
              :items="driveNames"
                @change="configureDrive"
            ></v-select>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-dialog
      v-model="dialog"
      width="500"
    >
      <v-card>
        <v-card-title
          class="headline grey lighten-2"
          primary-title
        >
          Reformat Drive
        </v-card-title>

        <v-card-text v-if="!reformatting"><br>

          To continue, we need to reformat this drive and ALL data on it will be lost. Do you want to continue?

        </v-card-text>
        <v-card-text v-if="reformatting"><br>
          Reformatting . . .
          <v-progress-linear indeterminate height="6px" value="100" background-color="#eeeeee" color="accent"/>


        </v-card-text>

        <v-divider v-if="!reformatting"></v-divider>

        <v-card-actions v-if="!reformatting">
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            text
            @click="dialog = false"
          >
            No, don't do that
          </v-btn>
          <v-btn
            color="primary"
            text
            @click="reformatDrive"
          >
            Yes, reformat my drive
          </v-btn>
        </v-card-actions>
      </v-card>
       </v-dialog>

  </div>
</template>

<script>
  import { mapGetters, mapActions } from 'vuex'


  export default {
    name: 'AdminHome',
    data: () => ({
      //toggle_disciplines: [],
      drives: [],
      driveNames: [],
      selectedDrive: '',
      driveToFormat: '',
      noDriveMessage: 'No USB detected',
      dialog: false,
      reformatting: false,
    }),
    computed: {...mapGetters(['admin', 'drive', 'partitionName', 'longDriveName'])},
    mounted() {
      this.getDrives()
      this.setDoneCreating(0);
    },
    methods: {...mapActions(['setAdmin', 'setDrive', 'setDoneCreating', 'setPartitionName', 'setPartitionSize', 'setLongDriveName']),
      async getDrives() {

        let resp = await this.$axios.$get("/admin/usb")
        this.drives = resp.data
        console.log('these drives are available ', this.drives)
        this.driveNames = [];
        for (let i = 0; i < this.drives.length; i++) {
          this.driveNames.push(this.fullDriveName(this.drives[i]))
        }
        if (this.driveNames.length === 0) {
          this.driveNames.push(this.noDriveMessage)
        }
      },
      configureDrive() {
        if (this.selectedDrive !== '' && this.selectedDrive !== null && this.selectedDrive !== this.noDriveMessage ) { //we have a real drive
          for (let i = 0; i < this.drives.length; i++) {
            if (this.fullDriveName(this.drives[i]) === this.selectedDrive) {

              this.setDrive(this.drives[i])
            }
          }
         if (typeof this.drive != 'undefined') {//and we found the object representing it
           let thePartition
           this.setLongDriveName(this.drive.name + ' ' + this.drive.model + ' ' + this.drive.size)
           for (let j = 0; j < this.drive['partitions'].length; j++) { //if there are partitions
             console.log(this.drive['partitions'][j])
             if (this.drive['partitions'][j]['filesystem'] === 'exfat') { //and one has the correct formatting . . .
                thePartition = this.drive['partitions'][j]
                this.setPartitionName(thePartition.name);
                console.log('the partition ', thePartition)
               this.setPartitionSize(thePartition.sizeKB / (1024 * 1024));
               this.$router.push({
                 path: '/configure'
               })
                //break;

             }
           }

           if (typeof thePartition == 'undefined') { //ask to reformat
              this.dialog = true
           }
         } else {
           alert ('Error connecting to the drive')
         }
        } else {
          alert('Plug in a USB drive first')
        }
      },
      async reformatDrive() {
        this.reformatting = true;
        console.log('drive to format ', this.drive.name)
        let resp = await this.$axios.$post("/admin/usb", {'drive': this.drive.name})

        console.log('the response, ', resp)
        if (resp.code === 200 ) {
          console.log('this.drive ', this.drive)
          this.setPartitionName(this.drive['partitions'][0].name)
          this.setPartitionSize(this.drive['partitions'][0].sizeKB / (1024 * 1024))
          this.$router.push({
            path: '/configure'
          })
        } else {
          alert('Sorry, reformatting didn\'t work')
        }
        this.dialog = false
        this.reformatting = false

      },
      fullDriveName(drive){
        return drive.name + ' ' + drive.model + ' ' + drive.size
      },

    }
  }
</script>

<style>

  .device-tile p {
    margin-bottom: -2px;
  }

  .device-tile {
    position: relative;
    background: linear-gradient(135deg, #006179, #0885A5);
    height: 130px;
    margin: 8px;
    padding: 8px;
    border-radius: 3px;
    color: #fff;
  }

  .device-tile:hover {
    cursor: pointer;
    box-shadow: 0px 2px 4px 0px rgba(0,0,0,0.5);
  }

  .device-name {
    font-size: 16px;
    font-weight: bold;
  }

  .build-in-progress {
    position: absolute;
    margin-left: -8px;
    margin-top: 8px;
    border-radius: 0 0 3px 3px;
  }

  .create-device {
    background-color:#f5f5f5;
    border: 1px solid #cfcfcf;
    border-radius: 3px;
    height: 130px;
    margin: 8px;
  }

  .select-button {
    height: 40px !important;
  }

  .new-drive-heading {
    margin-bottom: 0px;
    font-weight: bold;
  }

  .v-btn-toggle {
    margin-right: 1em;
    margin-bottom: 1em;
  }

  .v-btn-toggle--dense > .v-btn.v-btn {
    padding: 24px;
  }
</style>

