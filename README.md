# The JSTOR Prison Education Project (PEP)

## Project Overview

## Project Documentation

[Online documentation for installing and using the application](https://ithaka.github.io/PEP/site/)

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

## License

## Contact Info