<template>
  <div class="row">
    <!-- Left Column -->
    <div class="left" style="width:33.3333%;">

      <!-- profile -->
      <div>
        <img :src="avatar_url" style="width:100%" alt="Avatar" @error="setDefaultAvatar">
        <div class="container">
          <h3>{{ name }}</h3>
          <p>{{ bio }}</p>
          <p><i class="fa fa-briefcase fa-fw"></i>{{ profession }}</p>
          <p><i class="fa fa-home fa-fw"></i>{{ location }}</p>
          <p><i class="fa fa-envelope fa-fw"></i>{{ email }}</p>
        </div>
      </div>

      <!-- skills -->
      <div class="container">
        <h3><i class="fa fa-asterisk fa-fw"></i>Skills</h3>
        <span v-for="skill in skills" :key="skill.name">
          {{ skill.name }}
        </span>
      </div>

      <div class="container">
        <h3><i class="fa fa-globe fa-fw"></i>Languages</h3>
        <span v-for="language in languages" :key="language.name">
          {{ language.name }}
        </span>
      </div>
    </div>
    <!-- End Left Column -->

    <!-- Right Column -->
    <div class="right" style="width:66.6666%;">

      <!-- work experience -->
      <div class="container">
        <h3><i class="fa fa-suitcase fa-fw"></i>Work Experience</h3>
        <div v-for="(we, index) in workExperience" :key="we.company">
          <h3 class="span">{{ we.company }}</h3>
          <p v-if="we.current === true" class="span right text-right opacity-plus">{{ we.date_from.slice(0, 7) }} ~ Current</p>
          <p v-else class="span right text-right opacity-plus">{{ we.date_from.slice(0, 7) }} ~ {{ we.date_to.slice(0, 7) }}</p>
          <p>{{ we.position }}</p>
          <hr v-if="index !== workExperience.length-1" class="opacity-plus">
        </div>
      </div>

      <div class="container">
        <h3><i class="fa fa-book fa-fw"></i>Education</h3>
        <div v-for="(e, index) in education" :key="index">
          <h3 class="span">{{ e.institution }} <div class="span no-weight">({{ e.degree }})</div></h3>
          <p v-if="e.current === true" class="span right text-right opacity-plus">{{ e.date_from.slice(0, 7) }} ~ Current</p>
          <p v-else class="span right text-right opacity-plus">{{ e.date_from.slice(0, 7) }} ~ {{ e.date_to.slice(0, 7) }}</p>
          <p>{{ e.department }} {{ e.major }}</p>
          <hr v-if="index !== workExperience.length-1" class="opacity-plus">
        </div>
      </div>

      <div class="container">
        <h3><i class="fa fa-certificate fa-fw"></i>Qualifications</h3>
        <div v-for="(q, index) in qualifications" :key="index">
          <h3 class="span left">{{ q.name }}</h3>
          <p class="span right text-right opacity-plus">{{ q.date.slice(0, 7) }}</p>
        </div>
      </div>

    </div>
    <!-- End Right Column -->

  </div>
</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      name: '',
      avatar_url: '',
      bio: '',
      profession: '',
      location: '',
      email: '',
      skills: [],
      languages: [],
      workExperience: [],
      education: [],
      qualifications: []
    }
  },
  created () {
    axios.get('http://api.lmm.local/profile').then((res) => {
      let profile = res.data
      this.name = profile.name
      this.avatar_url = profile.avatar_url
      this.bio = profile.bio
      this.profession = profile.profession
      this.location = profile.location
      this.email = profile.email
      if (profile.skills) {
        this.skills = this.skills.concat(profile.skills)
      }
      if (profile.languages) {
        this.languages = this.languages.concat(profile.languages)
      }
      if (profile.work_experience) {
        this.workExperience = this.workExperience.concat(profile.work_experience)
      }
      if (profile.education) {
        this.education = this.education.concat(profile.education)
      }
      if (profile.qualifications) {
        this.qualifications = this.qualifications.concat(profile.qualifications)
      }
    }).catch((e) => {
      console.log(e)
    })
  },
  methods: {
    setDefaultAvatar () {
      this.avatar_url = 'https://avatars3.githubusercontent.com/u/17140497?s=400&u=636be90e7798e07230fa5f37af1a0f5070fa23a6&v=4'
    }
  }
}
</script>
