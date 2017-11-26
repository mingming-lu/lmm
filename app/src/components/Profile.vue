<template>
  <div class="lmm-row">
      
    <!-- Left Column -->
    <div class="lmm-left" style="text-align:left; width:33.3333%; display:inline-block">
      <div class="lmm-white lmm-margin lmm-card-4">
        <div class="lmm-display-container">
          <img :src="avatar_url" style="width:100%" alt="Avatar">
          <div class="lmm-display-bottomleft lmm-container">
            <h2>{{ name }}</h2>
          </div>
        </div>
        <div class="lmm-container">
          <p><i class="fa fa-briefcase fa-fw lmm-margin-right lmm-large"></i>{{ profession }}</p>
          <p><i class="fa fa-home fa-fw lmm-margin-right lmm-large"></i>{{ location }}</p>
          <p><i class="fa fa-envelope fa-fw lmm-margin-right lmm-large"></i>{{ email }}</p>
          <hr>

          <p><b><i class="fa fa-asterisk fa-fw lmm-margin-right"></i>Skills</b></p>
          <div v-for="(skill, index) in skills" :key="index">
            <p>{{ skill }}</p>
          </div>

          <p><b><i class="fa fa-globe fa-fw lmm-margin-right"></i>Languages</b></p>
          <div v-for="(language, index) in languages" :key="index">
            <p>{{ language }}</p>
          </div>
          <br>
        </div>
      </div><br>
    <!-- End Left Column -->
    </div>

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
