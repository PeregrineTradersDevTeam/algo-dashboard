<template>
  <div
    v-bind:class="['pnl-counter', val >= 0 ? 'pnl-positive' : 'pnl-negative']"
  >
    PnL: {{ val }}
  </div>
</template>

<script>
export default {
  name: "PnlBadge",
  computed: {
    val() {
      return Math.trunc(this.$store.state.lastPnL);
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.pnl-counter {
  display: block;
  height: 22px;
  font-weight: bold;
  font-size: 14px;
  line-height: 22px;
  padding: 1px 6px 1px 6px;
  color: white;
  overflow: none;
  white-space: nowrap;
}

.pnl-positive {
  background-color: green;
}

.pnl-negative {
  background-color: red;
}
</style>
