<template>
  <img
    :src="url"
    :alt="tags.join(' ')"
    :srcset="`${photoThumbnailUrl(url, 1280)} 1280w,
              ${photoThumbnailUrl(url, 960)} 960w,
              ${photoThumbnailUrl(url, 640)} 640w,
              ${photoThumbnailUrl(url, 320)} 320w`"
    sizes="(min-width: 800px) 50vw, 100vw"
    @error="fallbackOnThumbnailFailure($event, url)"
  >
</template>

<script>
export default {
  props: {
    tags: {
      type: Array,
      default: () => {
        ;[]
      }
    },
    url: {
      type: String,
      default: '' // TODO: replace by default image
    }
  },
  methods: {
    photoThumbnailUrl(originalURL, width) {
      const idx = originalURL.lastIndexOf('/')
      const prefix = originalURL.substring(0, idx)
      const suffix = originalURL.substring(idx)
      return `${prefix}/thumbnail/w${width}${suffix}`
    },
    fallbackOnThumbnailFailure(event, url) {
      // sleep 500ms on error
      setTimeout(() => {}, 500)
      event.target.srcset = url
    }
  }
}
</script>
