<template>
  <div class="iw-body-chart" @click="$emit('open-chart')">
    <img
      style="padding-top:1px"
      v-bind:src="url"
      :width="width"
      :height="height"
    />
  </div>
</template>

<script>
export default {
  name: "InstanceWidgetBodyChart",
  components: {},
  props: {
    width: {
      type: Number,
      default: 184,
    },
    height: {
      type: Number,
      default: 68,
    },
    id: String,
    mode: String,
  },
  data: function() {
    return {
      smallCurves: "pnl,op,lo,hi,l,sahp", // the order of curves is important
      reqNum: Date.now(),
    };
  },
  computed: {
    url() {
      return `/api/charter/M:${this.id}/img?w=${this.width}&h=${this.height}&a=${this.smallCurves}&req=${this.reqNum}`;
    },
  },
};
</script>

v<style scoped>
.iw-body-chart {
  padding-left: 1px;
  margin: 0;
  display: block;
  height: 70px;
  background-color: #363639;
}
</style>
