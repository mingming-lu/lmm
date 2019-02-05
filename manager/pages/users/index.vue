<template>
  <v-data-table
    :headers="headers"
    :items="users"
    :pagination.sync="pagination"
    :rows-per-page-items="pagination.rowsPerPageItems"
    :total-items="pagination.totalItems"
    :loading="loading"
    class="elevation-1"
  >
    <template slot="items" slot-scope="props">
      <td
        width="40%"
        :class="{'primary--text': props.item.name === me}"
        >
        {{ props.item.name }}
      </td>
      <td width="40%">
        <v-flex xs12 sm6 d-flex>
          <v-select
            :disabled="props.item.name === me"
            :items="roles"
            v-model="props.item.role"
            @input="assignUserRole(props.item.name, props.item.role, $event)"
          >
          </v-select>
        </v-flex>
      </td>
      <td width="20%">{{ formattedTime(props.item.registered_date) }}</td>
    </template>
  </v-data-table>
</template>

<script>
import {
  buildURLEncodedQueryString,
  formattedDateFromTimeStamp,
} from '~/assets/js/utils'

const usersHandler = httpClient => {
  return {
    fetchByQuery: (queryParams) => {
      let query = buildURLEncodedQueryString(queryParams)
      query = query ? `?${query}` : ''

      return httpClient.get(`/v1/users${query}`, {
        headers: {
          Authorization: `Bearer ${window.localStorage.getItem('accessToken')}`,
        },
      })
    },
    changeRole: (username, role) => {
      return httpClient.put(`/v1/users/${username}/role`, {
        role: role,
      }, {
        headers: {
          Authorization: `Bearer ${window.localStorage.getItem('accessToken')}`,
        },
      })
    },
  }
}

export default {
  data () {
    return {
      roles: ['admin', 'ordinary'],
      headers: [
        { text: 'Name',            value: 'name'            },
        { text: 'Role',            value: 'role'            },
        { text: 'Registered Date', value: 'registered_date' },
      ],
      users: [],
      me: '',
      loading: true,
      pagination: {
        descending:       true,
        page:             1,
        sortBy:           'registered_date',
        totalItems:       0,
        rowsPerPage:      50,
        rowsPerPageItems: [50, 100],
      },
    }
  },
  methods: {
    formattedTime(dtString) {
      return formattedDateFromTimeStamp(dtString)
    },
    assignUserRole(userName, oldRole, newRole) {
      if (!confirm(`going to change the role of "${userName}" into ${newRole}`)) {
        return
      }
      usersHandler(this.$axios).changeRole(userName, oldRole, newRole).catch(err => {
        if (!err.response) {
          console.log(err)
          return
        }
        alert(err.response.data)
      })
    }
  },
  watch: {
    pagination: {
      handler(newOne, oldOne) {
        this.loading = true
        const query = {
          page:    newOne.page,
          count:   newOne.rowsPerPage,
          sort_by: newOne.sortBy, 
          sort:    newOne.descending ? 'desc' : 'asc',
        }
        usersHandler(this.$axios).fetchByQuery(query).then(res => {
          this.me                    = res.headers['x-lmm-user']
          this.users                 = res.data.users
          this.pagination.sortBy     = res.data.sort_by
          this.pagination.descending = res.data.sort === 'desc' ? true : false
          this.pagination.page       = res.data.page
          this.loading               = false
        })
      },
      deep: true,
    },
  },
}
</script>
