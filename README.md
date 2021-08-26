# Prison Education: Offline Access Initiative

## Project Overview

### Background
In 2007, at the request of the [Bard Prison Initiative](https://bpi.bard.edu/), JSTOR created a tool to help incarcerated students who do not have access to the internet.  The tool was an offline browser that provided a searchable index of content on JSTOR, but not the full-text of documents.  Students use it to find content to request of instructors and librarians, who are then able to get the content outside of the prison and bring in a printed version.  Nearly twenty prison education programs now have this offline index.  

### Project goal
The goal of the project is to create a next-generation tool to support incarcerated students conducting research without access to the internet.  The solution will be tested with a cohort of prison education programs, in order to make a recommendation regarding how to provide full access to JSTOR to as many higher education in prison programs as possible.

This work performed is paired with and informed by new research on prison education being conducted by [Ithaka S+R](https://sr.ithaka.org/).

### The work
1. A one-day workshop was conducted with an advisory committee of prison education program leaders, librarians, educators, former students and department of corrections representatives.
2. With guidance from the advisory committee, a test cohort of five prison education programs was selected to implement and test the new solution. The five programs are Calvin Prison Initiative, Cornell Prison Education Program, Freedom Education Project Puget Sound, Prison University Project, and Community University Project at Stetson University.
3. A one-day workshop with the test cohort was conducted, during which features for the proposed solution were refined and prioritized.
4. The new solution was implemented.

The new solution is currently being shipping to the test cohort, who will use it with their students and teachers.  The feedback for this test cohort will be used to assess how best to expand access to the program.

### The solution
The new solution offers a number of advancements that will help increase access for incarcerated researchers to scholarly material. Installed on a small server called a NUC that is not much bigger than a paperback book, the solution contains the JSTOR search index, the full-text content for all open access content, as well as workflow software both to configure the system and to manage students’ requests. When they access it, students conduct searches much as they would on the main JSTOR platform, except that when they find an article, rather than click through to read it, they request the article, where it gets moved into an administrative request queue. This allows administrators to review and fulfill the requests. When connected to the internet, the device will automatically update both functionality and content, dramatically improving the ability to maintain and improve these systems moving forward.

Each device contains JSTOR's solution for prison education. 

### Project Documentation

[Documentation for installing and using the Prison Education: Offline Access Initiative application](https://ithaka.github.io/PEP/site/)

## Code

### Architecture Overview

This project is designed to run in multiple environments, the primary of which is an Intel NUC with a customized Arch Linux installtion.  The software is also designed to run off of a thumbdrive on Window's computers for the purposes of providing access in environments where remote access to an Intel NUC is not possible.  The software is broken up into two primary components with some helper software to support the two primary components.

#### App

The [app](app) folder contains a javascript website built on Vue.js that provides the primary user interface for the project.  It allows students to search, browse, and request documents for approval which can then be provided to them.  In addition to student access, it provides a password protected administration panel for administrators to view and manage requests, and build a thumb drive containing a chosen subset of the content on the NUC for use in situations where the NUC itself cannot be used.  All of this functionality is powered by the API.

#### API

The [API](api) package contains the API which provides communication between the web interface, the Elasticsearch index, PDFs repository, and a requests database.  It facilitates taking user search requests, properly formulating an Elasticsearch query, and returning those results to the web app.  Additionally, it provides an interface to modify the requests database, build thumb drives based on a set of parameters from the web app, and retrieve documents if available.

The API is written in Go to facilitate simple and easy cross-platform function, without the need of an extra interpreter or virtual machine on the target machine.

#### Firmware Builder

The (firmware builder)[pep-firmware-builder] is a helper utility to customize and build a new device image for a NUC.

#### Startup tool

The last component, the (startup tool)[pep-linux-startup] is a small program which is included in the operating system image to ensure a sane environment upon system bootup.

#### Elasticsearch

Elasticsearch 7.x comes bundled on the NUC, and is included in the operating system built by the Firmware Builder.  Additionally, a version for Windows is placed on the thumb drives produced by the API, and is accompanied by a Windows compatible copy of OpenJDK to facilitate running Elasticsearch.  Elasticsearch is loaded with an index of JSTOR content for the purposes of providing search functionality.

### Issue Tracking

:beetle:&nbsp;&nbsp;[Problem reports and enhancement requests](https://github.com/JSTOR-Labs/pep/issues)

### Discussion

Additional information can be found on the [PEP Discussions](https://github.com/JSTOR-Labs/juncture/discussions) site.  This includes sections for:

📢 &nbsp;&nbsp;[Announcements](https://github.com/JSTOR-Labs/pep/discussions/categories/announcements)  
❓ &nbsp;&nbsp;[Q&A with the PEP team and development community](https://github.com/JSTOR-Labs/pep/discussions/categories/q-a)  
💡 &nbsp;&nbsp;[Suggestions for new features](https://github.com/JSTOR-Labs/pep/discussions/categories/ideas)  

## About Us

[JSTOR](https://about.jstor.org/) is part of [ITHAKA](https://www.ithaka.org/), a not-for-profit dedicated to expanding access to knowledge and education worldwide.

The [Prison Education: Offline Access Initiative](https://labs.jstor.org/projects/prison-education/) was developed by the [JSTOR Labs Team](https://labs.jstor.org) in collaboration with [Ithaka S+R](http://labs.jstor.org/projects/prison-education/) under funding provided by [The Andrew W. Mellon Foundation](https://www.mellon.org/).

## License

This project is available under the MIT license.
For more information, [view the full license and copyright notice](./LICENSE).

Copyright 2021 Ithaka Harbors, Inc.
