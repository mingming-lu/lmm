<template>
  <section class="container pagination">
    <input
      type="button"
      class="button"
      :class="{enable: page > 1}"
      @click="prevPage"
      value="Prev"
      />
    <div v-for="item in items" :key="item">
      <input
        type="button"
        :class="{active: item === page, enable: item !== page}"
        class="button page"
        :page="item"
        :value="item"
        />
    </div>
    <input
      type="button"
      class="button"
      :class="{enable: page < total}"
      @click="nextPage"
      value="Next"
      />
  </section>
</template>
<script>
import { range } from '~/assets/js/utils'

const maxPipe = 5

export default {
  props: {
    page: {
      type: Number,
      default: 1,
      validator: v => {
        return Number.isInteger(v) && v > 0
      },
    },
    total: {
      type: Number,
      default: 1,
      validator: v => {
        return Number.isInteger(v) && v > 0
      }
    }
  },
  data: () => {
    return {
      items: [],
    }
  },
  created() {
    this.items = range(1, this.total)
  },
  methods: {
    prevPage() {
      if (this.page > 1) {
        --this.page
      }
    },
    nextPage() {
      if (this.page < this.total) {
        ++this.page
      }
    }
  }
}
</script>
<style lang="scss" scoped>
@import '~/assets/scss/styles.scss';
.button {
  background-color: transparent;
  border: none;
  margin: 8px;
  font-size: 1.1em;
  opacity: 0.2;
  &:focus {
    outline: 0;
    box-shadow: none;
  }
  &.active {
    color: $color_accent;
    opacity: 1;
  }
  &.enable {
    opacity: 1;
    &:hover {
      background-color: $color_accent;
      color: white;
      cursor: pointer;
    }
  }
}
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>

