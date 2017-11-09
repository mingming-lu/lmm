<template>
  <div>
    <p v-if="profile">{{ profile.result.name }}</p>
  </div>
</template>

<script>
import * as request from '@/request'
export default {
  data () {
    request.get('http://localhost:8081/profile', (response) => {
      this.profile = response
    })
    return {
      profile: undefined
    }
  }
}
</script>
