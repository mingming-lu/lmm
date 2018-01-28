<template>
  <div class="container">
    <!-- Left Column -->
    <div :class="{ 'left': !isMobile }">

      <!-- profile -->
      <div class="container">
        <img :src="avatar_url" style="width:100%" alt="Avatar" @error="setDefaultAvatar">
        <h3>{{ name }}</h3>
        <p>{{ bio }}</p>
        <p><i class="fa fa-briefcase fa-fw"></i>{{ profession }}</p>
        <p><i class="fa fa-map-marker fa-fw"></i>{{ location }}</p>
        <p><i class="fa fa-envelope fa-fw"></i>{{ email }}</p>
      </div>

      <!-- skills -->
      <div class="container">
        <h3><i class="fa fa-asterisk fa-fw"></i>Skills</h3>
        <router-link to="" v-for="skill in skills" :key="skill.name" class="link tag">
          {{ skill.name }}
        </router-link>
        <hr class="transparent">
      </div>

      <div class="container">
        <h3><i class="fa fa-globe fa-fw"></i>Languages</h3>
        <router-link to="" v-for="language in languages" :key="language.name" class="link tag">
          {{ language.name }}
        </router-link>
        <hr class="transparent">
      </div>
    </div>
    <!-- End Left Column -->

    <!-- Right Column -->
    <div :class="{ 'right': !isMobile }">

      <!-- work experience -->
      <div class="container">
        <h3><i class="fa fa-building fa-fw"></i>Work Experience</h3>
        <br>
        <div v-for="we in workExperience" :key="we.company">
          <div class="row">
            <div class="left h">{{ we.company }}</div>
            <div v-if="we.current === true" class="right opacity">{{ we.date_from.slice(0, 7) }} ~ Current</div>
            <div v-else class="right opacity">{{ we.date_from.slice(0, 7) }} ~ {{ we.date_to.slice(0, 7) }}</div>
            <hr class="opacity-plus connection">
          </div>
          <p class="opacity">{{ we.position }}</p>
          <br>
        </div>
      </div>

      <div class="container">
        <h3><i class="fa fa-graduation-cap fa-fw"></i>Education</h3>
        <br>
        <div v-for="(e, index) in education" :key="index">
          <div class="row">
            <div class="left h">{{ e.institution }}</div>
            <div v-if="e.current === true" class="right opacity">{{ e.date_from.slice(0, 7) }} ~ Current</div>
            <div v-else class="right opacity">{{ e.date_from.slice(0, 7) }} ~ {{ e.date_to.slice(0, 7) }}</div>
            <hr class="opacity-plus connection">
          </div>
          <div class="opacity">
            <p>{{ e.degree }}</p>
            <p>{{ e.department }}</p>
            <p>{{ e.major }}</p>
          </div>
          <br>
        </div>
      </div>

      <div class="container">
        <h3><i class="fa fa-certificate fa-fw"></i>Qualifications</h3>
        <br>
        <div v-for="(q, index) in qualifications" :key="index">
          <div class="row">
            <div class="left h">{{ q.name }}</div>
            <div class="right opacity">{{ q.date.slice(0, 7) }}</div>
            <hr class="opacity-plus connection">
          </div>
          <br>
        </div>
        <br>
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
      qualifications: [],
      isMobile: false
    }
  },
  created () {
    this.calcIsMobile()
    window.addEventListener('resize', this.calcIsMobile)
    axios.get('http://api.lmm.im/users/1/profile').then((res) => {
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
  beforeDestroy () {
    window.removeEventListener('resize', this.calcIsMobile)
  },
  methods: {
    setDefaultAvatar () {
      this.avatar_url = 'https://avatars3.githubusercontent.com/u/17140497?s=400&u=636be90e7798e07230fa5f37af1a0f5070fa23a6&v=4'
    },
    calcIsMobile () {
      this.isMobile = window.innerWidth <= 825
    }
  }
}
</script>

<style scoped>
.container > .left {
  width: 33.333333%;
}
.container > .right {
  width: 66.666666%;
}
.row {
  width: 100%;
  border: 1px solid white;
}
.row .left {
  margin-right:8px;
}
.row .right {
  margin-left:8px;
}
.connection {
  display: block;
}
.h {
  font-size: 1.17em;
}
</style>

