<template>
  <div class="lmm-row">
      
    <!-- Left Column -->
    <div class="lmm-left" style="text-align:left; width:33.3333%; display:inline-block">
      <div class="lmm-margin lmm-card-4">
        <div class="lmm-display-container">
          <img :src="avatar_url" style="width:100%" alt="Avatar">
        </div>
        <div class="lmm-container">
          <p class="lmm-xlarge"><b>{{ name }}</b></p>
          <p>{{ bio }}</p>
        </div>
        <div class="lmm-container">
          <p><i class="fa fa-briefcase fa-fw lmm-margin-right lmm-large"></i>{{ profession }}</p>
          <p><i class="fa fa-home fa-fw lmm-margin-right lmm-large"></i>{{ location }}</p>
          <p><i class="fa fa-envelope fa-fw lmm-margin-right lmm-large"></i>{{ email }}</p>
        </div>
      </div>

      <div class="lmm-container lmm-white lmm-margin lmm-card-4">
        <p class="lmm-large"><b><i class="fa fa-asterisk fa-fw lmm-margin-right"></i>Skills</b></p>
        <div v-for="(skill, index) in skills" :key="index">
          <p>
            <div style="width: 100%; height: 16px; border: 1px solid white;">
              <div class="lmm-left" style="margin-right:8px">{{ skill.name }}</div>
              <div class="lmm-small lmm-right lmm-level-opacity" style="margin-left:8px;">{{ skill.level }}</div>
              <hr class="lmm-level-opacity" style="display: block">
            </div>
          </p>
        </div>
      </div>

      <div class="lmm-container lmm-white lmm-margin lmm-card-4">
        <p class="lmm-large"><b><i class="fa fa-globe fa-fw lmm-margin-right"></i>Languages</b></p>
        <div v-for="(language, index) in languages" :key="index">
          <p>
            <div style="width: 100%; height: 16px; border: 1px solid white;">
              <div class="lmm-left" style="margin-right:8px">{{ language.name }}</div>
              <div class="lmm-small lmm-right lmm-level-opacity" style="margin-left:8px;">{{ language.level }}</div>
              <hr class="lmm-level-opacity" style="display: block">
            </div>
          </p>
        </div>
      </div>
    </div>
    <!-- End Left Column -->

    <!-- Right Column -->
    <div class="lmm-right" style="text-align:left; width:66.6666%; display:inline-block">
    
      <div class="lmm-container lmm-white lmm-margin lmm-card-4">
        <p class="lmm-large"><b><i class="fa fa-suitcase fa-fw lmm-margin-right"></i>Work Experience</b></p>
        <br>
        <div v-for="(we, index) in workExperience" :key="index">
          <p v-if=" we.current === true" class="lmm-opacity"><i class="fa fa-calendar fa-fw lmm-margin-right"></i>{{ we.year_from }}/{{ we.month_from }} ~ current</p>
          <p v-else class="lmm-opacity"><i class="fa fa-calendar fa-fw lmm-margin-right"></i>{{ we.year_from }}/{{ we.month_from }} ~ {{ we.year_to }}/{{ we.month_to }}</p>
          <p><b>{{ we.company }}</b> - {{ we.position }} ({{ we.status }})</p>
          <br v-if="index === workExperience.length - 1">
          <hr v-else class="lmm-level-opacity">
        </div>
      </div>

      <div class="lmm-container lmm-card-4 lmm-margin lmm-white">
        <p class="lmm-large"><i class="fa fa-certificate fa-fw lmm-margin-right"></i><b>Education</b></p>
        <br>
        <div v-for="(e, index) in education" :key="index">
          <p v-if=" e.current === true" class="lmm-opacity"><i class="fa fa-calendar fa-fw lmm-margin-right"></i>{{ e.year_from }}/{{ e.month_from }} ~ current</p>
          <p v-else class="lmm-opacity"><i class="fa fa-calendar fa-fw lmm-margin-right"></i>{{ e.year_from }}/{{ e.month_from }} ~ {{ e.year_to }}/{{ e.month_to }}</p>
          <p><b>{{ e.institution }}</b> ({{ e.degree }})</p>
          <p>{{ e.department }} {{ e.major }}</p>
          <br v-if="index === education.length - 1">
          <hr v-else class="lmm-level-opacity">
        </div>
      </div>

      <div class="lmm-container lmm-card-4 lmm-margin lmm-white">
        <p class="lmm-large"><i class="fa fa-certificate fa-fw lmm-margin-right"></i><b>Qualifications</b></p>
        <br>
        <div v-for="(q, index) in qualifications" :key="index">
          <p class="lmm-opacity"><i class="fa fa-calendar fa-fw lmm-margin-right"></i>{{ q.year }}/{{ q.month }}</p>
          <p><b>{{ q.name }}</b></p>
          <br v-if="index === qualifications.length - 1">
          <hr v-else class="lmm-level-opacity">
        </div>
      </div>

    </div>
    <!-- End Right Column -->

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
      this.workExperience = this.workExperience.concat(response.work_experience)
      this.education = this.education.concat(response.education)
      this.qualifications = this.qualifications.concat(response.qualifications)
    })
    return {
      name: undefined,
      avatar_url: undefined,
      bio: undefined,
      profession: undefined,
      location: undefined,
      email: undefined,
      skills: [],
      languages: [],
      workExperience: [],
      education: [],
      qualifications: []
    }
  }
}
</script>
