<template>
  <div class="lmm-row">
      
    <!-- Left Column -->
    <div class="lmm-left" style="text-align:left; width:33.3333%; display:inline-block">
      <div class="lmm-white lmm-margin lmm-card-4">
        <div class="lmm-display-container">
          <img :src="avatar_url" style="width:100%" alt="Avatar">
        </div>
        <div class="lmm-container">
          <p class="lmm-large"><b>{{ name }}</b></p>
          <p>{{ bio }}</p>
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
      this.name = response.name
      this.avatar_url = response.avatar_url
      this.bio = response.bio
      this.profession = response.profession
      this.location = response.location
      this.email = response.email
      this.skills = this.skills.concat(response.skills)
      this.languages = this.languages.concat(response.languages)
    })
    return {
      name: undefined,
      avatar_url: undefined,
      bio: undefined,
      profession: undefined,
      location: undefined,
      email: undefined,
      skills: [],
      languages: []
    }
  }
}
</script>
