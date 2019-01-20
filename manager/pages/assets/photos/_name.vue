<template>
  <img :src="src" :alt="alt">
</template>

<script>

const photoFetcher = (httpClient) => {
  return {
    fetch: (name) => {
      return httpClient.get(`/v1/assets/photos/${name}`, {
        headers: {
          Authorization: `Bearer ${window.localStorage.getItem('accessToken')}`,
        },
      })
    },
  }
}

export default {
  asyncData({ $axios, params }) {
    return photoFetcher($axios).fetch(params.name)
      .then(res => {
        return {
          src:  `${process.env.ASSET_URL}/photos/${res.data.name}`,
          alt: res.data.alts.join(' '),
        }
      })
  },
}
</script>

<style scoped>
</style>
