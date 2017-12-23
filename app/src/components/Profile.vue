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
          <p>
            <div style="width: 100%; height: 16px; border: 1px solid white;">
              <div class="lmm-left" style="margin-right:8px"><b>{{ we.company }}</b></div>
              <div v-if="we.current === true" class="lmm-small lmm-right lmm-level-opacity" style="margin-left:8px;">{{ we.date_from.slice(0, 7) }} ~ Current</div>
              <div v-else class="lmm-small lmm-right lmm-level-opacity" style="margin-left:8px;">{{ we.date_from.slice(0, 7) }} ~ {{ we.date_to.slice(0, 7) }}</div>
              <hr class="lmm-level-opacity" style="display: block">
            </div>
            <p>{{ we.position }}</p>
            <br>
          </p>
        </div>
      </div>

      <div class="lmm-container lmm-card-4 lmm-margin lmm-white">
        <p class="lmm-large"><i class="fa fa-certificate fa-fw lmm-margin-right"></i><b>Education</b></p>
        <br>
        <div v-for="(e, index) in education" :key="index">
          <p>
            <div style="width: 100%; height: 16px; border: 1px solid white;">
              <div class="lmm-left" style="margin-right:8px"><b>{{ e.institution }}</b> ({{ e.degree }})</div>
              <div v-if="e.current === true" class="lmm-small lmm-right lmm-level-opacity" style="margin-left:8px;">{{ e.date_from.slice(0, 7) }} ~ Current</div>
              <div v-else class="lmm-small lmm-right lmm-level-opacity" style="margin-left:8px;">{{ e.date_from.slice(0, 7) }} ~ {{ e.date_to.slice(0, 7) }}</div>
              <hr class="lmm-level-opacity" style="display: block">
            </div>
            <p>{{ e.department }} {{ e.major }}</p>
            <br>
          </p>
        </div>
      </div>

      <div class="lmm-container lmm-card-4 lmm-margin lmm-white">
        <p class="lmm-large"><i class="fa fa-certificate fa-fw lmm-margin-right"></i><b>Qualifications</b></p>
        <br>
        <div v-for="(q, index) in qualifications" :key="index">
          <p>
            <div style="width: 100%; height: 16px; border: 1px solid white;">
              <div class="lmm-left" style="margin-right:8px"><b>{{ q.name }}</b></div>
              <div class="lmm-small lmm-right lmm-level-opacity" style="margin-left:8px;">{{ q.date.slice(0, 7) }}</div>
              <hr class="lmm-level-opacity" style="display: block">
            </div>
            <br>
          </p>
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
      if (response.skills) {
        this.skills = this.skills.concat(response.skills)
      }
      if (response.languages) {
        this.languages = this.languages.concat(response.languages)
      }
      if (response.work_experience) {
        this.workExperience = this.workExperience.concat(response.work_experience)
      }
      if (response.education) {
        this.education = this.education.concat(response.education)
      }
      if (response.qualifications) {
        this.qualifications = this.qualifications.concat(response.qualifications)
      }
    })
    return {
      name: '',
      avatar_url: '//:0',
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
  }
}
</script>
