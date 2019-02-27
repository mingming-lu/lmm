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
    <template 
      slot="items" 
      slot-scope="props">
      <td
        :class="{'primary--text': props.item.name === me}"
        width="40%"
      >
        {{ props.item.name }}
      </td>
      <td width="40%">
        <v-flex 
          xs12 
          sm6 
          d-flex>
          <v-select
            :disabled="props.item.name === me"
            :items="roles"
            v-model="props.item.role"
            @input="assignUserRole(props.item.name, props.item.role, $event)"
          />
        </v-flex>
      </td>
      <td width="20%">{{ formattedTime(props.item.registered_date) }}</td>
    </template>
  </v-data-table>
</template>

<script>
import {
  buildURLEncodedQueryString,
  formattedDateFromTimeStamp
} from '~/assets/js/utils'

const usersHandler = httpClient => {
  return {
    fetchByQuery: queryParams => {
      let query = buildURLEncodedQueryString(queryParams)
      query = query ? `?${query}` : ''

      return httpClient.get(`/v1/users${query}`, {
        headers: {
          Authorization: `Bearer ${window.localStorage.getItem('accessToken')}`
        }
      })
    },
    changeRole: (username, role) => {
      return httpClient.put(
        `/v1/users/${username}/role`,
        {
          role: role
        },
        {
          headers: {
            Authorization: `Bearer ${window.localStorage.getItem(
              'accessToken'
            )}`
          }
        }
      )
    }
  }
}

export default {
  asyncData({ $axios }) {
    return usersHandler($axios)
      .fetchByQuery({})
      .then(res => {
        return {
          loading: false,
          users: res.data.users,
          me: res.headers['x-lmm-user'],
          pagination: {
            sortBy: res.data.sort_by,
            descending: res.data.sort === 'desc' ? true : false,
            page: res.data.page,
            rowsPerPage: res.data.count,
            totalItems: res.data.total,
            rowsPerPageItems: [50, 100]
          }
        }
      })
  },
  data() {
    return {
      roles: ['admin', 'ordinary'],
      headers: [
        { text: 'Name', value: 'name' },
        { text: 'Role', value: 'role' },
        { text: 'Registered Date', value: 'registered_date' }
      ],
      loading: true
    }
  },
  watch: {
    pagination: {
      handler(newOne, oldOne) {
        if (Object.keys(newOne).every(key => newOne[key] === oldOne[key])) {
          return
        }
        this.loading = true
        const query = {
          page: newOne.page,
          count: newOne.rowsPerPage,
          sort_by: newOne.sortBy,
          sort: newOne.descending ? 'desc' : 'asc'
        }
        usersHandler(this.$axios)
          .fetchByQuery(query)
          .then(res => {
            this.me = res.headers['x-lmm-user']
            this.users = res.data.users
            this.pagination.sortBy = res.data.sort_by
            this.pagination.descending = res.data.sort === 'desc' ? true : false
            this.pagination.page = res.data.page
            this.pagination.rowsPerPage = res.data.count
            this.pagination.totalItems = res.data.total
            this.loading = false
          })
      },
      deep: true
    }
  },
  methods: {
    formattedTime(dtString) {
      return formattedDateFromTimeStamp(dtString)
    },
    assignUserRole(userName, oldRole, newRole) {
      if (
        !confirm(`going to change the role of "${userName}" into ${newRole}`)
      ) {
        return
      }
      usersHandler(this.$axios)
        .changeRole(userName, oldRole, newRole)
        .catch(err => {
          if (!err.response) {
            console.log(err)
            return
          }
          alert(err.response.data)
        })
    }
  }
}
</script>
