<template>
  <div>
    <div v-for="line in msg.split('\n')" :key="line">
      {{ line }}
    </div>
    <br>
    <br>
  </div>
</template>

<script>
export default {
  data () {
    return {
      msg: `不要膜蛤，去学习，人生没有那么多时间让你去习习蛤蛤。
      苟も国家に利すれば生死を以てす、豈禍福に因りて之を避趨するや。`
    }
  }
}
</script>