# JSTOR Access in Prison Initiative

## Chromebook Installation Guide

1. In the settings menu, select Advanced>Developers.
    ![Screenshot with arrows pointing at the "Advanced" and "Developers" options.](images/step1.png)
2. In the card for Linux Development Environment, select Turn On.
    ![Screenshot with an arrow pointing at the "Turn On" button](images/step2.png)
3. When prompted to set up the Linux development environment, select Next.
    ![Screenshot with an arrow pointing at the "Next" button.](images/step3.png)

4. Choose a username and set a disk size of 35GB, then select Install. Make a note of the username you selected. 
    ![Screenshot showing the options for a Linux development environment. Arrows point at the username field, the disk size slider, and then "Install" button.](images/step4a.png)
    The installation may take several minutes.
    ![Screenshot showing the progress monitor with a partially filled bar waiting for the Linux environment to finish installing.](images/step4b.png)
     Once the installation is complete, a terminal window will appear. You can close or ignore that for the moment. Open the Files app. You'll notice that there is now a new space for "Linux files."
    ![Screenshot with an arrow pointing at the new "Linux files" option.](images/step4c.png)
5. Plug in the usb drive provided by JSTOR. Drag all the files on the drive over to Linux files, except the one that begins with "elasticsearch-." 
    ````
    Chromebooks can be a little finicky about transferring files. You may note that the time remaining for the transfer swings wildly from several minutes to hundreds of hours. Don’t worry. It’s typically less than an hour. You may want to make sure the computer is plugged in and active during this process (see the guide below for instructions on keeping the chromebook from sleeping during the transfer). It can take a while, and you don’t want to have any errors. If you do, delete the files and start again.
    ````
    ![Screenshot with an arrow pointing at the highlighted files that need to be moved. An arrow points to the Linux files folder where they need to be dragged. A third arrow points to the progress monitor indicating over an hour until the transfer is complete.](images/step5.png)
6. While the files are transfering, let's install Elasticsearch! Double click the file that begins with "elasticsearch-" to begin the installation. You'll be prompted to move forward. Click "Install" then click "OK."
    ![Screenshot with the Elasticsearch installer highlighted.](images/step6a.png)
    ![Screenshot of the first installation step with the "Install" button.](images/step6b.png)
    ![Screenshot of the second installation step with the "OK" button.](images/step6c.png)
    The progress monitor for the Elasticsearch installation may show 0% and not move. That’s okay! When it’s done, the progress monitor will disappear.
    ![Screenshot showing the progress monitor with 0% progress.](images/step6d.png)
7. After the file transfer is complete, you can open the Terminal again. If you left the Terminal window open after step 4, you can disregard this step. Open the Terminal app and select Penguin. 
    ![Screenshot showing the app selection pane with arrows pointing to the the three steps. 1. Open the app selection pane. 2. Type "Terminal" in the search bar. 3. Select the Terminal app.](images/step7a.png)
    ![Screenshot showing where the option for "Penguin" is.](images/step7b.png)

8. In the Terminal, enter the following commands, and hit enter after each one (replace username in the second command with the username you chose in step four, which is also visible in the Terminal prompt). In the screenshot below, the username is "ryan".:
    ````
    chmod +x ~/init.sh
    sudo ~/init.sh -u username
    ````
    ![Screenshot of the terminal window showing the commands and the results.](images/step8.png)
9. You're almost done! Just close the terminal, and then re-open it, using the process described in step 7. After a moment, you'll see some text appear. It will take a couple of minutes to run through the initialization process. At the end, you'll see a note that an http server has started on port 1323.
    ![Screenshot of the terminal text showing the end of the startup process.](images/step9.png)
10. Open a browser and enter the following for the URL:
    ````
    localhost:1323
    ````
    ![Screenshot of the landing page with the term "Rachel Carson" placed in the search bar.](images/step10.png)

11. You're now ready to use JSTOR!
    ![Screenshot showing the search results page with results for "Rachel Carson."](images/step11.png)

## How to keep your chromebook awake during file transfer
1. Make sure your chromebook is plugged in and charging.
2. In the settings menu, select Device>Power.
    ![Screenshot with arrows indicating the options for "Device" and "Power."](images/sleepstep2.png)
3. From the dropdown next to "While charging," select "Keep display on." That's it!
    ![Screenshot showing the dropdown options with "Keep display on" selected.](images/sleepstep3.png)


