<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div class="container">
    <h1>Tickets</h1>
    <b-form @submit="search" class="search-form">
      <b-container>
        <b-row>
          <div class="col-11">
            <b-form-input v-model="query" @keyup.enter="search"
              placeholder="Type to Search" />
          </div>
          <div class="col-1">
            <b-btn @click="search" type="submit" variant="outline-success">
              Search
            </b-btn>
          </div>
        </b-row>
      </b-container>
    </b-form>
    <ticket-list :tickets="tickets" showColumnPicker="true"></ticket-list>
  </div>
</template>

<script>
 import TicketList from '@/components/Tickets/List'
 import Axios from 'axios'

 export default {
   components: {
     TicketList
   },

   data: function () {
     return {
       'query': '',
       'error': '',
       'tickets': null
     }
   },

   methods: {
     search: function () {
       this.$router.push({ name: 'Tickets/SearchPage', query: { q: this.query } })
       this.loadTickets()
     },

     loadTickets: function () {
       let url = '/api/tickets'
       let inst = this

       if (this.query && this.query !== '') {
         url += '?q=' + this.query
       }

       Axios.get(url)
            .then((res) => {
              inst.tickets = res.data
            })
            .catch((err) => {
              if (err.response.status === 404) {
                inst.tickets = []
                return
              }

              // TODO: Visually que the user that there's been an error
              console.log('ERROR', err)
            })
     }
   },

   created: function () {
     this.query = this.$router.currentRoute.query.q
     this.loadTickets()
   }
 }
</script>

<style>
 .search-form {
   margin-bottom: 1rem;
 }
</style>
